package transport

import (
	"github.com/wednesday-solution/go-boiler"
)

// Login request
// swagger:parameters login
type swaggLoginReq struct {
	// in:body
	Body credentials
}

// Login response
// swagger:response loginResp
type swaggLoginResp struct {
	// in:body
	Body struct {
		*goboiler.AuthToken
	}
}

// Token refresh response
// swagger:response refreshResp
type swaggRefreshResp struct {
	// in:body
	Body struct {
		*goboiler.RefreshToken
	}
}
