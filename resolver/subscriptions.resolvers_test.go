package resolver_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/joho/godotenv"
	"github.com/stretchr/testify/assert"
	fm "github.com/wednesday-solutions/go-template/graphql_models"
	"github.com/wednesday-solutions/go-template/resolver"
	"github.com/wednesday-solutions/go-template/testutls"
)

func TestUserNotification(t *testing.T) {
	cases := []struct {
		name     string
		wantResp <-chan *fm.User
		wantErr  bool
	}{
		{
			name:     "Success",
			wantResp: make(chan *fm.User, 1),
			wantErr:  false,
		},
	}

	observers := map[string]chan *fm.User{}
	resolver1 := resolver.Resolver{Observers: observers}
	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			err := godotenv.Load("../.env.local")
			if err != nil {
				fmt.Print("error loading .env file")
			}

			c := context.Background()
			ctx := context.
				WithValue(c, testutls.UserKey, testutls.MockUser())
			response, err := resolver1.Subscription().UserNotification(ctx)
			if tt.wantResp != nil {
				tt.wantResp = response
				assert.Equal(t, tt.wantResp, response)
			}
			assert.Equal(t, tt.wantErr, err != nil)
		})
	}
}
