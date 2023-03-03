package resolver_test

import (
	"context"
	"fmt"
	"testing"

	fm "go-template/gqlmodels"
	"go-template/internal/config"
	"go-template/resolver"
	"go-template/testutls"

	"github.com/stretchr/testify/assert"
)

func TestUserNotification(
	t *testing.T,
) {
	cases := []struct {
		name     string
		wantResp <-chan *fm.User
		wantErr  bool
	}{
		{
			name: SuccessCase,
			wantResp: make(
				chan *fm.User,
				1,
			),
			wantErr: false,
		},
	}

	observers := map[string]chan *fm.User{}
	resolver1 := resolver.Resolver{
		Observers: observers,
	}
	for _, tt := range cases {
		t.Run(
			tt.name,
			func(t *testing.T) {
				err := config.LoadEnv()
				if err != nil {
					fmt.Print("error loading .env file")
				}

				c := context.Background()
				ctx := context.WithValue(c, testutls.UserKey, testutls.MockUser())
				response, err := resolver1.Subscription().UserNotification(ctx)
				if tt.wantResp != nil {
					tt.wantResp = response
					assert.Equal(t, tt.wantResp, response)
				}
				assert.Equal(t, tt.wantErr, err != nil)
			},
		)
	}
}
