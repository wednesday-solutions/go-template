package service_test

import (
	"log"
	"os"
	"testing"

	. "github.com/agiledragon/gomonkey/v2"
	"github.com/stretchr/testify/assert"
	"github.com/wednesday-solutions/go-template/internal/config"
	"github.com/wednesday-solutions/go-template/internal/service"
	"github.com/wednesday-solutions/go-template/testutls"
)

func TestSecure(t *testing.T) {
	type args struct {
		cfg *config.Configuration
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "Success",
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
			name: "Success",
			args: args{
				cfg: testutls.MockConfig(),
			},
		},
	}
	ApplyFunc(os.Getenv, func(s string) string {
		return "1234567890123456789012345678901234567890123456789012345678901234567890"
	})
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
