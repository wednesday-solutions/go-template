// Package user contains user application services
package user

import (
	"context"
	"github.com/labstack/echo"
	"github.com/volatiletech/null"
	"github.com/volatiletech/sqlboiler/boil"
	"github.com/volatiletech/sqlboiler/queries/qm"
	"github.com/wednesday-solution/go-boiler/models"

	"github.com/wednesday-solution/go-boiler"
)

// Create creates a new user account
func (u User) Create(c echo.Context, user models.User) (models.User, error) {
	user.Password = null.StringFrom(u.sec.Hash(user.Password.String))
	err := user.Insert(context.Background(), u.db, boil.Infer())
	return user, err
}

// List returns list of users
func (u User) List(c echo.Context, p goboiler.Pagination) (models.UserSlice, error) {
	users, err := models.Users(qm.Select()).All(context.Background(), u.db)
	if err != nil {
		return nil, err
	}
	return users, err
}

// View returns single user
func (u User) View(c echo.Context, id int) (*models.User, error) {
	return models.FindUser(context.Background(), u.db, id)
}

// Delete deletes a user
func (u User) Delete(c echo.Context, id int) error {
	user, err := models.FindUser(context.Background(), u.db, id)
	if err != nil {
		return err
	}
	_, err = user.Delete(context.Background(), u.db)
	return err
}

// Update updates user's contact information
func (u User) Update(c echo.Context, user models.User) (models.User, error) {
	_, err := user.Update(context.Background(), u.db, boil.Infer())
	return user, err
}
