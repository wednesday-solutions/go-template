package jwt_test

import (
	"log"
	"strings"
	"testing"

	jwtgo "github.com/dgrijalva/jwt-go"
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
		error        string
	}{
		"invalid algo": {
			algo:         "invalid",
			wantErr:      true,
			minSecretLen: 1,
			secret:       "g0r$kt3$t1ng",
			error:        "invalid jwt signing method: invalid",
		},
		"invalid secret length": {
			algo:    "HS256",
			secret:  "123",
			wantErr: true,
			error:   "jwt secret length is 3, which is less than required 128",
		},
		"invalid secret length with min defined": {
			algo:         "HS256",
			minSecretLen: 4,
			secret:       "123",
			wantErr:      true,
			error:        "jwt secret length is 3, which is less than required 4",
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
			if err != nil {
				assert.Equal(t, tt.error, err.Error())
			}
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

func TestParseToken(t *testing.T) {
	algo := "HS256"
	cases := map[string]struct {
		authHeader string
		error      string
		algo       string
	}{
		"Failure_InvalidToken": {
			authHeader: "bearer 123",
			error:      "token contains an invalid number of segments",
			algo:       algo,
		},
		"Failure_NoAuth": {
			authHeader: "",
			error:      "generic error",
			algo:       algo,
		},
		"Failure_MismatchTokenMethod": {
			authHeader: "bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjM0NTY3ODkwIi" +
				"wibmFtZSI6IkpvaG4gRG9lIiwiaWF0IjoxNTE2MjM5MDIyfQ.SflKxwRJSMeKKF2QT4fwpMeJf36POk6yJV_adQssw5c",
			algo:  "ES256",
			error: "generic error",
		},
		"Success": {
			authHeader: "bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjM0NTY3ODkwIi" +
				"wibmFtZSI6IkpvaG4gRG9lIiwiaWF0IjoxNTE2MjM5MDIyfQ.SflKxwRJSMeKKF2QT4fwpMeJf36POk6yJV_adQssw5c",
			algo: algo,
		},
	}

	for name, tt := range cases {
		t.Run(name, func(t *testing.T) {
			jwtSvc, err := jwt.New(tt.algo, "g0r$kt3$t1ng", 60, 1)
			if err != nil {
				log.Fatal(err)
			}
			token, err := jwtSvc.ParseToken(tt.authHeader)
			if len(tt.error) != 0 {
				assert.Equal(t, tt.error, err.Error())
			} else {
				assert.Equal(t, "John Doe", token.Claims.(jwtgo.MapClaims)["name"])
			}
		})
	}

}
