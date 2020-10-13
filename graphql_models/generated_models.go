// Code generated by github.com/99designs/gqlgen, DO NOT EDIT.

package graphql_models

type BooleanFilter struct {
	IsTrue  *bool `json:"isTrue"`
	IsFalse *bool `json:"isFalse"`
	IsNull  *bool `json:"isNull"`
}

type ChangePasswordResponse struct {
	Ok bool `json:"ok"`
}

type FloatFilter struct {
	EqualTo           *float64  `json:"equalTo"`
	NotEqualTo        *float64  `json:"notEqualTo"`
	LessThan          *float64  `json:"lessThan"`
	LessThanOrEqualTo *float64  `json:"lessThanOrEqualTo"`
	MoreThan          *float64  `json:"moreThan"`
	MoreThanOrEqualTo *float64  `json:"moreThanOrEqualTo"`
	In                []float64 `json:"in"`
	NotIn             []float64 `json:"notIn"`
}

type IDFilter struct {
	EqualTo    *string  `json:"equalTo"`
	NotEqualTo *string  `json:"notEqualTo"`
	In         []string `json:"in"`
	NotIn      []string `json:"notIn"`
}

type IntFilter struct {
	EqualTo           *int  `json:"equalTo"`
	NotEqualTo        *int  `json:"notEqualTo"`
	LessThan          *int  `json:"lessThan"`
	LessThanOrEqualTo *int  `json:"lessThanOrEqualTo"`
	MoreThan          *int  `json:"moreThan"`
	MoreThanOrEqualTo *int  `json:"moreThanOrEqualTo"`
	In                []int `json:"in"`
	NotIn             []int `json:"notIn"`
}

type LoginResponse struct {
	Token        string `json:"token"`
	RefreshToken string `json:"refreshToken"`
}

type RefreshTokenResponse struct {
	Token string `json:"token"`
}

type Role struct {
	ID          string  `json:"id"`
	AccessLevel int     `json:"accessLevel"`
	Name        string  `json:"name"`
	UpdatedAt   *int    `json:"updatedAt"`
	DeletedAt   *int    `json:"deletedAt"`
	CreatedAt   *int    `json:"createdAt"`
	Users       []*User `json:"users"`
}

type RoleCreateInput struct {
	AccessLevel int    `json:"accessLevel"`
	Name        string `json:"name"`
}

type RoleDeletePayload struct {
	ID string `json:"id"`
}

type RoleFilter struct {
	Search *string    `json:"search"`
	Where  *RoleWhere `json:"where"`
}

type RolePagination struct {
	Limit int `json:"limit"`
	Page  int `json:"page"`
}

type RolePayload struct {
	Role *Role `json:"role"`
}

type RoleUpdateInput struct {
	AccessLevel *int    `json:"accessLevel"`
	Name        *string `json:"name"`
	UpdatedAt   *int    `json:"updatedAt"`
	DeletedAt   *int    `json:"deletedAt"`
	CreatedAt   *int    `json:"createdAt"`
}

type RoleWhere struct {
	ID          *IDFilter     `json:"id"`
	AccessLevel *IntFilter    `json:"accessLevel"`
	Name        *StringFilter `json:"name"`
	UpdatedAt   *IntFilter    `json:"updatedAt"`
	DeletedAt   *IntFilter    `json:"deletedAt"`
	CreatedAt   *IntFilter    `json:"createdAt"`
	Users       *UserWhere    `json:"users"`
	Or          *RoleWhere    `json:"or"`
	And         *RoleWhere    `json:"and"`
}

type RolesCreateInput struct {
	Roles []*RoleCreateInput `json:"roles"`
}

type RolesDeletePayload struct {
	Ids []string `json:"ids"`
}

type RolesPayload struct {
	Roles []*Role `json:"roles"`
}

type RolesUpdatePayload struct {
	Ok bool `json:"ok"`
}

type StringFilter struct {
	EqualTo            *string  `json:"equalTo"`
	NotEqualTo         *string  `json:"notEqualTo"`
	In                 []string `json:"in"`
	NotIn              []string `json:"notIn"`
	StartWith          *string  `json:"startWith"`
	NotStartWith       *string  `json:"notStartWith"`
	EndWith            *string  `json:"endWith"`
	NotEndWith         *string  `json:"notEndWith"`
	Contain            *string  `json:"contain"`
	NotContain         *string  `json:"notContain"`
	StartWithStrict    *string  `json:"startWithStrict"`
	NotStartWithStrict *string  `json:"notStartWithStrict"`
	EndWithStrict      *string  `json:"endWithStrict"`
	NotEndWithStrict   *string  `json:"notEndWithStrict"`
	ContainStrict      *string  `json:"containStrict"`
	NotContainStrict   *string  `json:"notContainStrict"`
}

type User struct {
	ID                 string  `json:"id"`
	FirstName          *string `json:"firstName"`
	LastName           *string `json:"lastName"`
	Username           *string `json:"username"`
	Password           *string `json:"password"`
	Email              *string `json:"email"`
	Mobile             *string `json:"mobile"`
	Phone              *string `json:"phone"`
	Address            *string `json:"address"`
	Active             *bool   `json:"active"`
	LastLogin          *int    `json:"lastLogin"`
	LastPasswordChange *int    `json:"lastPasswordChange"`
	Token              *string `json:"token"`
	Role               *Role   `json:"role"`
	CreatedAt          *int    `json:"createdAt"`
	DeletedAt          *int    `json:"deletedAt"`
	UpdatedAt          *int    `json:"updatedAt"`
}

type UserCreateInput struct {
	FirstName *string `json:"firstName"`
	LastName  *string `json:"lastName"`
	Username  *string `json:"username"`
	Password  *string `json:"password"`
	Email     *string `json:"email"`
	RoleID    *string `json:"roleId"`
}

type UserDeletePayload struct {
	ID string `json:"id"`
}

type UserFilter struct {
	Search *string    `json:"search"`
	Where  *UserWhere `json:"where"`
}

type UserPagination struct {
	Limit int `json:"limit"`
	Page  int `json:"page"`
}

type UserPayload struct {
	User *User `json:"user"`
}

type UserUpdateInput struct {
	FirstName *string `json:"firstName"`
	LastName  *string `json:"lastName"`
	Mobile    *string `json:"mobile"`
	Phone     *string `json:"phone"`
	Address   *string `json:"address"`
}

type UserUpdatePayload struct {
	Ok bool `json:"ok"`
}

type UserWhere struct {
	ID                 *IDFilter      `json:"id"`
	FirstName          *StringFilter  `json:"firstName"`
	LastName           *StringFilter  `json:"lastName"`
	Username           *StringFilter  `json:"username"`
	Password           *StringFilter  `json:"password"`
	Email              *StringFilter  `json:"email"`
	Mobile             *StringFilter  `json:"mobile"`
	Phone              *StringFilter  `json:"phone"`
	Address            *StringFilter  `json:"address"`
	Active             *BooleanFilter `json:"active"`
	LastLogin          *IntFilter     `json:"lastLogin"`
	LastPasswordChange *IntFilter     `json:"lastPasswordChange"`
	Token              *StringFilter  `json:"token"`
	Role               *RoleWhere     `json:"role"`
	CreatedAt          *IntFilter     `json:"createdAt"`
	DeletedAt          *IntFilter     `json:"deletedAt"`
	UpdatedAt          *IntFilter     `json:"updatedAt"`
	Or                 *UserWhere     `json:"or"`
	And                *UserWhere     `json:"and"`
}

type UsersCreateInput struct {
	Users []*UserCreateInput `json:"users"`
}

type UsersDeletePayload struct {
	Ids []string `json:"ids"`
}

type UsersPayload struct {
	Users []*User `json:"users"`
}
