package password

import (
	"context"
	"fmt"
	"github.com/volatiletech/null"
	"github.com/volatiletech/sqlboiler/boil"
	"github.com/wednesday-solutions/go-boiler/models"
	"github.com/wednesday-solutions/go-boiler/pkg/utl"
	"net/http"

	"github.com/labstack/echo"
)

// Custom errors
var (
	ErrIncorrectPassword = echo.NewHTTPError(http.StatusBadRequest, "incorrect old password")
	ErrInsecurePassword  = echo.NewHTTPError(http.StatusBadRequest, "insecure password")
)

// Change changes user's password
func (p Password) Change(c echo.Context, userID int, oldPass, newPass string) error {
	u, err := models.FindUser(context.Background(), p.db, userID)
	if err != nil {
		return err
	}
	fmt.Print("\n\n\n\npass\n\n\nuserID: ",userID)
	if !p.sec.HashMatchesPassword(utl.FromNullableString(u.Password), oldPass) {
		return ErrIncorrectPassword
	}

	if !p.sec.Password(newPass, utl.FromNullableString(u.FirstName), utl.FromNullableString(u.LastName), utl.FromNullableString(u.Username), utl.FromNullableString(u.Email)) {
		return ErrInsecurePassword
	}

	u.Password = null.StringFrom(p.sec.Hash(newPass))
	_, err = u.Update(context.Background(), p.db, boil.Infer())
	return err
}
