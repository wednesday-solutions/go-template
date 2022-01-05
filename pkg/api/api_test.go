// Package api Go-Template
package api

import (
	"context"
	"database/sql"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"reflect"
	"testing"
	"time"

	graphql2 "github.com/99designs/gqlgen/graphql"
	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/handler/transport"
	. "github.com/agiledragon/gomonkey/v2"
	"github.com/gorilla/websocket"
	"github.com/labstack/echo"
	_ "github.com/lib/pq"
	"github.com/stretchr/testify/assert"
	"github.com/volatiletech/sqlboiler/boil"
	graphql "github.com/wednesday-solutions/go-template/graphql_models"
	"github.com/wednesday-solutions/go-template/internal/config"
	"github.com/wednesday-solutions/go-template/internal/server"
	"github.com/wednesday-solutions/go-template/resolver"
	"github.com/wednesday-solutions/go-template/testutls"
)

func TestStart(t *testing.T) {
	type args struct {
		cfg *config.Configuration
	}
	tests := []struct {
		name        string
		args        args
		wantErr     bool
		setDbCalled bool

		getTransportCalled           bool
		postTransportCalled          bool
		optionsTransportCalled       bool
		multipartFormTransportCalled bool
		websocketTransportCalled     bool
	}{
		{
			name: "Success",
			args: args{
				cfg: testutls.MockConfig(),
			},
			wantErr: false,
		},
		{
			name: "Test_AddTransport",
			args: args{
				cfg: testutls.MockConfig(),
			},
			wantErr: false,

			getTransportCalled:           false,
			postTransportCalled:          false,
			optionsTransportCalled:       false,
			multipartFormTransportCalled: false,
			websocketTransportCalled:     false,
		},
	}

	ApplyFunc(os.Getenv, func(key string) (value string) {
		if key == "JWT_SECRET" {
			return testutls.MockJWTSecret
		}
		return ""
	})
	ApplyFunc(sql.Open, func(driverName string, dataSourceName string) (*sql.DB, error) {
		fmt.Print("sql.Open called\n")
		return nil, nil
	})
	ApplyFunc(server.Start, func(e *echo.Echo, cfg *server.Config) {
		fmt.Print("Fake server started\n")
	})
	e := echo.New()

	ApplyFunc(server.New, func() *echo.Echo {
		return e
	})

	observers := map[string]chan *graphql.User{}
	graphqlHandler := handler.New(graphql.NewExecutableSchema(graphql.Config{
		Resolvers: &resolver.Resolver{Observers: observers},
	}))

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ApplyFunc(boil.SetDB, func(db boil.Executor) {
				fmt.Print("boil.SetDB called\n")
				tt.setDbCalled = true
			})

			if tt.getTransportCalled || tt.postTransportCalled ||
				tt.optionsTransportCalled || tt.multipartFormTransportCalled {

				ApplyMethod(reflect.TypeOf(graphqlHandler), "AddTransport", func(s *handler.Server, t graphql2.Transport) {

					transportGET := transport.GET{}
					transportOptions := transport.Options{}
					transportMultipartForm := transport.MultipartForm{}
					transportPOST := transport.POST{}
					transportWebsocket := transport.Websocket{
						KeepAlivePingInterval: 10 * time.Second,
						InitFunc: func(ctx context.Context, initPayload transport.InitPayload) (context.Context, error) {
							return ctx, nil
						},
						Upgrader: websocket.Upgrader{
							CheckOrigin: func(r *http.Request) bool {
								return true
							},
						},
					}

					if t == transportGET {
						tt.getTransportCalled = true
					}
					if t == transportOptions {
						tt.optionsTransportCalled = true
					}
					if t == transportMultipartForm {
						tt.multipartFormTransportCalled = true
					}
					if t == transportPOST {
						tt.postTransportCalled = true
					}

					if reflect.TypeOf(t) == reflect.TypeOf(transportWebsocket) {
						tt.websocketTransportCalled = true
					}

				})
				_, err := Start(tt.args.cfg)
				if err != nil != tt.wantErr {
					t.Errorf("Start() error = %v, wantErr %v", err, tt.wantErr)
				}

				assert.Equal(t, tt.getTransportCalled, true)
				assert.Equal(t, tt.optionsTransportCalled, true)
				assert.Equal(t, tt.multipartFormTransportCalled, true)
				assert.Equal(t, tt.postTransportCalled, true)
				assert.Equal(t, tt.websocketTransportCalled, true)

			} else {
				_, err := Start(tt.args.cfg)
				if err != nil != tt.wantErr {
					t.Errorf("Start() error = %v, wantErr %v", err, tt.wantErr)
				}
				jsonRes, err := testutls.MakeRequest(testutls.RequestParameters{
					E:           e,
					Pathname:    "/graphql",
					HttpMethod:  "POST",
					RequestBody: testutls.MockWhitelistedQuery,
					IsGraphQL:   false,
				})

				if err != nil {
					log.Fatal(err)
				}

				assert.Equal(t, tt.setDbCalled, true)

				// check if it returns schema correctly
				assert.NotNil(t, jsonRes["data"].(map[string]interface{})["__schema"])

				_, res, err := testutls.SimpleMakeRequest(testutls.RequestParameters{
					E:          e,
					Pathname:   "/playground",
					HttpMethod: "GET",

					IsGraphQL: false,
				})
				if err != nil {
					log.Fatal(err)
				}
				bodyBytes, _ := io.ReadAll(res.Body)

				// check if the playground is returned
				assert.Contains(t, string(bodyBytes), "GraphQLPlayground.init")
			}

		})
	}
}
