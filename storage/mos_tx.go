package storage

import (
	"context"
	"fmt"
	"log"
	"sync"

	"tschwaa.com/api/models"
	"tschwaa.com/api/utils"
)

type UpdateSessionMembersParams struct {
	OrganizationID uint64
	SessionID      uint64
	Memberships    []models.Membership
}

type insertMOSResponse struct {
	MemberId uint64
	MOS      *models.MembersOfSession
	Error    string
}

func (store *SQLStorage) UpdateSessionMembersTx(ctx context.Context, arg UpdateSessionMembersParams) ([]*models.MembersOfSession, error) {
	responses := make([]*models.MembersOfSession, 0, len(arg.Memberships))

	err := store.execTx(ctx, func(q *Queries) error {
		// 0. Delete organization'smembers of that session id
		err := store.RemoveAllMembersFromSession(ctx, RemoveAllMembersFromSessionParams{
			OrganizationID: arg.OrganizationID,
			SessionID:      arg.SessionID,
		})
		if err != nil {
			return utils.Fail(
				fmt.Sprintf("error when removing all organization[%d]'s members from session[%d]", arg.OrganizationID, arg.SessionID),
				"ERR_UPD_SESS_MBR_01",
				err,
			)
		}

		// 1. Insert the (membership, session) into MembersOfSession
		wg := new(sync.WaitGroup)
		wg.Add(len(arg.Memberships))
		insertMOSChannel := make(chan insertMOSResponse)
		for _, membership := range arg.Memberships {
			go func(membership models.Membership, channel chan insertMOSResponse, wg *sync.WaitGroup) {
				defer wg.Done()
				mos, err := store.AddMemberToSession(ctx, AddMemberToSessionParams{
					MembershipID: membership.ID,
					SessionID:    arg.SessionID,
				})
				if err != nil {
					channel <- insertMOSResponse{
						MemberId: membership.MemberID,
						Error:    err.Error(),
						// Membership: membership,
						MOS: nil,
					}
				} else {
					channel <- insertMOSResponse{
						MemberId: membership.MemberID,
						Error:    "",
						// Membership: membership,
						MOS: mos,
					}
				}
			}(membership, insertMOSChannel, wg)
		}

		go func() {
			wg.Wait()
			close(insertMOSChannel)
		}()

		for val := range insertMOSChannel {
			log.Println("Channel : ", val)
			responses = append(responses, val.MOS)
		}

		return nil
	})

	return responses, err
}
