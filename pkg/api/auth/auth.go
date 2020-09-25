package auth

import (
	"context"
	"github.com/volatiletech/null"
	"github.com/volatiletech/sqlboiler/boil"
	"github.com/volatiletech/sqlboiler/queries/qm"
	"github.com/wednesday-solutions/go-boiler/models"
	"net/http"

	"github.com/labstack/echo"

	"github.com/wednesday-solutions/go-boiler"
)

// Custom errors
var (
	ErrInvalidCredentials = echo.NewHTTPError(http.StatusUnauthorized, "Username or password does not exist")
)

// Authenticate tries to authenticate the user provided by username and password
func (a Auth) Authenticate(c echo.Context, user, pass string) (goboiler.AuthToken, error) {
	u, err := models.Users(qm.Where("username=?", user)).One(context.Background(), a.db)
	if err != nil {
		return goboiler.AuthToken{}, err
	}

	if !u.Password.Valid || (!a.sec.HashMatchesPassword(u.Password.String, pass)) {
		return goboiler.AuthToken{}, ErrInvalidCredentials
	}

	if !u.Active.Valid || (!u.Active.Bool) {
		return goboiler.AuthToken{}, goboiler.ErrUnauthorized
	}

	token, err := a.tg.GenerateToken(u)
	if err != nil {
		return goboiler.AuthToken{}, goboiler.ErrUnauthorized
	}

	refreshToken := a.sec.Token(token)
	u.Token = null.StringFrom(refreshToken)
	_, err = u.Update(context.Background(), a.db, boil.Infer())

	return goboiler.AuthToken{Token: token, RefreshToken: refreshToken}, err
}

// Refresh refreshes jwt token and puts new claims inside
func (a Auth) Refresh(c echo.Context, refreshToken string) (string, error) {
	user, err := models.Users(qm.Where("token=?", refreshToken)).One(context.Background(), a.db)
	if err != nil {
		return "", err
	}
	return a.tg.GenerateToken(user)
}

// Me returns info about currently logged user
func (a Auth) Me(c echo.Context) (*models.User, error) {
	return models.FindUser(context.Background(), a.db, c.Get("id").(int))
}
