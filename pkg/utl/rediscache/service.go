package rediscache

import (
	"context"
	"encoding/json"
	"fmt"
	"math"
	"time"

	"go-template/daos"
	"go-template/models"
	resultwrapper "go-template/pkg/utl/resultwrapper"

	redigo "github.com/gomodule/redigo/redis"
)

// Service ...
type Service interface {
	GetUser(id int, ctx context.Context) (models.User, error)
	GetRole(id int, ctx context.Context) (models.Role, error)
	IncVisits(path string) (int, error)
	StartVisits(path string, exp time.Duration) error
}

// GetUser gets user from redis, if present, else from the database
func GetUser(userID int, ctx context.Context) (*models.User, error) {
	// get user cache key
	cachedUserValue, err := GetKeyValue(fmt.Sprintf("user%d", userID))
	if err != nil {
		return nil, err
	}
	var user *models.User
	if cachedUserValue != nil {
		b := cachedUserValue.([]byte)
		err = json.Unmarshal(b, &user)
		if err != nil {
			return nil, fmt.Errorf("%s", err)
		}
		return user, err
	}
	user, err = daos.FindUserByID(userID, ctx)
	if err != nil {
		return nil, resultwrapper.ResolverSQLError(err, "data")
	}
	// setting user cache key
	err = SetKeyValue(fmt.Sprintf("user%d", userID), user)
	if err != nil {
		return nil, fmt.Errorf("%s", err)
	}
	return user, nil
}

// GetRole gets role from redis, if present, else from the database
func GetRole(roleID int, ctx context.Context) (*models.Role, error) {
	// get role cache key
	cachedRoleValue, err := GetKeyValue(fmt.Sprintf("role%d", roleID))
	if err != nil {
		return nil, err
	}
	var role *models.Role
	if cachedRoleValue != nil {
		b := cachedRoleValue.([]byte)
		err = json.Unmarshal(b, &role)
		if err != nil {
			return nil, fmt.Errorf("%s", err)
		}
		return role, err
	}
	role, err = daos.FindRoleByID(roleID, ctx)
	if err != nil {
		return nil, resultwrapper.ResolverSQLError(err, "data")
	}
	// setting role cache key
	err = SetKeyValue(fmt.Sprintf("role%d", roleID), role)
	if err != nil {
		return nil, fmt.Errorf("%s", err)
	}
	return role, nil
}

// IncVisits Increases the no. of visits by a particular visitor on a
// particular graphQL path by one, or returns 1 if visiting 1st time.
func IncVisits(path string) (int, error) {
	conn, err := redisDial()
	if err != nil {
		return 0, fmt.Errorf("error in redis connection %s", err)
	}
	defer conn.Close()

	return redigo.Int(conn.Do("INCR", path))
}

// StartVisits is called when the visiter is first time entering the
// given path or no entry of the visiter is present because of time-out, It sets the path with expiry as exp
func StartVisits(path string, exp time.Duration) error {
	conn, err := redisDial()
	if err != nil {
		return fmt.Errorf("error in redis connection %s", err)
	}
	defer conn.Close()

	ttl := math.Ceil(exp.Seconds())

	_, err = conn.Do("SETEX", path, int(ttl), 1)
	if err != nil {
		return err
	}
	return nil
}
