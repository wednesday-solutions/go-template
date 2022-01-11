package rediscache

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"
	"reflect"
	"regexp"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	. "github.com/agiledragon/gomonkey/v2"
	redigo "github.com/gomodule/redigo/redis"
	redis "github.com/gomodule/redigo/redis"
	redigomock "github.com/rafaeljusto/redigomock/v3"
	"github.com/volatiletech/sqlboiler/boil"
	"github.com/wednesday-solutions/go-template/models"
	"github.com/wednesday-solutions/go-template/testutls"
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
							"first_name", "last_name", "username", "email", "mobile", "phone", "address",
						}).AddRow(
							testutls.MockUser().FirstName,
							testutls.MockUser().LastName,
							testutls.MockUser().Username,
							testutls.MockUser().Email,
							testutls.MockUser().Mobile,
							testutls.MockUser().Phone,
							testutls.MockUser().Address,
						),
					},
				},
			},
			want: testutls.MockUser(),
		},
	}
	conn := redigomock.NewConn()
	ApplyFunc(redigo.Dial, func(string, string, ...redis.DialOption) (redis.Conn, error) {
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
