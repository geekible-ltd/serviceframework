package frameworkutils

import (
	"errors"
	"fmt"
	"time"

	frameworkdto "github.com/geekible-ltd/serviceframework/framework-dto"
	"github.com/golang-jwt/jwt/v5"
)

func GenerateJWT(userID, tenantID any, email, firstName, lastName, role string, jwtSecret []byte) (string, error) {
	claims := jwt.MapClaims{
		"sub":        userID,
		"tenant_id":  tenantID,
		"email":      email,
		"first_name": firstName,
		"last_name":  lastName,
		"role":       role,
		"exp":        time.Now().Add(10 * time.Hour).Unix(),
		"iat":        time.Now().Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(jwtSecret)
}

func ParseJWT(tokenString string, jwtSecret []byte) (frameworkdto.TokenDTO, error) {
	tok, err := jwt.Parse(tokenString, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", t.Header["alg"])
		}
		return jwtSecret, nil
	})
	if err != nil {
		return frameworkdto.TokenDTO{}, err
	}
	if !tok.Valid {
		return frameworkdto.TokenDTO{}, errors.New("invalid token")
	}
	claims, ok := tok.Claims.(jwt.MapClaims)
	if !ok {
		return frameworkdto.TokenDTO{}, errors.New("invalid claims")
	}

	// JWT claims store numbers as float64, need to convert properly
	subRaw, ok := claims["sub"]
	if !ok {
		return frameworkdto.TokenDTO{}, errors.New("sub claim missing or invalid")
	}
	var sub string
	switch v := subRaw.(type) {
	case string:
		sub = v
	case float64:
		sub = fmt.Sprintf("%.0f", v)
	default:
		return frameworkdto.TokenDTO{}, errors.New("sub claim invalid type")
	}

	tenantIDRaw, ok := claims["tenant_id"]
	if !ok {
		return frameworkdto.TokenDTO{}, errors.New("tenant_id claim missing or invalid")
	}
	tenantIDFloat, ok := tenantIDRaw.(float64)
	if !ok {
		return frameworkdto.TokenDTO{}, errors.New("tenant_id claim invalid type")
	}
	tenantID := uint(tenantIDFloat)

	email, ok := claims["email"].(string)
	if !ok {
		return frameworkdto.TokenDTO{}, errors.New("email claim missing or invalid")
	}

	firstName, ok := claims["first_name"].(string)
	if !ok {
		return frameworkdto.TokenDTO{}, errors.New("first_name claim missing or invalid")
	}

	lastName, ok := claims["last_name"].(string)
	if !ok {
		return frameworkdto.TokenDTO{}, errors.New("last_name claim missing or invalid")
	}

	roleRaw, ok := claims["role"]
	if !ok {
		return frameworkdto.TokenDTO{}, errors.New("role claim missing or invalid")
	}
	role, ok := roleRaw.(string)
	if !ok {
		return frameworkdto.TokenDTO{}, errors.New("role claim invalid type")
	}

	expFloat, ok := claims["exp"].(float64)
	if !ok {
		return frameworkdto.TokenDTO{}, errors.New("exp claim missing or invalid")
	}
	exp := int64(expFloat)

	iatFloat, ok := claims["iat"].(float64)
	if !ok {
		return frameworkdto.TokenDTO{}, errors.New("iat claim missing or invalid")
	}
	iat := int64(iatFloat)

	return frameworkdto.TokenDTO{
		Sub:       sub,
		TenantID:  tenantID,
		Email:     email,
		FirstName: firstName,
		LastName:  lastName,
		Role:      role,
		Exp:       exp,
		Iat:       iat,
	}, nil
}
