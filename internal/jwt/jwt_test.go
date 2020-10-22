package jwt_test

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/volatiletech/null"
	"github.com/wednesday-solutions/go-template/internal/jwt"
	"github.com/wednesday-solutions/go-template/models"
)

func TestNew(t *testing.T) {
	cases := map[string]struct {
		algo         string
		secret       string
		minSecretLen int
		req          models.User
		wantErr      bool
		want         jwt.Service
	}{
		"invalid algo": {
			algo:    "invalid",
			wantErr: true,
		},
		"invalid secret length": {
			algo:    "HS256",
			secret:  "123",
			wantErr: true,
		},
		"invalid secret length with min defined": {
			algo:         "HS256",
			minSecretLen: 4,
			secret:       "123",
			wantErr:      true,
		},
		"success": {
			algo:         "HS256",
			secret:       "g0r$kt3$t1ng",
			minSecretLen: 1,
			req: models.User{
				Username: null.StringFrom("johndoe"),
				Email:    null.StringFrom("johndoe@mail.com"),
			},
			want: jwt.Service{},
		},
	}

	for name, tt := range cases {
		t.Run(name, func(t *testing.T) {
			_, err := jwt.New(tt.algo, tt.secret, 60, tt.minSecretLen)
			assert.Equal(t, tt.wantErr, err != nil)
		})
	}
}

func TestGenerateToken(t *testing.T) {
	cases := map[string]struct {
		algo         string
		secret       string
		minSecretLen int
		req          models.User
		wantErr      bool
		want         string
	}{
		"invalid algo": {
			algo:    "invalid",
			wantErr: true,
		},
		"secret not set": {
			algo:    "HS256",
			wantErr: true,
		},
		"invalid secret length": {
			algo:    "HS256",
			secret:  "123",
			wantErr: true,
		},
		"invalid secret length with min defined": {
			algo:         "HS256",
			minSecretLen: 4,
			secret:       "123",
			wantErr:      true,
		},
		"success": {
			algo:         "HS256",
			secret:       "g0r$kt3$t1ng",
			minSecretLen: 1,
			req: models.User{
				Username: null.StringFrom("johndoe"),
				Email:    null.StringFrom("johndoe@mail.com"),
			},
			want: "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9",
		},
	}

	for name, tt := range cases {
		t.Run(name, func(t *testing.T) {
			jwtSvc, err := jwt.New(tt.algo, tt.secret, 60, tt.minSecretLen)
			assert.Equal(t, tt.wantErr, err != nil)
			if err == nil && !tt.wantErr {
				token, _ := jwtSvc.GenerateToken(&tt.req)
				assert.Equal(t, tt.want, strings.Split(token, ".")[0])
			}
		})
	}
}
