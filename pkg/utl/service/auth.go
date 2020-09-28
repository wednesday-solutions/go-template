package service

import (
	"crypto/sha1"
	"github.com/wednesday-solutions/go-boiler/pkg/utl/config"
	"github.com/wednesday-solutions/go-boiler/pkg/utl/jwt"
	"github.com/wednesday-solutions/go-boiler/pkg/utl/secure"
	"os"
)

// Secure returns new secure service
func Secure(cfg *config.Configuration) *secure.Service {
	return secure.New(cfg.App.MinPasswordStr, sha1.New())
}

// JWT returns new JWT service
func JWT(cfg *config.Configuration) (jwt.Service, error) {
	return jwt.New(cfg.JWT.SigningAlgorithm, os.Getenv("JWT_SECRET"), cfg.JWT.DurationMinutes, cfg.JWT.MinSecretLength)
}
