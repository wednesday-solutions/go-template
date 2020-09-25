// Copyright 2017 Emir Ribic. All rights reserved.
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE file.

// Package api Go-Boiler
//
// API Docs for GO-BOILER v1
//
// 	 Terms Of Service:  N/A
//     Schemes: http
//     Version: 2.0.0
//     License: MIT http://opensource.org/licenses/MIT
//     Host: localhost:9000
//
//     Consumes:
//     - application/json
//
//     Produces:
//     - application/json
//
//     Security:
//     - bearer: []
//
//     SecurityDefinitions:
//     bearer:
//          type: apiKey
//          name: Authorization
//          in: header
//
// swagger:meta
package api

import (
	"context"
	"crypto/sha1"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/labstack/echo"
	"os"

	"github.com/volatiletech/sqlboiler/boil"
	"github.com/wednesday-solutions/go-boiler/pkg/utl/postgres"
	"github.com/wednesday-solutions/go-boiler/pkg/utl/zlog"

	"github.com/wednesday-solutions/go-boiler/pkg/api/auth"
	al "github.com/wednesday-solutions/go-boiler/pkg/api/auth/logging"
	at "github.com/wednesday-solutions/go-boiler/pkg/api/auth/transport"
	"github.com/wednesday-solutions/go-boiler/pkg/api/password"
	pl "github.com/wednesday-solutions/go-boiler/pkg/api/password/logging"
	pt "github.com/wednesday-solutions/go-boiler/pkg/api/password/transport"
	"github.com/wednesday-solutions/go-boiler/pkg/api/user"
	ul "github.com/wednesday-solutions/go-boiler/pkg/api/user/logging"
	ut "github.com/wednesday-solutions/go-boiler/pkg/api/user/transport"

	_ "github.com/lib/pq" // here
	"github.com/wednesday-solutions/go-boiler/pkg/utl/config"
	"github.com/wednesday-solutions/go-boiler/pkg/utl/jwt"
	authMw "github.com/wednesday-solutions/go-boiler/pkg/utl/middleware/auth"
	"github.com/wednesday-solutions/go-boiler/pkg/utl/secure"
	"github.com/wednesday-solutions/go-boiler/pkg/utl/server"
)

// Start starts the API service
func Start(cfg *config.Configuration) error {
	db, err := postgres.Connect()
	if err != nil {
		return err
	}
	boil.SetDB(db)

	sec := secure.New(cfg.App.MinPasswordStr, sha1.New())
	jwt, err := jwt.New(cfg.JWT.SigningAlgorithm, os.Getenv("JWT_SECRET"), cfg.JWT.DurationMinutes, cfg.JWT.MinSecretLength)
	if err != nil {
		return err
	}

	log := zlog.New()

	e := server.New()
	e.Static("/swaggerui", cfg.App.SwaggerUIPath)

	authMiddleware := authMw.Middleware(jwt)

	at.NewHTTP(al.New(auth.Initialize(db, jwt, sec), log), e, authMiddleware)

	v1 := e.Group("/v1")
	v1.Use(authMiddleware)

	ut.NewHTTP(ul.New(user.Initialize(db, sec), log), v1)
	pt.NewHTTP(pl.New(password.Initialize(db, sec), log), v1)

	playgroundHandler := playground.Handler("GraphQL playground", "/graphql")

	// graphql playground
	e.GET("/playground", func(c echo.Context) error {
		req := c.Request()
		res := c.Response()
		playgroundHandler.ServeHTTP(res, req)
		return nil
	})
	server.Start(e, &server.Config{
		Port:                cfg.Server.Port,
		ReadTimeoutSeconds:  cfg.Server.ReadTimeout,
		WriteTimeoutSeconds: cfg.Server.WriteTimeout,
		Debug:               cfg.Server.Debug,
	})
	return nil
}

// CustomContext ...
type CustomContext struct {
	echo.Context
	ctx context.Context
}

// A private key for context that only this package can access. This is important
// to prevent collisions between different context uses
var userCtxKey = &contextKey{"user"}

type contextKey struct {
	name string
}

// Process is the middleware function.
func Process(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		ctx := context.WithValue(c.Request().Context(), userCtxKey, c)
		c.SetRequest(c.Request().WithContext(ctx))

		cc := &CustomContext{c, ctx}

		return next(cc)
	}
}
