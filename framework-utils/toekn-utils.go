package frameworkutils

import (
	"errors"

	frameworkconstants "github.com/geekible-ltd/serviceframework/framework-constants"
	frameworkdto "github.com/geekible-ltd/serviceframework/framework-dto"
	"github.com/gin-gonic/gin"
)

func GetTokenDTO(c *gin.Context) (frameworkdto.TokenDTO, error) {
	tokenRaw, ok := c.Get(frameworkconstants.TokenKey)
	if !ok {
		return frameworkdto.TokenDTO{}, errors.New("token not found")
	}

	tokenDTO, ok := tokenRaw.(frameworkdto.TokenDTO)
	if !ok {
		return frameworkdto.TokenDTO{}, errors.New("invalid token format")
	}

	return tokenDTO, nil
}
