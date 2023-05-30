package service_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/coma/coma/src/domains/auth/dto"
	"github.com/coma/coma/src/domains/auth/model"
	"github.com/coma/coma/src/domains/auth/repository/repositoryfakes"
	"github.com/coma/coma/src/domains/auth/service"
	"github.com/stretchr/testify/assert"
)

func TestValidateToken(t *testing.T) {
	type fakes struct {
		repoReaderFake *repositoryfakes.FakeRepositoryReader
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
				fake.repoReaderFake.FindTokenByTokenReturns(model.Apikey{
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
				fake.repoReaderFake.FindTokenByTokenReturns(model.Apikey{}, fmt.Errorf("err: token is not found"))
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
				fake.repoReaderFake.FindTokenByTokenReturns(model.Apikey{}, fmt.Errorf("err: token is not found"))
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
			repo := repositoryfakes.FakeRepositoryReader{}
			test.fakes(fakes{
				repoReaderFake: &repo,
			})

			svc := service.New(&repo, map[dto.Method]service.AuthServicer{
				dto.Apikey: service.NewApiKey(&repo),
				dto.Oauth:  service.NewOauth(&repo),
			})

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
