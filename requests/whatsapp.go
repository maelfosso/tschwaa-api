package requests

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	"tschwaa.com/api/models"
	"tschwaa.com/api/utils"
)

func getWhatsappRequestURL(url string) string {
	return fmt.Sprintf(
		url,
		os.Getenv("WHATSAPP_API_VERSION"), os.Getenv("WHATSAPP_PHONE_NUMBER_ID"),
	)
}

var jsonForSendMessageText = `{
	"messaging_product": "whatsapp",
	"recipient_type": "individual",
	"to": "%s",
	"type": "text",
	"text": {
			"preview_url": false,
			"body": "%s"
	}
}`

func SendMessageText(to, message string) (*WhatsappSendMessageResponse, error) {
	jsonBody := []byte(fmt.Sprintf(jsonForSendMessageText, to, message))
	bodyReader := bytes.NewReader(jsonBody)
	log.Println("send message text data : ", string(jsonBody))

	requestUrl := getWhatsappRequestURL("https://graph.facebook.com/%s/%s/messages")
	log.Println("send message text request url: ", requestUrl)

	req, err := http.NewRequest(
		http.MethodPost,
		requestUrl,
		bodyReader,
	)
	if err != nil {
		return nil, fmt.Errorf("client: could not create request: %s", err)
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", os.Getenv("WHATSAPP_USER_ACCESS_TOKEN")))

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("client: error making http request: %s", err)
	}

	log.Println("Client: got response!")
	log.Println("client: status code: %d", res.StatusCode)

	resBody, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, fmt.Errorf("client: could not read response body: %s", err)
	}

	log.Println("client: response body: %s", string(resBody))

	var data WhatsappSendMessageResponse
	err = json.Unmarshal(resBody, &data)
	if err != nil {
		return nil, fmt.Errorf("error when unmarshilling response body : %s", err)
	}

	return &data, nil
}

var jsonForSendMessageTemplateText = `{
	"messaging_product": "whatsapp",
	"recipient_type": "individual",
	"to": "%s",
	"type": "template",
	"template": {
		"name": "%s",
		"language": {
			"code": "%s"
		},
		"components": %s
	}
}`

func sendMessageTextFromTemplate(to, template, language, components string) (*WhatsappSendMessageResponse, error) {
	jsonBody := []byte(
		fmt.Sprintf(
			jsonForSendMessageTemplateText, to, template, language, components,
		),
	)
	bodyReader := bytes.NewReader(jsonBody)

	requestUrl := getWhatsappRequestURL("https://graph.facebook.com/%s/%s/messages")
	log.Println("request url: ", requestUrl)

	req, err := http.NewRequest(
		http.MethodPost,
		requestUrl,
		bodyReader,
	)
	if err != nil {
		return nil, utils.Fail("client: could not create request", "ERR_SMSG_TPL_01", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", os.Getenv("WHATSAPP_USER_ACCESS_TOKEN")))

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("client: could not read response body: %w", err)
	}

	log.Println("client: got response!")
	log.Println("client: status code: %d", res.StatusCode)

	resBody, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, utils.Fail("client: could not read response body", "ERR_SMSG_TPL_02", err)
	}

	log.Println("client: response body: %s", string(resBody))

	if res.StatusCode != http.StatusOK {
		log.Println("\nError: \n\n", string(resBody))

		var errData WhatsappSendMessageErrorResponse
		err = json.Unmarshal(resBody, &errData)
		if err != nil {
			return nil, utils.Fail("error when unmarshelling response body error", "ERR_SMSG_TPL_04", err)
		}

		return nil, fmt.Errorf("Unbable to send whatsapp message")
	}
	var data WhatsappSendMessageResponse
	err = json.Unmarshal(resBody, &data)
	if err != nil {
		return nil, utils.Fail("error when unmarshelling response body", "ERR_SMSG_TPL_03", err)
	}

	return &data, nil
}

func SendTschwaaOtp(to, language, pinCode string) (*WhatsappSendMessageResponse, error) {
	parameters := fmt.Sprintf(`
		[
			{
				"type": "body",
				"parameters": [
						{
								"type": "text",
								"text": "%s"
						}
				]
			},
			{
				"type": "button",
				"sub_type": "url",
				"index": 0,
				"parameters": [
						{
						"type": "text",
						"text": "%s"
						}
				]
			}
		]
	`, pinCode, pinCode)
	template := "tschwaa_otp"

	return sendMessageTextFromTemplate(to, template, language, parameters)
}

func getMemberName(member models.Member, language string) string {
	if len(member.FirstName) > 0 || len(member.LastName) > 0 {
		return fmt.Sprintf("%s %s", member.FirstName, member.LastName)
	}

	return "Membre"
}

func SendInvitationToJoinOrganization(member models.Member, organizationName, joinId, organizationReps string) (*WhatsappSendMessageResponse, error) {
	log.Println("SendInvitationToJoinOrganization ", member)
	linkToJoin := fmt.Sprintf("https://tschwaa.com/join/%s", joinId)

	language := "fr"
	template := "tschwaa_invite_member_to_join"
	parameters := fmt.Sprintf(`[
		{
			"type": "body",
			"parameters": [
				{
					"type": "text",
					"text": "%s"
				},
				{
					"type": "text",
					"text": "%s"
				},
				{
					"type": "text",
					"text": "%s"
				},
				{
					"type": "text",
					"text": "%s"
				}
			]
		}
	]`, getMemberName(member, language), organizationName, linkToJoin, organizationReps)
	return sendMessageTextFromTemplate(member.Phone, template, language, parameters)
}
