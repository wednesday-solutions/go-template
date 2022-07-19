package rediscache

import (
	"encoding/json"
	"fmt"
	"reflect"
	"testing"

	. "github.com/agiledragon/gomonkey/v2"
	"github.com/gomodule/redigo/redis"
	redigo "github.com/gomodule/redigo/redis"
	redigomock "github.com/rafaeljusto/redigomock/v3"
)

var redigoConn = redigomock.NewConn()

func Test_redisDial(t *testing.T) {

	tests := []struct {
		name    string
		want    redigo.Conn
		wantErr bool
	}{
		{
			name: "Success",
			want: redigoConn,
		},
		{
			name:    "Failure",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.wantErr {
				ApplyFunc(redigo.Dial, func(string, string, ...redis.DialOption) (redis.Conn, error) {
					return nil, fmt.Errorf("some error")
				})
			} else {
				ApplyFunc(redigo.Dial, func(string, string, ...redis.DialOption) (redis.Conn, error) {
					return redigoConn, nil
				})
			}
			got, err := redisDial()
			if (err != nil) != tt.wantErr {
				t.Errorf("redisDial() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("redisDial() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSetKeyValue(t *testing.T) {
	type args struct {
		key  string
		data interface{}
	}
	// user := testutls.MockUser()
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "Success",
			args: args{
				key:  "user10",
				data: 1,
			},
		},
		{
			name: "Failure",
			args: args{
				key:  "user10",
				data: 1,
			},
			wantErr: true,
		},
	}
	ApplyFunc(redigo.Dial, func(string, string, ...redis.DialOption) (redis.Conn, error) {
		return redigoConn, nil
	})

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var patches *Patches
			b, _ := json.Marshal(tt.args.data)
			if tt.wantErr {
				patches = ApplyFunc(redigo.Dial, func(string, string, ...redis.DialOption) (redis.Conn, error) {
					return nil, fmt.Errorf("some error")
				})
			}
			redigoConn.Command("SET", tt.args.key, string(b)).Expect("something")

			if err := SetKeyValue(tt.args.key, tt.args.data); (err != nil) != tt.wantErr {
				t.Errorf("SetKeyValue() error = %v, wantErr %v", err, tt.wantErr)
			}
			if patches != nil {
				patches.Reset()
			}
		})
	}
}

func TestGetKeyValue(t *testing.T) {
	type args struct {
		key string
	}
	tests := []struct {
		name    string
		args    args
		want    interface{}
		wantErr bool
	}{
		{
			name: "Success",
			args: args{
				key: "user10",
			},
			want: "user",
		},
		{
			name: "Failure",
			args: args{
				key: "user10",
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var patches *Patches
			if tt.wantErr {
				patches = ApplyFunc(redigo.Dial, func(string, string, ...redis.DialOption) (redis.Conn, error) {
					return nil, fmt.Errorf("some error")
				})
			}
			redigoConn.Command("GET", tt.args.key).Expect(tt.want)

			got, err := GetKeyValue(tt.args.key)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetKeyValue() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetKeyValue() = %v, want %v", got, tt.want)
			}
			if patches != nil {
				patches.Reset()
			}
		})
	}
}
