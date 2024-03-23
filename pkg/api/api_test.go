// Package api Go-Template
package api

import (
	"context"
	"database/sql"
	"fmt"
	"io"
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
	. "github.com/agiledragon/gomonkey/v2"
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
type testCases struct {
	name                         string
	args                         args
	wantErr                      bool
	setDbCalled                  bool
	getTransportCalled           bool
	postTransportCalled          bool
	optionsTransportCalled       bool
	multipartFormTransportCalled bool
	websocketTransportCalled     bool
}

func TestStart(t *testing.T) {
	tests := initializeTestCases()

	patches := applyPatches()
	defer patches.Reset()

	mockEchoServer, _ := prepareMocks()

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			setupMockFunctions()

			if hasTransportCalls(tt) {
				checkTransportCalls(tt)
			} else {
				checkGraphQLServer(t, tt, mockEchoServer)
			}
		})
	}
}

func initializeTestCases() []testCases {
	return []testCases{
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
			wantErr:                      false,
			getTransportCalled:           false,
			postTransportCalled:          false,
			optionsTransportCalled:       false,
			multipartFormTransportCalled: false,
			websocketTransportCalled:     false,
		},
	}
}

func applyPatches() *gomonkey.Patches {
	return ApplyFunc(os.Getenv, mockGetenv).
		ApplyFunc(sql.Open, mockSqlOpen).
		ApplyFunc(server.Start, mockServerStart).
		ApplyFunc(server.New, mockServerNew)
}

func mockGetenv(key string) string {
	if key == "JWT_SECRET" {
		return testutls.MockJWTSecret
	}
	return ""
}

func mockSqlOpen(driverName string, dataSourceName string) (*sql.DB, error) {
	fmt.Print("sql.Open called\n")
	return nil, nil
}

func mockServerStart(e *echo.Echo, cfg *server.Config) {
	fmt.Print("Fake server started\n")
}

func mockServerNew() *echo.Echo {
	return echo.New()
}

func prepareMocks() (*echo.Echo, *handler.Server) {
	e := echo.New()
	graphqlHandler := handler.New(graphql.NewExecutableSchema(graphql.Config{
		Resolvers: &resolver.Resolver{Observers: make(map[string]chan *graphql.User)},
	}))
	return e, graphqlHandler
}

func setupMockFunctions() {
	ApplyFunc(boil.SetDB, mockSetDB)
}

func mockSetDB(db boil.Executor, tt testCases) {
	fmt.Print("boil.SetDB called\n")
	tt.setDbCalled = true
}

func hasTransportCalls(tt testCases) bool {
	return tt.getTransportCalled || tt.postTransportCalled ||
		tt.optionsTransportCalled || tt.multipartFormTransportCalled
}

func checkTransportCalls(tt testCases) {
	observers := map[string]chan *graphql.User{}
	graphqlHandler := handler.New(graphql.NewExecutableSchema(graphql.Config{
		Resolvers: &resolver.Resolver{Observers: observers},
	}))
	if tt.getTransportCalled || tt.postTransportCalled ||
		tt.optionsTransportCalled || tt.multipartFormTransportCalled {
		ApplyMethod(reflect.TypeOf(graphqlHandler), "AddTransport", func(s *handler.Server, t graphql2.Transport) {
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
	}
}
func checkGraphQLServer(t *testing.T, tt testCases, e *echo.Echo) {
	jsonRes, err := testutls.MakeRequest(testutls.RequestParameters{
		E:           e,
		Pathname:    graphQLPathname,
		HttpMethod:  "POST",
		RequestBody: testutls.MockWhitelistedQuery,
		IsGraphQL:   false,
	})
	if err != nil {
		t.Fatalf("Failed to make request to GraphQL server: %v", err)
	}

	// Assert that database was set
	assert.True(t, tt.setDbCalled, "Expected database to be set")

	// Assert GraphQL schema is returned
	assert.NotNil(t, jsonRes["data"].(map[string]interface{})["__schema"],
		"Expected GraphQL schema to be returned")

	// Simulate request to the GraphQL playground
	_, res, err := testutls.SimpleMakeRequest(testutls.RequestParameters{
		E:          e,
		Pathname:   "/playground",
		HttpMethod: "GET",
		IsGraphQL:  false,
	})
	if err != nil {
		t.Fatalf("Failed to make request to GraphQL playground: %v", err)
	}

	bodyBytes, err := io.ReadAll(res.Body)
	if err != nil {
		t.Fatalf("Failed to read response body: %v", err)
	}

	// Assert that GraphQL playground is returned
	assert.Contains(t, string(bodyBytes), "GraphiQL.createFetcher",
		"Expected GraphQL playground to be returned")

	// Connect to the WebSocket endpoint
	ts := httptest.NewServer(e)
	defer ts.Close()
	u := "ws" + strings.TrimPrefix(ts.URL+graphQLPathname, "http")
	ws, _, err := websocket.DefaultDialer.Dial(u, nil)
	if err != nil {
		t.Fatalf("Failed to connect to WebSocket: %v", err)
	}
	defer ws.Close()

	// Send connection initiation message
	if err := ws.WriteMessage(websocket.TextMessage, []byte(`{"type":"connection_init","payload":`+
		`{"authorization":"bearer ABC"}}`)); err != nil {
		t.Fatalf("Failed to send connection initiation message: %v", err)
	}

	// Read response from WebSocket
	_, p, err := ws.ReadMessage()
	if err != nil {
		t.Fatalf("Failed to read WebSocket message: %v", err)
	}

	// Assert that connection acknowledgement is received
	assert.Contains(t, string(p), "{\"type\":\"connection_ack\"}\n",
		"Expected connection acknowledgement from WebSocket")
}
