package validate_test

import (
	"gopkg.in/orivil/validate.v0"
	"net/url"
	"log"
	"fmt"
)

const (
	// the input names
	username = "username"
	email = "email"
	password = "password"
)

func ExampleValidate() {

	// valid "url.Values" data
	values, err := url.ParseQuery(
		"username=zhangsan&" +
		"email=example@xmail.com",
	)
	if err != nil {
		log.Fatal(err)
	}


	// Setp 1. new validator
	var validator = validate.Validate{
		Required: map[string]string{
			username: "please input user name!",
			email: "please input email!",
			password: "please input password!",
		},
	}

	// Setp 2. check if has any errors, this step should be done in "init()" function
	validator.MustCheck()

	// Setp 3. valid
	if msg := validator.Valid(values); msg != "" {

		fmt.Println(msg)
	}

	// Output:
	// please input password!
}
