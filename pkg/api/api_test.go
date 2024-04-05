// Package api Go-Template
package api

import (
	"context"
	"database/sql"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"reflect"
	"strings"
	"testing"
	"time"

	graphql "go-template/gqlmodels"
	"go-template/internal/config"
	"go-template/internal/server"
	"go-template/resolver"
	"go-template/testutls"

	graphql2 "github.com/99designs/gqlgen/graphql"
	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/handler/transport"
	"github.com/agiledragon/gomonkey/v2"
	"github.com/gorilla/websocket"
	"github.com/labstack/echo/v4"
	_ "github.com/lib/pq"
	"github.com/stretchr/testify/assert"
	"github.com/volatiletech/sqlboiler/v4/boil"
)

var graphQLPathname = "/graphql"

type args struct {
	cfg *config.Configuration
}
type testStartServerType struct {
	name                         string
	args                         args
	want                         bool
	wantErr                      bool
	setDbCalled                  bool
	getTransportCalled           bool
	postTransportCalled          bool
	optionsTransportCalled       bool
	multipartFormTransportCalled bool
	websocketTransportCalled     bool
	init                         func(e *echo.Echo, tt testStartServerType) *gomonkey.Patches
}

func testStartServerSuccessCase() testStartServerType {
	return testStartServerType{
		name: "Success",
		args: args{
			cfg: testutls.MockConfig(),
		},
		setDbCalled: true,
		wantErr:     false,
		init: func(e *echo.Echo, tt testStartServerType) *gomonkey.Patches {
			return gomonkey.ApplyFunc(os.Getenv, func(key string) (value string) {
				if key == "JWT_SECRET" {
					return testutls.MockJWTSecret
				}
				return ""
			}).ApplyFunc(sql.Open, func(driverName string, dataSourceName string) (*sql.DB, error) {
				fmt.Print("sql.Open called\n")
				return nil, nil
			}).ApplyFunc(server.Start, func(e *echo.Echo, cfg *server.Config) {
				fmt.Print("Fake server started\n")
			}).ApplyFunc(server.New, func() *echo.Echo {
				return e
			}).ApplyFunc(boil.SetDB, func(db boil.Executor) {
				fmt.Print("boil.SetDB called\n", tt.setDbCalled)
			})
		},
	}
}
func testWithTransportCase() testStartServerType {
	return testStartServerType{
		name: "Test_AddTransport",
		args: args{
			cfg: testutls.MockConfig(),
		},
		wantErr:                      false,
		want:                         true,
		getTransportCalled:           true,
		postTransportCalled:          true,
		optionsTransportCalled:       true,
		multipartFormTransportCalled: true,
		websocketTransportCalled:     true,
		init: func(e *echo.Echo, tt testStartServerType) *gomonkey.Patches {
			observers := map[string]chan *graphql.User{}
			graphqlHandler := handler.New(graphql.NewExecutableSchema(graphql.Config{
				Resolvers: &resolver.Resolver{Observers: observers},
			}))
			return gomonkey.ApplyFunc(os.Getenv, func(key string) (value string) {
				if key == "JWT_SECRET" {
					return testutls.MockJWTSecret
				}
				return ""
			}).ApplyFunc(sql.Open, func(driverName string, dataSourceName string) (*sql.DB, error) {
				fmt.Print("sql.Open called\n")
				return nil, nil
			}).ApplyFunc(server.Start, func(e *echo.Echo, cfg *server.Config) {
				fmt.Print("Fake server started\n")
			}).ApplyFunc(server.New, func() *echo.Echo {
				return e
			}).ApplyFunc(boil.SetDB, func(db boil.Executor) {
				//mocking  boil.setdb
			}).ApplyMethod(reflect.TypeOf(graphqlHandler), "AddTransport", func(s *handler.Server, t graphql2.Transport) {
				transportGET := transport.GET{}
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
		},
	}
}
func loadTestStartCases() []testStartServerType {
	return []testStartServerType{
		testStartServerSuccessCase(),
		testWithTransportCase(),
	}
}

func checkGraphQLGetResponse(e *echo.Echo) (*http.Response, error) {
	_, res, err := testutls.SimpleMakeRequest(testutls.RequestParameters{
		E:          e,
		Pathname:   "/playground",
		HttpMethod: "GET",
		IsGraphQL:  false,
	})
	return res, err
}
func checkGraphQLPostResponse(e *echo.Echo) (map[string]interface{}, error) {
	res, err := testutls.MakeRequest(testutls.RequestParameters{
		E:           e,
		Pathname:    graphQLPathname,
		HttpMethod:  "POST",
		RequestBody: testutls.MockWhitelistedQuery,
		IsGraphQL:   false,
	})
	return res, err
}
func TestStart(t *testing.T) {
	tests := loadTestStartCases()
	e := echo.New()
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			patches := tt.init(e, tt)
			if tt.getTransportCalled || tt.postTransportCalled ||
				tt.optionsTransportCalled || tt.multipartFormTransportCalled {
				testWithTransportCalls(t, tt)
			} else {
				testWithoutTransportCalls(t, tt, e)
			}
			if patches != nil {
				patches.Reset()
			}
		})
	}
}

func testWithTransportCalls(t *testing.T, tt testStartServerType) {
	_, err := Start(tt.args.cfg)
	if err != nil != tt.wantErr {
		assert.Equal(t, err, tt.wantErr)
	}
	assert.Equal(t, tt.getTransportCalled, tt.want)
	assert.Equal(t, tt.multipartFormTransportCalled, tt.want)
	assert.Equal(t, tt.postTransportCalled, tt.want)
	assert.Equal(t, tt.websocketTransportCalled, tt.want)
}

func testWithoutTransportCalls(t *testing.T, tt testStartServerType, e *echo.Echo) {
	_, err := Start(tt.args.cfg)
	if err != nil != tt.wantErr {
		assert.Equal(t, err, tt.wantErr)
	}
	jsonRes, err := checkGraphQLPostResponse(e)
	if err != nil {
		log.Fatal(err)
	}
	assert.Equal(t, true, jsonRes["data"].(map[string]interface{})["__schema"] != nil, "Schema is nil")
	res, err := checkGraphQLGetResponse(e)
	if err != nil {
		log.Fatal(err)
	}
	bodyBytes, _ := io.ReadAll(res.Body)
	assert.Contains(t, string(bodyBytes), "GraphiQL.createFetcher", "Playground not found")
	ts := httptest.NewServer(e)
	u := "ws" + strings.TrimPrefix(ts.URL+graphQLPathname, "http")
	ws, _, err := websocket.DefaultDialer.Dial(u, nil)
	if err != nil {
		t.Fatalf("%v", err)
	}
	if err := ws.WriteMessage(websocket.TextMessage, []byte(`{"type":"connection_init","payload":`+
		`{"authorization":"bearer ABC"}}`)); err != nil {
		t.Fatalf("%v", err)
	}
	_, p, err := ws.ReadMessage()
	if err != nil {
		t.Fatalf("%v", err)
	}
	assert.Contains(t, string(p), "{\"type\":\"connection_ack\"}\n", "Connection ack not found")
	ts.Close()
	ws.Close()
}
