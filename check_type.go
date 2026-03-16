package main

import (
	"fmt"
	"reflect"
	emailverifier "github.com/AfterShip/email-verifier"
)

func main() {
	v := emailverifier.NewVerifier()
	res, _ := v.Verify("test@example.com")
	f, _ := reflect.TypeOf(*res).FieldByName("Gravatar")
	fmt.Println(f.Type.Kind().String())
}
