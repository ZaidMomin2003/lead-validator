package main

import (
	"fmt"
	"reflect"
	emailverifier "github.com/AfterShip/email-verifier"
)

func main() {
	v := emailverifier.NewVerifier()
	res, _ := v.Verify("test@example.com")
	val := reflect.ValueOf(*res)
	for i := 0; i < val.NumField(); i++ {
		fmt.Println(val.Type().Field(i).Name)
	}
}
