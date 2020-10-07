package gotemplate

import (
	"encoding/json"
	"fmt"
	"github.com/wednesday-solutions/go-template/daos"
	"github.com/wednesday-solutions/go-template/models"
	"github.com/wednesday-solutions/go-template/pkg/utl"
	resultwrapper "github.com/wednesday-solutions/go-template/pkg/utl/result_wrapper"
)

// Service ...
type Service interface {
	GetUser(id int) (models.User, error)
}

// GetUser gets user from redis, if present, else from the database
func GetUser(userID int) (*models.User, error) {
	// get user cache key
	cachedUserValue, err := utl.GetKeyValue(fmt.Sprintf("user%d", userID))
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
	err = utl.SetKeyValue(fmt.Sprintf("user%d", userID), user)
	if err != nil {
		return nil, fmt.Errorf("%s", err)
	}
	return user, nil
}
