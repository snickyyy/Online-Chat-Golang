package validators

import (
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
)

func InitValidators() {
	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		if err := v.RegisterValidation("username", ValidateUsername); err != nil {
			panic(err)
		}
		if err := v.RegisterValidation("password", ValidatePassword); err != nil {
			panic(err)
		}
	}
}
