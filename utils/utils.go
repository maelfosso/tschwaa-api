package utils

import (
	"fmt"
	"log"
	"reflect"
)

func Fail(msg, customError string, err error) error {
	if err != nil {
		log.Printf("\n%s: %s", msg, err)
		return fmt.Errorf(customError)
	}

	return nil
}

func CheckNilInterface(i interface{}) bool {
	iv := reflect.ValueOf(i)
	if !iv.IsValid() {
		return true
	}
	switch iv.Kind() {
	case reflect.Ptr, reflect.Slice, reflect.Map, reflect.Func, reflect.Interface:
		return iv.IsNil()
	default:
		return false
	}
}
