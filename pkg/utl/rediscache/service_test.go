package rediscache

import (
	"context"
	"database/sql/driver"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"math"
	"reflect"
	"regexp"
	"testing"
	"time"

	"go-template/internal/config"
	"go-template/models"
	"go-template/pkg/utl/convert"
	"go-template/testutls"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/agiledragon/gomonkey/v2"
	. "github.com/agiledragon/gomonkey/v2"
	"github.com/gomodule/redigo/redis"
	redigo "github.com/gomodule/redigo/redis"
	redigomock "github.com/rafaeljusto/redigomock/v3"
	"github.com/volatiletech/sqlboiler/v4/boil"
)

const (
	ErrorFromCacheUserValue = "CacheUserValueError"
	SuccessCacheMiss        = "Success_WithCacheMiss"
	ErrorFromJson           = "jsonError"
	ErrorRedisDial          = "Redis Dial Error"
	ErrorConnDo             = "Conn.Do Error"
	ErrMsgFromConnDo        = "Error From Conn Do"
	ErrMsgFromRedisDial     = "Error From Redis Dial"
	ErrorGetKeyValue        = "Get Key Value Error Case"
	ErrorSetKeyValue        = "Set Key Value Error Case"
	ErrorUnmarshal          = "Unmarshalling Error Case"
	ErrorMarshal            = "Marshalling Error Case"
	ErrorFindRoleById       = "Find Role By Id Error Case"
	ErrMsgGetKeyValue       = "Error From Get Key Value"
	ErrMsgSetKeyValue       = "Error From Set Key Value"
	ErrMsgUnmarshal         = "Error while Unmarshalling"
	ErrMsgMarshal           = "Error while Marshalling"
	ErrMsgFindRoleById      = "Error From Find Role By Id"
	ErrorDaos               = "Error daos"
)

var conn = redigomock.NewConn()
var dbQuerydata = sqlmock.NewRows([]string{
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

func TestGetUser(t *testing.T) {
	type args struct {
		userID    int
		cacheMiss bool
		dbQueries []testutls.QueryData
	}
	tests := []struct {
		name    string
		args    args
		want    *models.User
		wantErr bool
		errMsg  error
	}{
		{
			name: ErrorGetKeyValue,
			args: args{
				userID:    testutls.MockID,
				dbQueries: []testutls.QueryData{},
			},
			wantErr: true,
			errMsg:  fmt.Errorf(ErrMsgGetKeyValue),
		},
		{
			name: ErrorUnmarshal,
			args: args{
				userID:    testutls.MockID,
				dbQueries: []testutls.QueryData{},
			},
			wantErr: true,
		},
		{
			name: ErrorSetKeyValue,
			args: args{
				userID:    testutls.MockID,
				cacheMiss: true,
				dbQueries: []testutls.QueryData{
					{
						Actions:    &[]driver.Value{testutls.MockID},
						Query:      `select * from "users" where "id"=$1`,
						DbResponse: dbQuerydata,
					},
				},
			},
			wantErr: true,
			errMsg:  fmt.Errorf(ErrMsgSetKeyValue),
		},
		{
			name: ErrorDaos,
			args: args{
				userID:    testutls.MockID,
				cacheMiss: true,
				dbQueries: []testutls.QueryData{
					{
						Actions:    &[]driver.Value{testutls.MockID},
						Query:      `select * from "users" where "id"=$1`,
						DbResponse: dbQuerydata.RowError(0, fmt.Errorf("data error")),
					},
				},
			},
			wantErr: true,
			errMsg:  fmt.Errorf(ErrMsgSetKeyValue),
		},

		{
			name: SuccessCase,
			args: args{
				userID:    testutls.MockID,
				dbQueries: []testutls.QueryData{},
			},
			want: testutls.MockUser(),
		},
		{
			name: SuccessCacheMiss,
			args: args{
				userID:    testutls.MockID,
				cacheMiss: true,
				dbQueries: []testutls.QueryData{
					{
						Actions:    &[]driver.Value{testutls.MockID},
						Query:      `select * from "users" where "id"=$1`,
						DbResponse: dbQuerydata,
					},
				},
			},
			want: testutls.MockUser(),
		},
		{
			name: ErrorFromCacheUserValue,
			args: args{
				userID:    testutls.MockID,
				dbQueries: []testutls.QueryData{},
			},
		},
	}

	oldDB := boil.GetDB()
	err := config.LoadEnvWithFilePrefix(convert.StringToPointerString("./../../../"))
	if err != nil {
		log.Fatal(err)
	}
	mock, db, _ := testutls.SetupMockDB(t)

	for _, tt := range tests {
		ApplyFunc(redisDial, func() (redis.Conn, error) {
			return conn, nil
		})
		t.Run(tt.name, func(t *testing.T) {
			if tt.name == ErrorUnmarshal {
				patchJson := ApplyFunc(json.Marshal, func(v any) ([]byte, error) {
					return []byte{}, fmt.Errorf("json error")
				})
				defer patchJson.Reset()
			}
			if tt.args.cacheMiss {
				conn.Command("GET", fmt.Sprintf("user%d", tt.args.userID)).Expect(nil)
				b, _ := json.Marshal(tt.want)
				conn.Command("SET", fmt.Sprintf("user%d", tt.args.userID), string(b)).Expect(nil)
				for _, dbQuery := range tt.args.dbQueries {
					mock.ExpectQuery(regexp.QuoteMeta(dbQuery.Query)).
						WithArgs(*dbQuery.Actions...).
						WillReturnRows(dbQuery.DbResponse)
				}
			} else if tt.name == ErrorGetKeyValue {
				conn.Command("GET", fmt.Sprintf("role%d", tt.args.userID)).ExpectError(fmt.Errorf(ErrMsgGetKeyValue))
			} else if tt.name == ErrorSetKeyValue {
				conn.Command("GET", fmt.Sprintf("role%d", tt.args.userID)).Expect(nil)
			} else {
				b, _ := json.Marshal(tt.want)
				conn.Command("GET", fmt.Sprintf("user%d", tt.args.userID)).Expect(b)
			}
			got, err := GetUser(tt.args.userID, context.Background())
			if (err != nil) != tt.wantErr {
				t.Errorf("GetUser() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetUser() = %v, want %v", got, tt.want)
			}
		})
	}
	boil.SetDB(oldDB)
	db.Close()
}

var role = &models.Role{
	ID:          1,
	AccessLevel: 100,
	Name:        "Admin",
}

var DbResponse = sqlmock.NewRows([]string{
	"id", "access_level", "name",
}).AddRow(
	role.ID,
	role.AccessLevel,
	role.Name,
)

type GetRoleArgs struct {
	roleID    int
	cacheMiss bool
	dbQueries []testutls.QueryData
}

func getRoleTestCase(name string, args GetRoleArgs, want *models.Role, wantErr bool, errMsg error) struct {
	name    string
	args    GetRoleArgs
	want    *models.Role
	wantErr bool
	errMsg  error
} {
	return struct {
		name    string
		args    GetRoleArgs
		want    *models.Role
		wantErr bool
		errMsg  error
	}{
		name:    name,
		args:    args,
		want:    want,
		wantErr: wantErr,
		errMsg:  errMsg,
	}
}

func setupErrorCase(name string, roleID int, errMsg error) struct {
	name    string
	args    GetRoleArgs
	want    *models.Role
	wantErr bool
	errMsg  error
} {
	return struct {
		name    string
		args    GetRoleArgs
		want    *models.Role
		wantErr bool
		errMsg  error
	}{
		name: name,
		args: GetRoleArgs{
			roleID:    roleID,
			dbQueries: []testutls.QueryData{},
		},
		wantErr: true,
		errMsg:  errMsg,
	}
}

func getGetRoleTestCases() []struct {
	name    string
	args    GetRoleArgs
	want    *models.Role
	wantErr bool
	errMsg  error
} {
	role := &models.Role{ID: testutls.MockID, Name: testutls.MockUser().R.Role.Name}

	tests := []struct {
		name    string
		args    GetRoleArgs
		want    *models.Role
		wantErr bool
		errMsg  error
	}{
		setupErrorCase(ErrorGetKeyValue, testutls.MockID, errors.New(ErrMsgGetKeyValue)),
		setupErrorCase(ErrorUnmarshal, testutls.MockID, errors.New(ErrMsgUnmarshal)),
		setupErrorCase(ErrorSetKeyValue, testutls.MockID, errors.New(ErrMsgSetKeyValue)),
		setupErrorCase(ErrorFindRoleById, testutls.MockID, errors.New(ErrMsgFindRoleById)),
		getRoleTestCase(SuccessCase, GetRoleArgs{roleID: testutls.MockID},
			role, false, nil),
		getRoleTestCase(SuccessCacheMiss, GetRoleArgs{roleID: testutls.MockID,
			cacheMiss: true, dbQueries: []testutls.QueryData{
				{Actions: &[]driver.Value{role.ID},
					Query:      `select * from "roles" where "id"=$1`,
					DbResponse: DbResponse}}},
			role, false, nil),
	}
	return tests
}

func TestGetRole(t *testing.T) {
	tests := getGetRoleTestCases()
	oldDB := boil.GetDB()
	err := config.LoadEnvWithFilePrefix(convert.StringToPointerString("./../../../"))
	if err != nil {
		log.Fatal(err)
	}
	mock, db, _ := testutls.SetupMockDB(t)
	for _, tt := range tests {
		ApplyFunc(redisDial, func() (redis.Conn, error) {
			return conn, nil
		})
		t.Run(tt.name, func(t *testing.T) {
			if tt.name == ErrorUnmarshal {
				patchJson := ApplyFunc(json.Unmarshal, func(data []byte, v any) error {
					return fmt.Errorf(ErrMsgUnmarshal)
				})
				defer patchJson.Reset()
			}
			if tt.args.cacheMiss {
				conn.Command("GET", fmt.Sprintf("role%d", tt.args.roleID)).Expect(nil)
				b, _ := json.Marshal(tt.want)
				conn.Command("SET", fmt.Sprintf("role%d", tt.args.roleID), string(b)).Expect(nil)
				for _, dbQuery := range tt.args.dbQueries {
					mock.ExpectQuery(regexp.QuoteMeta(dbQuery.Query)).
						WithArgs(*dbQuery.Actions...).
						WillReturnRows(dbQuery.DbResponse)
				}
			} else if tt.name == ErrorGetKeyValue {
				conn.Command("GET", fmt.Sprintf("role%d", tt.args.roleID)).ExpectError(fmt.Errorf(ErrMsgGetKeyValue))
			} else if tt.name == ErrorSetKeyValue {
				conn.Command("GET", fmt.Sprintf("role%d", tt.args.roleID)).Expect(nil)
			} else {
				b, _ := json.Marshal(tt.want)
				conn.Command("GET", fmt.Sprintf("role%d", tt.args.roleID)).Expect(b)
			}
			got, err := GetRole(tt.args.roleID, context.Background())
			if (err != nil) != tt.wantErr {
				t.Errorf("GetRole() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetRole() = %v, want %v", got, tt.want)
			}
		})
	}
	boil.SetDB(oldDB)
	db.Close()
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

type args struct {
	path string
}

func TestStartVisits(t *testing.T) {
	tests := getTestCases()

	for _, tt := range tests {
		mockDependencies(tt.name)

		t.Run(tt.name, func(t *testing.T) {
			expectRedisCommand(tt.args.path, tt.want)
			err := StartVisits(tt.args.path, time.Second)

			verifyError(t, err, tt.wantErr)
		})
	}
}

func getTestCases() []struct {
	name    string
	args    args
	want    int
	wantErr bool
	errMsg  error
} {
	return []struct {
		name    string
		args    args
		want    int
		wantErr bool
		errMsg  error
	}{
		{
			name: ErrorConnDo,
			args: args{
				path: "test",
			},
			wantErr: true,
			errMsg:  fmt.Errorf(ErrMsgFromConnDo),
		},
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
}

func mockDependencies(name string) {
	if name == ErrorRedisDial {
		patch := gomonkey.ApplyFunc(redisDial, func() (redigo.Conn, error) {
			return nil, fmt.Errorf(ErrMsgFromRedisDial)
		})
		defer patch.Reset()
	} else if name == ErrorConnDo {
		patch := gomonkey.ApplyMethodFunc(redigomock.NewConn(), "Do",
			func(commandName string, args ...interface{}) (reply interface{}, err error) {
				return nil, fmt.Errorf(ErrMsgFromConnDo)
			})
		defer patch.Reset()
	}
}

func expectRedisCommand(path string, want int) {
	conn.Command("SETEX", path, int(math.Ceil(time.Second.Seconds())), 1).Expect(want)
}

func verifyError(t *testing.T, err error, wantErr bool) {
	if (err != nil) != wantErr {
		t.Errorf("StartVisits() error = %v, wantErr %v", err, wantErr)
	}
}
