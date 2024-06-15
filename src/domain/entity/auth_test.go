package entity_test

import (
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/nurcahyaari/coma/src/domain/entity"
	"github.com/stretchr/testify/assert"
)

func TestLocalUserAuthToken(t *testing.T) {
	t.Run("full token", func(t *testing.T) {
		// Key := "12345"

		now := time.Now()
		localUserAuthToken := entity.LocalUserAuthToken{
			RegisteredClaims: jwt.RegisteredClaims{
				ExpiresAt: jwt.NewNumericDate(now.Add(1 * time.Hour)),
				IssuedAt:  jwt.NewNumericDate(now),
			},
			Id:       "1",
			Type:     entity.AccessToken,
			UserType: "admin",
		}

		token, err := localUserAuthToken.GenerateJWTToken(nil)

		assert.NoError(t, err)
		assert.NotEqual(t, "", token)

		parseLocalUserAuthToken, err := entity.NewLocalUserAuthTokenFromToken(token, nil)
		assert.NoError(t, err)
		assert.Equal(t, localUserAuthToken, parseLocalUserAuthToken)
	})

	t.Run("empty token", func(t *testing.T) {
		// Key := "12345"
		_, err := entity.NewLocalUserAuthTokenFromToken("", nil)
		assert.Error(t, err)
	})
}
