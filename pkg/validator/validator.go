package validator

import "github.com/asaskevich/govalidator"

type ValidationErrors map[string][]string

func (ve ValidationErrors) Error() string { return "" }

func Validate(dest interface{}) error {
	var errors = ValidationErrors{}

	_, err := govalidator.ValidateStruct(dest)
	if err != nil {
		errs := govalidator.ErrorsByField(err)
		for field, message := range errs {
			errors[field] = append(errors[field], message)
		}

		return errors
	}

	return nil
}
