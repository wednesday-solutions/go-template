package password

import (
	"fmt"
	"github.com/wednesday-solution/go-boiler/pkg/utl/secure"
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
	if err := p.rbac.EnforceUser(c, userID); err != nil {
		return err
	}

	u, err := p.udb.View(p.db, userID)
	if err != nil {
		return err
	}

	fmt.Print(u.Password, "\n\n\n\n")
	fmt.Print(oldPass, "\n\n\n\n")
	sec := secure.New(5, nil)
	fmt.Print(sec.Hash(oldPass), "\n\n\n\n")

	if !p.sec.HashMatchesPassword(u.Password, oldPass) {
		return ErrIncorrectPassword
	}

	if !p.sec.Password(newPass, u.FirstName, u.LastName, u.Username, u.Email) {
		return ErrInsecurePassword
	}

	u.ChangePassword(p.sec.Hash(newPass))

	return p.udb.Update(p.db, u)
}
