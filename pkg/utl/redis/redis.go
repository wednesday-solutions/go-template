package redis

import (
	"encoding/json"
	"fmt"
	redigo "github.com/gomodule/redigo/redis"
	"os"
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

//type Storage interface {
//	Get(key string) []byte
//	Set(key string, content []byte, duration time.Duration)
//}
//
// Cached ...
//func Cached(duration string, handler func(w http.ResponseWriter, r *http.Request)) http.Handler {
//	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
//
//		var storage cache.Storage
//
//		strategy := flag.String("s", "memory", "Cache strategy (memory or redis)")
//		flag.Parse()
//
//		if *strategy == "memory" {
//			storage = memory.NewStorage()
//		} else if *strategy == "redis" {
//			var err error
//			if storage, err = redis.NewStorage(os.Getenv("REDIS_URL")); err != nil {
//				panic(err)
//			}
//		} else {
//
//		}
//		content := storage.Get(r.RequestURI)
//		if content != nil {
//			w.Write(content)
//		} else {
//			c := httptest.NewRecorder()
//			handler(c, r)
//
//			for k, v := range c.HeaderMap {
//				w.Header()[k] = v
//			}
//
//			w.WriteHeader(c.Code)
//			content := c.Body.Bytes()
//
//			if d, err := time.ParseDuration(duration); err == nil {
//				storage.Set(r.RequestURI, content, d)
//			}
//
//			_, _ = w.Write(content)
//		}
//
//	})
//}
