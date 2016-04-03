package validate

import (
	"net/url"
	"regexp"
	"fmt"
	"strconv"
	"strings"
)

var emailPattern = regexp.MustCompile("[\\w!#$%&'*+/=?^_`{|}~-]+(?:\\.[\\w!#$%&'*+/=?^_`{|}~-]+)*@(?:[\\w](?:[\\w-]*[\\w])?\\.)+[a-zA-Z0-9](?:[\\w-]*[\\w])?")

type Validate struct {
	Required    map[string]string                    // {inputName: msg}
	Email       map[string]string                    // {inputName: msg}
	Confirm     map[string]map[string]string         // {inputName1: {inputName2: msg}}
	SliceRange  map[string]map[string]string         // {inputName: {numRange: msg}}
	StringRange map[string]map[string]string         // {inputName: {lenRange: msg}}
	NumRange    map[string]map[string]string         // {inputName: {numRange: msg}}
	Min         map[string]map[string]string         // {inputName: {minNum: msg}}
	Max         map[string]map[string]string         // {inputName: {maxNum: msg}}
	Regexp      map[string]map[string]*regexp.Regexp // {inputName: {msg: *Regexp}}
	validators  []*Validate
}

// MustCheck for check if there any errors
func (v *Validate) MustCheck() {

	mustGetRange(v.SliceRange)
	mustGetRange(v.StringRange)
	mustGetRange(v.NumRange)

	mustGetNum(v.Min)
	mustGetNum(v.Max)
}

// AddValidator for add another validator to this validator
func (this *Validate) AddValidator(v *Validate) {
	this.validators = append(this.validators, v)
}

// Valid for validate the values, if return "", means valid success
// else return the error messages
func (v *Validate) Valid(vs url.Values) (msg string) {

	for _, validator := range v.validators {
		if msg := validator.Valid(vs); msg != "" {
			return msg
		}
	}

	for key, msg := range v.Required {
		if vs.Get(key) == "" {
			return msg
		}
	}

	for key, msg := range v.Email {
		value := vs.Get(key)
		if value != "" {
			if !emailPattern.MatchString(value) {
				return msg
			}
		}
	}

	for input1, msgs := range v.Confirm {
		for input2, msg := range msgs {
			if vs.Get(input1) != vs.Get(input2) {
				return msg
			}
		}
	}

	for key, msgs := range v.SliceRange {
		values := vs[key]
		if values != nil {
			length := len(values)
			for rang, msg := range msgs {
				small, big, sEqual, bEqual, _ := getRange(rang)
				if sEqual {
					if length < small {
						return msg
					}
				} else {
					if length <= small {
						return msg
					}
				}

				if bEqual {
					if length > big {
						return msg
					}
				} else {
					if length >= big {
						return msg
					}
				}
			}
		}
	}

	for key, msgs := range v.StringRange {
		value := vs.Get(key)
		if value != "" {
			length := len(value)
			for rang, msg := range msgs {
				small, big, sEqual, bEqual, _ := getRange(rang)
				if sEqual {
					if length < small {
						return msg
					}
				} else {
					if length <= small {
						return msg
					}
				}

				if bEqual {
					if length > big {
						return msg
					}
				} else {
					if length >= big {
						return msg
					}
				}
			}
		}
	}

	for key, msgs := range v.NumRange {
		value := vs.Get(key)
		if value != "" {
			num, err := strconv.Atoi(value)
			for rang, msg := range msgs {
				// if user input out of range, just print the message
				if err != nil {
					return msg
				}
				small, big, sEqual, bEqual, _ := getRange(rang)
				if sEqual {
					if num < small {
						return msg
					}
				} else {
					if num <= small {
						return msg
					}
				}

				if bEqual {
					if num > big {
						return msg
					}
				} else {
					if num >= big {
						return msg
					}
				}
			}
		}
	}

	for key, msgs := range v.Min {
		value := vs.Get(key)
		if value != "" {
			input, err := strconv.Atoi(value)
			for numStr, msg := range msgs {
				if err != nil {
					return msg
				}
				num, equal, _ := getNum(numStr)
				if equal {
					if input < num {
						return msg
					}
				} else {
					if input <= num {
						return msg
					}
				}
			}
		}
	}

	for key, msgs := range v.Max {
		value := vs.Get(key)
		if value != "" {
			input, err := strconv.Atoi(value)
			for numStr, msg := range msgs {
				if err != nil {
					return msg
				}
				num, equal, _ := getNum(numStr)
				if equal {
					if input > num {
						return msg
					}
				} else {
					if input >= num {
						return msg
					}
				}
			}
		}
	}

	for key, msgs := range v.Regexp {
		value := vs.Get(key)
		if value != "" {
			for msg, patten := range msgs {
				if !patten.MatchString(value) {
					return msg
				}
			}
		}
	}
	return
}

func getRange(rang string) (s, b int, sEqual, bEqual bool, err error) {
	strs := strings.Split(rang, "|")
	var small, big string
	l := len(strs)
	if l == 2 {
		// "8|18"
		small, big = strs[0], strs[1]
		sEqual = false
		bEqual = false
	} else if l == 4 {
		// "|8|18|"
		small = strs[1]
		big = strs[2]
		sEqual = true
		bEqual = true
	} else if l == 3 {
		if strs[0] == "" {
			// "|8|18"
			small = strs[1]
			big = strs[2]
			sEqual = true
			bEqual = false
		} else {
			// "8|18|"
			small = strs[0]
			big = strs[1]
			sEqual = false
			bEqual = true
		}
	}

	s, err = strconv.Atoi(small)
	b, err = strconv.Atoi(big)
	if s > b {
		s, b = b, s
		sEqual, bEqual = bEqual, sEqual
	}
	return
}

func getNum(num string) (numInt int, equal bool, err error) {
	strs := strings.Split(num, "|")
	var numStr string
	l := len(strs)
	if l == 1 {
		// "8"
		numStr = strs[0]
		equal = false
	} else if l == 3 {
		// "|8|"
		numStr = strs[1]
		equal = true
	}

	numInt, err = strconv.Atoi(numStr)
	return
}

var pattenENum = regexp.MustCompile(`^\|\d+\|$`)
var pattenNum = regexp.MustCompile(`^\d+$`)
var pattenRange = regexp.MustCompile(`^\|?\d+\|\d+\|?$`)

func mustGetRange(data map[string]map[string]string) {
	for _, msgs := range data {
		for rang, _ := range msgs {
			if !pattenRange.MatchString(rang) {
				panic(fmt.Errorf("validate.Validate: range should be like '8|16', '|8|16', '8|16|' or '|8|16|', got: %s\n", rang))
			}
			_, _, _, _, err := getRange(rang)
			if err != nil {
				panic(fmt.Errorf("validate.Validate: %s\n", err.Error()))
			}
		}
	}
}

func mustGetNum(data map[string]map[string]string) {
	for _, msgs := range data {
		for num, _ := range msgs {
			if !pattenNum.MatchString(num) && !pattenENum.MatchString(num) {
				panic(fmt.Errorf("validate.Validate: num should be like '8' or '|8|', got: %s\n", num))
			}

			_, _, err := getNum(num)
			if err != nil {
				panic(fmt.Errorf("validate.Validate: %s\n", err.Error()))
			}
		}
	}
}
