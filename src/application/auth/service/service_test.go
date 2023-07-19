package service_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/coma/coma/src/application/auth/dto"
	"github.com/coma/coma/src/application/auth/service"
	"github.com/coma/coma/src/domains/entity"
	"github.com/coma/coma/src/domains/repository/repositoryfakes"
	domainservice "github.com/coma/coma/src/domains/service"

	"github.com/stretchr/testify/assert"
)

func TestValidateToken(t *testing.T) {
	type fakes struct {
		repoReaderFake *repositoryfakes.FakeRepositoryAuthReader
	}

	testCases := []struct {
		name     string
		data     dto.RequestValidateToken
		fakes    func(fake fakes)
		expected func() dto.ResponseValidateKey
		withErr  bool
	}{
		{
			name: "Test1 - ApiKey",
			data: dto.RequestValidateToken{
				Method: dto.Apikey,
				Token:  "12345",
			},
			fakes: func(fake fakes) {
				fake.repoReaderFake.FindTokenByTokenReturns(entity.Apikey{
					Id:   1,
					Name: "service-a",
					Key:  "12345",
				}, nil)
			},
			expected: func() dto.ResponseValidateKey {
				return dto.ResponseValidateKey{
					Valid: true,
				}
			},
		},
		{
			name: "Test2 - ApiKey not found",
			data: dto.RequestValidateToken{
				Method: dto.Apikey,
				Token:  "123456",
			},
			fakes: func(fake fakes) {
				fake.repoReaderFake.FindTokenByTokenReturns(entity.Apikey{}, fmt.Errorf("err: token is not found"))
			},
			expected: func() dto.ResponseValidateKey {
				return dto.ResponseValidateKey{
					Valid: false,
				}
			},
			withErr: true,
		},
		{
			name: "Test3 - Oauth not implemented",
			data: dto.RequestValidateToken{
				Method: dto.Oauth,
				Token:  "12345",
			},
			fakes: func(fake fakes) {
				fake.repoReaderFake.FindTokenByTokenReturns(entity.Apikey{}, fmt.Errorf("err: token is not found"))
			},
			expected: func() dto.ResponseValidateKey {
				return dto.ResponseValidateKey{
					Valid: false,
				}
			},
			withErr: true,
		},
	}

	for _, test := range testCases {
		t.Run(test.name, func(t *testing.T) {
			repoReader := repositoryfakes.FakeRepositoryAuthReader{}
			repoWriter := repositoryfakes.FakeRepositoryAuthWriter{}

			test.fakes(fakes{
				repoReaderFake: &repoReader,
			})

			svc := service.New(service.SetAuthSvc(map[dto.Method]domainservice.AuthServicer{
				dto.Apikey: service.NewApiKey(service.SetApiKeyRepository(&repoReader, repoWriter)),
				dto.Oauth:  service.NewOauth(service.SetOauthRepository(&repoReader, repoWriter)),
			}), service.SetRepository(&repoReader, repoWriter))

			res, err := svc.ValidateToken(context.Background(), test.data)

			resExp := test.expected()

			if test.withErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}

			assert.Equal(t, resExp, res)
		})
	}
}
