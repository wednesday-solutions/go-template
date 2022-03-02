package rediscache

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"
	"math"
	"reflect"
	"regexp"
	"testing"
	"time"

	"go-template/models"
	"go-template/testutls"

	"github.com/DATA-DOG/go-sqlmock"
	. "github.com/agiledragon/gomonkey/v2"
	"github.com/gomodule/redigo/redis"
	redigomock "github.com/rafaeljusto/redigomock/v3"
	"github.com/volatiletech/sqlboiler/v4/boil"
)

var conn = redigomock.NewConn()

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
	}{
		{
			name: "Success",
			args: args{
				userID:    testutls.MockID,
				dbQueries: []testutls.QueryData{},
			},
			want: testutls.MockUser(),
		},
		{
			name: "Success_WithCacheMiss",
			args: args{
				userID:    testutls.MockID,
				cacheMiss: true,
				dbQueries: []testutls.QueryData{
					{
						Actions: &[]driver.Value{testutls.MockID},
						Query:   "select * from \"users\" where \"id\"=$1",
						DbResponse: sqlmock.NewRows([]string{
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
						),
					},
				},
			},
			want: testutls.MockUser(),
		},
	}

	ApplyFunc(redisDial, func() (redis.Conn, error) {
		return conn, nil
	})

	oldDB := boil.GetDB()
	mock, db, _ := testutls.SetupEnvAndDB(t, testutls.Parameters{EnvFileLocation: "../../../.env.local"})

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			if tt.args.cacheMiss {
				conn.Command("GET", fmt.Sprintf("user%d", tt.args.userID)).Expect(nil)

				b, _ := json.Marshal(tt.want)
				conn.Command("SET", fmt.Sprintf("user%d", tt.args.userID), string(b)).Expect(nil)
				for _, dbQuery := range tt.args.dbQueries {
					mock.ExpectQuery(regexp.QuoteMeta(dbQuery.Query)).
						WithArgs(*dbQuery.Actions...).
						WillReturnRows(dbQuery.DbResponse)
				}
			} else {
				b, _ := json.Marshal(tt.want)
				conn.Command("GET", fmt.Sprintf("user%d", tt.args.userID)).Expect(b)
			}

			got, err := GetUser(tt.args.userID)
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

func TestGetRole(t *testing.T) {
	type args struct {
		roleID    int
		cacheMiss bool
		dbQueries []testutls.QueryData
	}
	var role = &models.Role{
		ID:          1,
		AccessLevel: 100,
		Name:        "Admin",
	}
	tests := []struct {
		name    string
		args    args
		want    *models.Role
		wantErr bool
	}{
		{
			name: "Success",
			args: args{
				roleID:    testutls.MockID,
				dbQueries: []testutls.QueryData{},
			},
			want: role,
		},
		{
			name: "Success_WithCacheMiss",
			args: args{
				roleID:    testutls.MockID,
				cacheMiss: true,
				dbQueries: []testutls.QueryData{
					{
						Actions: &[]driver.Value{role.ID},
						Query:   "select * from \"roles\" where \"id\"=$1",
						DbResponse: sqlmock.NewRows([]string{
							"id", "access_level", "name",
						}).AddRow(
							role.ID,
							role.AccessLevel,
							role.Name,
						),
					},
				},
			},
			want: role,
		},
	}
	ApplyFunc(redisDial, func() (redis.Conn, error) {
		return conn, nil
	})

	oldDB := boil.GetDB()
	mock, db, _ := testutls.SetupEnvAndDB(t, testutls.Parameters{EnvFileLocation: "../../../.env.local"})

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			if tt.args.cacheMiss {
				conn.Command("GET", fmt.Sprintf("role%d", tt.args.roleID)).Expect(nil)

				b, _ := json.Marshal(tt.want)
				conn.Command("SET", fmt.Sprintf("role%d", tt.args.roleID), string(b)).Expect(nil)
				for _, dbQuery := range tt.args.dbQueries {
					mock.ExpectQuery(regexp.QuoteMeta(dbQuery.Query)).
						WithArgs(*dbQuery.Actions...).
						WillReturnRows(dbQuery.DbResponse)
				}
			} else {
				b, _ := json.Marshal(tt.want)
				conn.Command("GET", fmt.Sprintf("role%d", tt.args.roleID)).Expect(b)
			}
			got, err := GetRole(tt.args.roleID)
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
	}{
		{
			name: "Success",
			args: args{
				path: "test",
			},
			want: 10,
		},
	}

	ApplyFunc(redisDial, func() (redis.Conn, error) {
		return conn, nil
	})
	for _, tt := range tests {
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

func TestStartVisits(t *testing.T) {
	type args struct {
		path string
	}
	tests := []struct {
		name    string
		args    args
		want    int
		wantErr bool
	}{
		{
			name: "Success",
			args: args{
				path: "test",
			},
			want: 10,
		},
	}
	ApplyFunc(redisDial, func() (redis.Conn, error) {
		return conn, nil
	})

	testutls.SetupEnvAndDB(t, testutls.Parameters{EnvFileLocation: "../../../.env.local"}) //nolint
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			conn.Command("SETEX", tt.args.path, int(math.Ceil(time.Second.Seconds())), 1).Expect(1)
			err := StartVisits(tt.args.path, time.Second)

			if (err != nil) != tt.wantErr {
				t.Errorf("StartVisits() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

		})
	}
}
