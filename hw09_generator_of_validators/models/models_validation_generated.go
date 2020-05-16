/*
* CODE GENERATED AUTOMATICALLY WITH go-validate
* THIS FILE SHOULD NOT BE EDITED BY HAND
 */
//nolint:gomnd,gofmt,goimports
package models

import (
	"fmt"
	"regexp"
)

type ValidationError struct {
	Field string
	Err   error
}

func (user User) Validate() ([]ValidationError, error) {
	ve := []ValidationError{}

	//ID
	if len(user.ID) != 36 {
		ve = append(ve, ValidationError{
			Field: "ID",
			Err:   fmt.Errorf("len field ID with value:%v not equal with validate value %v", user.ID, 36),
		})
	}

	//Age
	if user.Age < 18 {
		ve = append(ve, ValidationError{
			Field: "Age",
			Err:   fmt.Errorf("field Age with value:%v smaller than min value %v", user.Age, 18),
		})
	}

	if user.Age > 50 {
		ve = append(ve, ValidationError{
			Field: "Age",
			Err:   fmt.Errorf("field Age with value:%v bigger than max value %v", user.Age, 50),
		})
	}

	//Email
	r := regexp.MustCompile(`^\w+@\w+\.\w+$`)
	if !r.MatchString(user.Email) {
		ve = append(ve, ValidationError{
			Field: "Email",
			Err:   fmt.Errorf("field Email with value:%v not match the regex %v", user.Email, `^\w+@\w+\.\w+$`),
		})
	}

	//Role
	if user.Role != "admin" && user.Role != "stuff" {
		ve = append(ve, ValidationError{
			Field: "Role",
			Err:   fmt.Errorf("field Role with value:%v is not one of  the values %v", user.Role, "[admin stuff]"),
		})
	}

	//Phones
	for _, s := range user.Phones {
		if len(s) != 11 {
			ve = append(ve, ValidationError{
				Field: "Phones",
				Err:   fmt.Errorf("len element of Phones with value:%v not equal with validate value %v", s, 11),
			})
		}
	}
	return ve, nil
}

func (app App) Validate() ([]ValidationError, error) {
	ve := []ValidationError{}

	//Version
	if len(app.Version) != 5 {
		ve = append(ve, ValidationError{
			Field: "Version",
			Err:   fmt.Errorf("len field Version with value:%v not equal with validate value %v", app.Version, 5),
		})
	}
	return ve, nil
}

func (response Response) Validate() ([]ValidationError, error) {
	ve := []ValidationError{}

	//Code
	if response.Code != 200 && response.Code != 404 && response.Code != 500 {
		ve = append(ve, ValidationError{
			Field: "Code",
			Err:   fmt.Errorf("field Code with value:%v is not one of  the values %v", response.Code, "[200 404 500]"),
		})
	}
	return ve, nil
}
