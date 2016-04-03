package another_example_test

import (
	"regexp"
	"github.com/orivil/validate"
	"net/url"
	"fmt"
)

const (
// the input names
	username = "username"
	email = "email"
	password = "password"
	confirm = "confirmPassword"
	age = "age"
)

var passwordPatten = regexp.MustCompile(`^[\w]+$`)

// new user login validator
var loginValidator = validate.Validate{

	Required: map[string]string{
		email: "please input email!",
		password: "please input password!",
	},

	Email: map[string]string{
		email: "email format incorrect！",
	},

	// field "SliceRange" and "NumRange" is just like "StringRange"
	StringRange: map[string]map[string]string{
		// "|6|16|" means "6 <= len <= 16"
		// "|6|16" means "6 <= len < 16"
		// "6|16" means "6 < len < 16"
		password: {"|6|16|": "password must have 6-16 characters"},
	},

	Regexp: map[string]map[string]*regexp.Regexp{
		password: {"The password must be alphanumeric characters and underscores": passwordPatten},
	},
}

// new user register validator
var registerValidator = validate.Validate{

	Required: map[string]string{
		username: "please input user name!",
		confirm: "please input confirm password!",
	},

	StringRange: map[string]map[string]string{
		username: {"|4|16|": "Username must have 4-16 characters"},
	},

	// field "Max" is just like "Min"
	Min: map[string]map[string]string{
		// "18" means "num > 18"
		// "|18|" means "num >= 18"
		age: {"|18|": "age must be bigger than 18"},
	},

	// Confirm means "equal"
	Confirm: map[string]map[string]string{
		password: {confirm: "confirm password failed!"},
	},
}

func init() {

	registerValidator.MustCheck()
	loginValidator.MustCheck()

	// add a "loginValidator" to the "registerValidator", so the "registerValidator"
	// do not have to valid the options which the "loginValidator" already valid
	registerValidator.AddValidator(loginValidator)
}

func ExampleValidate() {
	values, err := url.ParseQuery(
		"username=zhangsan&" +
		"email=example@xmail.com&" +
		"password=123456",
	)
	if err != nil {
		panic(err)
	}

	if msg := registerValidator.Valid(values); msg != "" {
		fmt.Println(msg)
	}
}