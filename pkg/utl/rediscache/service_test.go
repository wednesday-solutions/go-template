package rediscache

import (
	"context"
	"database/sql/driver"
	"encoding/json"
	"errors"
	"fmt"
	"math"
	"reflect"
	"regexp"
	"testing"
	"time"

	"go-template/models"
	"go-template/testutls"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/agiledragon/gomonkey/v2"
	. "github.com/agiledragon/gomonkey/v2"
	"github.com/gomodule/redigo/redis"
	redigo "github.com/gomodule/redigo/redis"
	redigomock "github.com/rafaeljusto/redigomock/v3"
)

const (
	ErrorFromCacheUserValue = "cacheUserValueError"
	SuccessCacheMiss        = "success_WithCacheMiss"
	ErrorFromJson           = "jsonError"
	ErrorRedisDial          = "redis Dial Error"
	ErrorConnDo             = "conn.Do Error"
	ErrMsgFromConnDo        = "error From Conn Do"
	ErrMsgFromRedisDial     = "error From Redis Dial"
	ErrorGetKeyValue        = "get Key Value Error Case"
	ErrorSetKeyValue        = "set Key Value Error Case"
	ErrorUnmarshal          = "unmarshalling Error Case"
	ErrorMarshal            = "marshalling Error Case"
	ErrorFindRoleById       = "find Role By Id Error Case"
	ErrMsgGetKeyValue       = "error From Get Key Value"
	ErrMsgSetKeyValue       = "error From Set Key Value"
	ErrMsgUnmarshal         = "error while Unmarshalling"
	ErrMsgMarshal           = "error while Marshalling"
	ErrMsgFindRoleById      = "error From Find Role By Id"
	ErrorDaos               = "error daos"
)

var conn = redigomock.NewConn()

func getDbQueryData() *sqlmock.Rows {
	return sqlmock.NewRows([]string{
		"id",
		"first_name",
		"last_name",
		"username",
		"email",
		"mobile",
		"address",
		"token",
		"password",
		"role_id",
		"active",
	}).AddRow(
		testutls.MockUser().ID,
		testutls.MockUser().FirstName,
		testutls.MockUser().LastName,
		testutls.MockUser().Username,
		testutls.MockUser().Email,
		testutls.MockUser().Mobile,
		testutls.MockUser().Address,
		testutls.MockUser().Token,
		testutls.MockUser().Password,
		testutls.MockUser().RoleID,
		testutls.MockUser().Active,
	)
}

type argsGetUser struct {
	userID int
	want   *models.User
}

type userTestCaseArgs struct {
	name string
	args argsGetUser

	wantErr bool
	errMsg  error
	init    func(sqlmock.Sqlmock, argsGetUser) *Patches
}

func getUserTestCases() []userTestCaseArgs {
	tests := []userTestCaseArgs{
		errorGetKeyValueCase(),
		errorUnmarshalCase(),
		errorSetKeyValueCase(),
		errorDaosCase(),
		getSuccessCase(),
		successCacheMissCase(),
		errorFromCacheUserValueCase(),
	}
	return tests
}

func errorGetKeyValueCase() userTestCaseArgs {
	return userTestCaseArgs{
		name: ErrorGetKeyValue,
		args: argsGetUser{
			userID: testutls.MockID,
		},
		wantErr: true,
		errMsg:  fmt.Errorf(ErrMsgGetKeyValue),
		init: func(mock sqlmock.Sqlmock, args argsGetUser) *Patches {
			conn.Command("GET", fmt.Sprintf("user%d", args.userID)).ExpectError(fmt.Errorf(ErrMsgGetKeyValue))
			return nil
		},
	}
}
func errorUnmarshalCase() userTestCaseArgs {
	return userTestCaseArgs{
		name: ErrorUnmarshal,
		args: argsGetUser{
			userID: testutls.MockID,
		},
		wantErr: true,
		init: func(s sqlmock.Sqlmock, user argsGetUser) *Patches {
			patchJson := ApplyFunc(json.Marshal, func(v any) ([]byte, error) {
				return []byte{}, fmt.Errorf("json error")
			})
			return patchJson
		},
	}
}
func errorSetKeyValueCase() userTestCaseArgs {
	return userTestCaseArgs{
		name: ErrorSetKeyValue,
		args: argsGetUser{
			userID: testutls.MockID,
		},
		init: func(mock sqlmock.Sqlmock, args argsGetUser) *Patches {
			conn.Command("GET", fmt.Sprintf("user%d", args.userID)).Expect(nil)
			b, _ := json.Marshal(testutls.MockUser())
			conn.Command("SET", fmt.Sprintf("user%d", args.userID), string(b)).ExpectError(fmt.Errorf("this is an error"))

			dbQueries := []testutls.QueryData{
				{
					Actions:    &[]driver.Value{testutls.MockID},
					Query:      `select * from "users" where "id"=$1`,
					DbResponse: getDbQueryData(),
				},
			}
			for _, dbQuery := range dbQueries {
				mock.ExpectQuery(regexp.QuoteMeta(dbQuery.Query)).
					WithArgs(*dbQuery.Actions...).
					WillReturnRows(dbQuery.DbResponse)
			}
			return nil
		},
		wantErr: true,
		errMsg:  fmt.Errorf(ErrMsgSetKeyValue),
	}
}
func errorDaosCase() userTestCaseArgs {
	return userTestCaseArgs{
		name: ErrorDaos,
		args: argsGetUser{
			userID: testutls.MockID,
		},
		init: func(mock sqlmock.Sqlmock, args argsGetUser) *Patches {
			conn.Command("GET", fmt.Sprintf("user%d", args.userID)).Expect(nil)
			dbQueries := []testutls.QueryData{
				{
					Actions:    &[]driver.Value{testutls.MockID},
					Query:      `select * from "users" where "id"=$1`,
					DbResponse: getDbQueryData().RowError(0, fmt.Errorf("data error")),
				},
			}
			for _, dbQuery := range dbQueries {
				mock.ExpectQuery(regexp.QuoteMeta(dbQuery.Query)).
					WithArgs(*dbQuery.Actions...).
					WillReturnRows(dbQuery.DbResponse)
			}
			return nil
		},
		wantErr: true,
		errMsg:  fmt.Errorf(ErrMsgSetKeyValue),
	}
}
func getSuccessCase() userTestCaseArgs {
	return userTestCaseArgs{
		name: SuccessCase,
		args: argsGetUser{
			userID: testutls.MockID,
			want:   testutls.MockUser(),
		},
		init: func(mock sqlmock.Sqlmock, args argsGetUser) *Patches {
			b, _ := json.Marshal(args.want)
			conn.Command("GET", fmt.Sprintf("user%d", args.userID)).Expect(b)
			return nil
		},
	}
}
func successCacheMissCase() userTestCaseArgs {
	return userTestCaseArgs{
		name: SuccessCacheMiss,
		args: argsGetUser{
			userID: testutls.MockID,
			want:   testutls.MockUser(),
		},
		init: func(mock sqlmock.Sqlmock, args argsGetUser) *Patches {
			conn.Command("GET", fmt.Sprintf("user%d", args.userID)).Expect(nil)
			b, _ := json.Marshal(args.want)
			conn.Command("SET", fmt.Sprintf("user%d", args.userID), string(b)).Expect(nil)
			dbQueries := []testutls.QueryData{
				{
					Actions:    &[]driver.Value{testutls.MockID},
					Query:      `select * from "users" where "id"=$1`,
					DbResponse: getDbQueryData(),
				},
			}
			for _, dbQuery := range dbQueries {
				mock.ExpectQuery(regexp.QuoteMeta(dbQuery.Query)).
					WithArgs(*dbQuery.Actions...).
					WillReturnRows(dbQuery.DbResponse)
			}
			return nil
		},
	}
}
func errorFromCacheUserValueCase() userTestCaseArgs {
	return userTestCaseArgs{
		name:    ErrorFromCacheUserValue,
		wantErr: true,
		args: argsGetUser{
			userID: testutls.MockID,
		},
		init: func(s sqlmock.Sqlmock, args argsGetUser) *Patches {
			conn.Command("GET", fmt.Sprintf("user%d", args.userID)).
				ExpectError(fmt.Errorf("error while getting from cache"))
			return nil
		},
	}
}

func TestGetUser(t *testing.T) {
	tests := getUserTestCases()
	redisDialPatch := ApplyFunc(redisDial, func() (redis.Conn, error) {
		return conn, nil
	})
	defer redisDialPatch.Reset()
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mock, cleanup, _ := testutls.SetupMockDB(t)
			patches := tt.init(mock, tt.args)
			got, err := GetUser(tt.args.userID, context.Background())
			fmt.Println(tt.name, "got", got)
			fmt.Println(tt.name, "err", err)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetUser() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.args.want) {
				t.Errorf("GetUser() = %v, want %v", got, tt.args.want)
			}
			if patches != nil {
				patches.Reset()
			}
			cleanup()
		})
	}
}

var role = &models.Role{
	ID:          1,
	AccessLevel: 100,
	Name:        "Admin",
}
var rowDbResponse = sqlmock.NewRows([]string{
	"id", "access_level", "name",
}).AddRow(
	role.ID,
	role.AccessLevel,
	role.Name,
)

type getRoleArgs struct {
	roleID    int
	cacheMiss bool
	want      *models.Role
}

type roleTestArgs struct {
	name string
	args getRoleArgs

	wantErr bool
	errMsg  error
	init    func(sqlmock.Sqlmock, getRoleArgs) *Patches
}

func getRoleTestCase(
	name string,
	args getRoleArgs,
	wantErr bool,
	errMsg error,
	init func(sqlmock.Sqlmock, getRoleArgs) *Patches) roleTestArgs {
	return roleTestArgs{
		name:    name,
		args:    args,
		wantErr: wantErr,
		errMsg:  errMsg,
		init:    init,
	}
}

func setupErrorCase(
	name string,
	roleID int,
	errMsg error,
	init func(sqlmock.Sqlmock, getRoleArgs) *Patches) roleTestArgs {
	return roleTestArgs{
		name: name,
		args: getRoleArgs{
			roleID: roleID,
		},
		wantErr: true,
		errMsg:  errMsg,
		init:    init,
	}
}

func loadGetRoleSuccessCase() roleTestArgs {
	return getRoleTestCase(SuccessCase, getRoleArgs{
		roleID: testutls.MockID,
		want:   role,
	},
		false,
		nil,
		func(mock sqlmock.Sqlmock, args getRoleArgs) *Patches {
			b, _ := json.Marshal(args.want)
			conn.Command("GET", fmt.Sprintf("role%d", args.roleID)).Expect(b)
			return nil
		})
}

func loadGetRoleCacheMissSuccessCase() roleTestArgs {
	return getRoleTestCase(
		SuccessCacheMiss,
		getRoleArgs{
			want:      role,
			roleID:    testutls.MockID,
			cacheMiss: true,
		},
		false,
		nil,
		func(mock sqlmock.Sqlmock, args getRoleArgs) *Patches {
			conn.Command("GET", fmt.Sprintf("role%d", args.roleID)).Expect(nil)
			b, _ := json.Marshal(args.want)
			conn.Command("SET", fmt.Sprintf("role%d", args.roleID), string(b)).Expect(nil)
			dbQueries := []testutls.QueryData{
				{
					Actions: &[]driver.Value{
						role.ID,
					},
					Query:      `select * from "roles" where "id"=$1`,
					DbResponse: rowDbResponse,
				},
			}
			for _, dbQuery := range dbQueries {
				mock.ExpectQuery(regexp.QuoteMeta(dbQuery.Query)).
					WithArgs(*dbQuery.Actions...).
					WillReturnRows(dbQuery.DbResponse)
			}
			return nil
		})
}
func getGetRoleTestCases() []roleTestArgs {
	tests := []roleTestArgs{
		setupErrorCase(ErrorGetKeyValue, testutls.MockID, errors.New(ErrMsgGetKeyValue),
			func(mock sqlmock.Sqlmock, args getRoleArgs) *Patches {
				conn.Command("GET", fmt.Sprintf("role%d", args.roleID)).ExpectError(fmt.Errorf(ErrMsgGetKeyValue))
				return nil
			}),
		setupErrorCase(ErrorUnmarshal, testutls.MockID, errors.New(ErrMsgUnmarshal),
			func(mock sqlmock.Sqlmock, args getRoleArgs) *Patches {
				return ApplyFunc(json.Unmarshal, func(data []byte, v any) error {
					return fmt.Errorf(ErrMsgUnmarshal)
				})
			}),
		setupErrorCase(ErrorSetKeyValue, testutls.MockID, errors.New(ErrMsgSetKeyValue),
			func(mock sqlmock.Sqlmock, args getRoleArgs) *Patches {
				conn.Command("GET", fmt.Sprintf("role%d", args.roleID)).Expect(nil)
				return nil
			}),
		setupErrorCase(ErrorFindRoleById, testutls.MockID, errors.New(ErrMsgFindRoleById),
			func(mock sqlmock.Sqlmock, args getRoleArgs) *Patches {
				conn.Command("GET", fmt.Sprintf("role%d", args.roleID)).ExpectError(fmt.Errorf("there was an error"))
				return nil
			}),
		loadGetRoleSuccessCase(),
		loadGetRoleCacheMissSuccessCase(),
	}
	return tests
}

func TestGetRole(t *testing.T) {
	tests := getGetRoleTestCases()
	mock, cleanup, _ := testutls.SetupMockDB(t)
	defer cleanup()
	redisDialPatches := ApplyFunc(redisDial, func() (redis.Conn, error) {
		return conn, nil
	})
	defer redisDialPatches.Reset()
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			patches := tt.init(mock, tt.args)

			got, err := GetRole(tt.args.roleID, context.Background())
			if (err != nil) != tt.wantErr {
				t.Errorf("GetRole() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.args.want) {
				t.Errorf("GetRole() = %v, want %v", got, tt.args.want)
			}
			if patches != nil {
				patches.Reset()
			}
		})
	}
}
func TestIncVisits(t *testing.T) {
	type args struct {
		path string
	}
	tests := []struct {
		name    string
		args    args
		want    int
		wantErr bool
		errMsg  error
	}{
		{
			name: SuccessCase,
			args: args{
				path: "test",
			},
			want: 10,
		},
		{
			name: ErrorRedisDial,
			args: args{
				path: "test",
			},
			wantErr: true,
			errMsg:  fmt.Errorf(ErrMsgFromRedisDial),
		},
	}

	ApplyFunc(redisDial, func() (redis.Conn, error) {
		return conn, nil
	})
	for _, tt := range tests {
		if tt.name == ErrorRedisDial {
			patch := gomonkey.ApplyFunc(redisDial, func() (redigo.Conn, error) {
				return nil, fmt.Errorf(ErrMsgFromRedisDial)
			})
			defer patch.Reset()
		}
		t.Run(tt.name, func(t *testing.T) {
			conn.Command("INCR", tt.args.path).Expect([]byte(fmt.Sprint(tt.want)))
			got, err := IncVisits(tt.args.path)
			if (err != nil) != tt.wantErr {
				t.Errorf("IncVisits() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("IncVisits() = %v, want %v", got, tt.want)
			}
		})
	}
}

type startVisitArgs struct {
	path string
	want int
}

type startVisitTestArgs struct {
	name string
	args startVisitArgs

	wantErr bool
	errMsg  error
	init    func(startVisitArgs) *Patches
}

func TestStartVisits(t *testing.T) {
	tests := getTestCases()

	for _, tt := range tests {
		patches := tt.init(tt.args)

		t.Run(tt.name, func(t *testing.T) {
			err := StartVisits(tt.args.path, time.Second)

			verifyError(t, err, tt.wantErr)
		})
		if patches != nil {
			patches.Reset()
		}
	}
}

func getTestCases() []startVisitTestArgs {
	return []startVisitTestArgs{
		{
			name: ErrorConnDo,
			args: startVisitArgs{
				path: "test",
			},
			wantErr: true,
			errMsg:  fmt.Errorf(ErrMsgFromConnDo),
			init: func(args startVisitArgs) *Patches {
				conn.Command("SETEX", args.path, int(math.Ceil(time.Second.Seconds())), 1).Expect(args.want)
				return gomonkey.ApplyMethodFunc(redigomock.NewConn(), "Do",
					func(commandName string, args ...interface{}) (reply interface{}, err error) {
						return nil, fmt.Errorf(ErrMsgFromConnDo)
					})
			},
		},
		{
			name: SuccessCase,
			args: startVisitArgs{
				path: "test",
				want: 10,
			},
			init: func(args startVisitArgs) *Patches {
				return gomonkey.ApplyFunc(redisDial, func() (redigo.Conn, error) {
					return conn, nil
				})
			},
		},
		{
			name: ErrorRedisDial,
			args: startVisitArgs{
				path: "test",
			},
			wantErr: true,
			errMsg:  fmt.Errorf(ErrMsgFromRedisDial),
			init: func(args startVisitArgs) *Patches {
				conn.Command("SETEX", args.path, int(math.Ceil(time.Second.Seconds())), 1).Expect(args.want)
				return gomonkey.ApplyFunc(redisDial, func() (redigo.Conn, error) {
					return nil, fmt.Errorf(ErrMsgFromRedisDial)
				})
			},
		},
	}
}

func verifyError(t *testing.T, err error, wantErr bool) {
	if (err != nil) != wantErr {
		t.Errorf("StartVisits() error = %v, wantErr %v", err, wantErr)
	}
}
