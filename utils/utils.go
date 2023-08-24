package utils

import (
	"fmt"
	"log"
)

func Fail(msg, customError string, err error) error {
	if err != nil {
		log.Printf("\n%s: %s", msg, err)
		return fmt.Errorf(customError)
	}

	return nil
}
