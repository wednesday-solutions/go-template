package password

import (
	goboiler "github.com/wednesday-solutions/go-boiler"
	"time"

	"github.com/labstack/echo"
	"github.com/wednesday-solutions/go-boiler/pkg/api/password"
)

// New creates new password logging service
func New(svc password.Service, logger goboiler.Logger) *LogService {
	return &LogService{
		Service: svc,
		logger:  logger,
	}
}

// LogService represents password logging service
type LogService struct {
	password.Service
	logger goboiler.Logger
}

const name = "password"

// Change logging
func (ls *LogService) Change(c echo.Context, id int, oldPass, newPass string) (err error) {
	defer func(begin time.Time) {
		ls.logger.Log(
			c,
			name, "Change password request", err,
			map[string]interface{}{
				"req":  id,
				"took": time.Since(begin),
			},
		)
	}(time.Now())
	return ls.Service.Change(c, id, oldPass, newPass)
}
