package entity_test

import (
	"testing"

	"github.com/coma/coma/src/domain/entity"
	"github.com/stretchr/testify/assert"
)

func TestLocalUserAuthToken(t *testing.T) {
	t.Run("full token", func(t *testing.T) {
		Key := "12345"

		localUserAuthToken := entity.LocalUserAuthToken{
			Id:       "1",
			Type:     entity.AccessToken,
			UserType: "admin",
		}

		token, err := localUserAuthToken.GenerateToken(Key)

		assert.NoError(t, err)
		assert.NotEqual(t, "", token)

		parseLocalUserAuthToken, err := entity.NewLocalUserAuthTokenFromToken(token, Key)
		assert.NoError(t, err)
		assert.Equal(t, localUserAuthToken, parseLocalUserAuthToken)
	})

	t.Run("empty token", func(t *testing.T) {
		Key := "12345"
		_, err := entity.NewLocalUserAuthTokenFromToken("", Key)
		assert.Error(t, err)
	})
}
