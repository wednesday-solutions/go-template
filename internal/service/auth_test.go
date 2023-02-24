package service_test

import (
	"log"
	"os"
	"testing"

	"go-template/internal/config"
	"go-template/internal/service"
	"go-template/testutls"

	. "github.com/agiledragon/gomonkey/v2"
	"github.com/stretchr/testify/assert"
)

const SuccessCase = "Success"

func TestSecure(t *testing.T) {
	type args struct {
		cfg *config.Configuration
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: SuccessCase,
			args: args{
				cfg: testutls.MockConfig(),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := service.Secure(tt.args.cfg)
			assert.NotNil(t, s)
		})
	}
}

func TestJWT(t *testing.T) {
	type args struct {
		cfg *config.Configuration
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: SuccessCase,
			args: args{
				cfg: testutls.MockConfig(),
			},
		},
	}
	patches := ApplyFunc(os.Getenv, func(s string) string {
		return testutls.MockJWTSecret
	})
	defer patches.Reset()
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := service.JWT(tt.args.cfg)
			if err != nil {
				log.Fatal(err)
			}
			assert.NotNil(t, got)
		})
	}
}
