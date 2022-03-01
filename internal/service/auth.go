package service

import (
	"crypto/sha1"
	"os"

	"go-template/internal/config"
	"go-template/internal/jwt"
	"go-template/pkg/utl/secure"
)

// Secure returns new secure service
func Secure(cfg *config.Configuration) *secure.Service {
	return secure.New(cfg.App.MinPasswordStr, sha1.New())
}

// JWT returns new JWT service
func JWT(cfg *config.Configuration) (jwt.Service, error) {
	return jwt.New(cfg.JWT.SigningAlgorithm, os.Getenv("JWT_SECRET"), cfg.JWT.DurationMinutes, cfg.JWT.MinSecretLength)
}
