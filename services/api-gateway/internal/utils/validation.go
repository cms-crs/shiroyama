package utils

import (
	"regexp"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

var (
	hexColorRegex = regexp.MustCompile(`^#([A-Fa-f0-9]{6}|[A-Fa-f0-9]{3})$`)
)

func BindAndValidate(c *gin.Context, obj interface{}) error {
	if err := c.ShouldBindJSON(obj); err != nil {
		return err
	}
	return nil
}

func BindQueryAndValidate(c *gin.Context, obj interface{}) error {
	if err := c.ShouldBindQuery(obj); err != nil {
		return err
	}
	return nil
}

func BindURIAndValidate(c *gin.Context, obj interface{}) error {
	if err := c.ShouldBindUri(obj); err != nil {
		return err
	}
	return nil
}

func RegisterCustomValidators(v *validator.Validate) {
	v.RegisterValidation("hexcolor", validateHexColor)
}

func validateHexColor(fl validator.FieldLevel) bool {
	color := fl.Field().String()
	return hexColorRegex.MatchString(color)
}

func GetParamID(c *gin.Context, param string) string {
	return c.Param(param)
}

func GetUserIDFromContext(c *gin.Context) (string, bool) {
	userID, exists := c.Get("user_id")
	if !exists {
		return "", false
	}
	return userID.(string), true
}
