package rediscache

import (
	"encoding/json"
	"fmt"
	"github.com/wednesday-solutions/go-template/daos"
	"github.com/wednesday-solutions/go-template/models"
	resultwrapper "github.com/wednesday-solutions/go-template/pkg/utl/result_wrapper"
)

// Service ...
type Service interface {
	GetUser(id int) (models.User, error)
	GetRole(id int) (models.Role, error)
}

// GetUser gets user from redis, if present, else from the database
func GetUser(userID int) (*models.User, error) {
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
	user, err = daos.FindUserByID(userID)
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
func GetRole(roleID int) (*models.Role, error) {
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
	role, err = daos.FindRoleByID(roleID)
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
