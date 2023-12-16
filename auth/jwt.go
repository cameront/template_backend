package auth

import (
	"fmt"
	"time"

	"github.com/cristalhq/jwt/v5"
)

var key = "set_me!"

func InitJWT(secretKey string) {
	key = secretKey
}

func buildToken(id, name, email, role string, expires time.Time) (string, error) {
	key := []byte(key)
	signer, err := jwt.NewSignerHS(jwt.HS256, key)
	if err != nil {
		return "", fmt.Errorf("error creating signer: %v", err)
	}

	builder := jwt.NewBuilder(signer)

	claims := &UserClaims{
		RegisteredClaims: jwt.RegisteredClaims{
			Audience:  []string{role},
			Subject:   email,
			ID:        id,
			ExpiresAt: jwt.NewNumericDate(expires),
		},
		Name:  name,
		Email: email,
	}

	token, err := builder.Build(claims)
	if err != nil {
		return "", fmt.Errorf("error building token claims: %v", err)
	}

	return token.String(), nil
}

func validateToken(tokenStr string) (*UserClaims, error) {
	verifier, err := jwt.NewVerifierHS(jwt.HS256, []byte(key))
	if err != nil {
		return nil, fmt.Errorf("error building verifier: %v", err)
	}

	var claims UserClaims
	if err = jwt.ParseClaims([]byte(tokenStr), verifier, &claims); err != nil {
		return nil, fmt.Errorf("error verifying claims: %v", err)
	}

	// TODO: consumers need to validate that token is not expired
	return &claims, nil
}
