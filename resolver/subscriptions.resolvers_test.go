package resolver_test

import (
	"context"
	"fmt"
	"os"
	"testing"

	"github.com/joho/godotenv"
	"github.com/stretchr/testify/assert"
	"github.com/volatiletech/null"
	fm "github.com/wednesday-solutions/go-template/graphql_models"
	"github.com/wednesday-solutions/go-template/models"
	"github.com/wednesday-solutions/go-template/resolver"
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
			err := godotenv.Load(fmt.Sprintf("../.env.%s", os.Getenv("ENVIRONMENT_NAME")))
			if err != nil {
				fmt.Print("Error loading .env file")
			}

			c := context.Background()
			ctx := context.WithValue(c, userKey, models.User{ID: 1, FirstName: null.StringFrom("First"), LastName: null.StringFrom("Last"), Username: null.StringFrom("username"), Email: null.StringFrom("mac@wednesday.is"), Mobile: null.StringFrom("+911234567890"), Phone: null.StringFrom("05943-1123"), Address: null.StringFrom("22 Jump Street")})
			response, err := resolver1.Subscription().UserNotification(ctx)
			if tt.wantResp != nil {
				tt.wantResp = response
				assert.Equal(t, tt.wantResp, response)
			}
			assert.Equal(t, tt.wantErr, err != nil)
		})
	}
}
