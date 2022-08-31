package http

import (
	"context"
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"

	"github.com/Scalingo/go-scalingo/v5/http/tokensservicemock"
)

func TestAPITokenGenerator_GetAccessToken(t *testing.T) {
	ctx := context.Background()
	apiToken := "tk-token-test"

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	cases := map[string]struct {
		tokenTTL time.Duration
		expect   func(t *testing.T, s *tokensservicemock.MockTokensService, token string)
		call     func(t *testing.T, gen *APITokenGenerator, token string)
	}{
		"it should return a JWT queried to the TokenService": {
			tokenTTL: time.Hour,
			expect: func(t *testing.T, s *tokensservicemock.MockTokensService, token string) {
				s.EXPECT().TokenExchange(gomock.Any(), apiToken).Return(token, nil)
			},
			call: func(t *testing.T, gen *APITokenGenerator, token string) {
				accessToken, err := gen.GetAccessToken(ctx)
				require.NoError(t, err)
				require.Equal(t, token, accessToken)
			},
		},
		"it should return twice the same JWT with one call to TokenService if stil valid": {
			tokenTTL: time.Hour,
			expect: func(t *testing.T, s *tokensservicemock.MockTokensService, token string) {
				s.EXPECT().TokenExchange(gomock.Any(), apiToken).Return(token, nil)
			},
			call: func(t *testing.T, gen *APITokenGenerator, token string) {
				accessToken, err := gen.GetAccessToken(ctx)
				require.NoError(t, err)
				require.Equal(t, token, accessToken)

				accessToken, err = gen.GetAccessToken(ctx)
				require.NoError(t, err)
				require.Equal(t, token, accessToken)
			},
		},
		"it should requery a another JWT if the token is less than 5 minutes to expire": {
			tokenTTL: 4 * time.Minute,
			expect: func(t *testing.T, s *tokensservicemock.MockTokensService, token string) {
				s.EXPECT().TokenExchange(gomock.Any(), apiToken).Return(token, nil)

				claims := &jwt.RegisteredClaims{
					ExpiresAt: jwt.NewNumericDate(time.Now().Add(5 * time.Minute)),
				}
				jwtToken := jwt.NewWithClaims(jwt.SigningMethodNone, claims)
				jwt, err := jwtToken.SignedString(jwt.UnsafeAllowNoneSignatureType)
				require.NoError(t, err)

				s.EXPECT().TokenExchange(gomock.Any(), apiToken).Return(jwt, nil)
			},
			call: func(t *testing.T, gen *APITokenGenerator, token string) {
				accessToken, err := gen.GetAccessToken(ctx)
				require.NoError(t, err)
				require.Equal(t, token, accessToken)

				accessToken, err = gen.GetAccessToken(ctx)
				require.NoError(t, err)
				require.NotEmpty(t, accessToken)
				require.NotEqual(t, token, accessToken)
			},
		},
	}

	for title, c := range cases {
		t.Run(title, func(t *testing.T) {
			s := tokensservicemock.NewMockTokensService(ctrl)
			gen := NewAPITokenGenerator(s, apiToken)

			claims := &jwt.RegisteredClaims{
				ExpiresAt: jwt.NewNumericDate(time.Now().Add(c.tokenTTL)),
			}
			jwtToken := jwt.NewWithClaims(jwt.SigningMethodNone, claims)
			jwt, err := jwtToken.SignedString(jwt.UnsafeAllowNoneSignatureType)
			require.NoError(t, err)

			c.expect(t, s, jwt)
			c.call(t, gen, jwt)
		})
	}
}
