package rediscache

import (
	"encoding/json"
	"fmt"
	"os"

	redigo "github.com/gomodule/redigo/redis"
)

func redisDial() (redigo.Conn, error) {
	conn, err := redigo.Dial("tcp", os.Getenv("REDIS_ADDRESS"))
	// Connection error handling
	if err != nil {
		return conn, err
	}
	return conn, err
}

// SetKeyValue ...
func SetKeyValue(key string, data interface{}) error {
	conn, err := redisDial()
	if err != nil {
		return fmt.Errorf("error in redis connection %s", err)
	}
	b, err := json.Marshal(data)
	if err != nil {
		return err
	}
	_, err = conn.Do("SET", key, string(b))
	return err
}

// GetKeyValue ...
func GetKeyValue(key string) (interface{}, error) {
	conn, err := redisDial()
	if err != nil {
		return nil, fmt.Errorf("error in redis connection %s", err)
	}

	reply, err := conn.Do("GET", key)
	return reply, err
}
