package jwt

import (
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v4"
	portainer "github.com/portainer/portainer/api"
	i "github.com/portainer/portainer/api/internal/testhelpers"
	"github.com/stretchr/testify/assert"
)

func TestGenerateSignedToken(t *testing.T) {
	dataStore := i.NewDatastore(i.WithSettingsService(&portainer.Settings{}))
	svc, err := NewService("24h", dataStore)
	assert.NoError(t, err, "failed to create a copy of service")

	token := &portainer.TokenData{
		Username: "Joe",
		ID:       1,
		Role:     1,
	}
	expiresAt := time.Now().Add(1 * time.Hour)

	generatedToken, err := svc.generateSignedToken(token, expiresAt, defaultScope)
	assert.NoError(t, err, "failed to generate a signed token")

	parsedToken, err := jwt.ParseWithClaims(generatedToken, &claims{}, func(token *jwt.Token) (any, error) {
		return svc.secrets[defaultScope], nil
	})
	assert.NoError(t, err, "failed to parse generated token")

	tokenClaims, ok := parsedToken.Claims.(*claims)
	assert.Equal(t, true, ok, "failed to claims out of generated ticket")

	assert.Equal(t, token.Username, tokenClaims.Username)
	assert.Equal(t, int(token.ID), tokenClaims.UserID)
	assert.Equal(t, int(token.Role), tokenClaims.Role)
	assert.Equal(t, jwt.NewNumericDate(expiresAt), tokenClaims.ExpiresAt)
}

func TestGenerateSignedToken_InvalidScope(t *testing.T) {
	dataStore := i.NewDatastore(i.WithSettingsService(&portainer.Settings{}))
	svc, err := NewService("24h", dataStore)
	assert.NoError(t, err, "failed to create a copy of service")

	token := &portainer.TokenData{
		Username: "Joe",
		ID:       1,
		Role:     1,
	}
	expiresAt := time.Now().Add(1 * time.Hour)

	_, err = svc.generateSignedToken(token, expiresAt, "testing")
	assert.Error(t, err)
	assert.Equal(t, "invalid scope: testing", err.Error())
}
