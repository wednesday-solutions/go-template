package daos_test

import (
	"fmt"
	"regexp"
	"testing"

	"go-template/daos"
	"go-template/models"
	"go-template/testutls"

	"github.com/DATA-DOG/go-sqlmock"
	_ "github.com/go-sql-driver/mysql"
	"github.com/joho/godotenv"
	"github.com/stretchr/testify/assert"
	"github.com/volatiletech/sqlboiler/v4/boil"
)

func TestCreateRoleTx(t *testing.T) {

	cases := []struct {
		name string
		req  models.Role
		err  error
	}{
		{
			name: "Passing role type value",
			req:  models.Role{},
			err:  nil,
		},
	}

	for _, tt := range cases {
		err := godotenv.Load("../.env.local")
		if err != nil {
			fmt.Print("error loading .env file")
		}

		db, mock, err := sqlmock.New()
		if err != nil {
			t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
		}
		// Inject mock instance into boil.
		oldDB := boil.GetDB()
		defer func() {
			db.Close()
			boil.SetDB(oldDB)
		}()
		boil.SetDB(db)

		query := regexp.QuoteMeta("INSERT INTO `roles` (`access_level`,`name`,`created_at`,`updated_at`,`deleted_at`)" +
			" VALUES (?,?,?,?,?)")
		mock.ExpectExec(query).
			WithArgs(tt.req.AccessLevel, tt.req.Name, testutls.AnyTime{}, testutls.AnyTime{}, nil).
			WillReturnResult(sqlmock.NewResult(1, 1))

		t.Run(tt.name, func(t *testing.T) {
			res, err := daos.CreateRoleTx(tt.req, nil)
			if err != nil {
				assert.Equal(t, true, tt.err != nil)
			} else {
				assert.Equal(t, err, tt.err)
				assert.Equal(t, true, tt.req.ID != res.ID)
			}
		})
	}
}

func TestFindRoleByID(t *testing.T) {

	cases := []struct {
		name string
		req  int
		err  error
	}{
		{
			name: "Passing a user_id",
			req:  1,
			err:  nil,
		},
	}

	for _, tt := range cases {
		err := godotenv.Load("../.env.local")
		if err != nil {
			fmt.Print("error loading .env file")
		}

		db, mock, err := sqlmock.New()
		if err != nil {
			t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
		}
		// Inject mock instance into boil.
		oldDB := boil.GetDB()
		defer func() {
			db.Close()
			boil.SetDB(oldDB)
		}()
		boil.SetDB(db)

		rows := sqlmock.NewRows([]string{"id"}).AddRow(1)
		query := regexp.
			QuoteMeta("select * from `roles` where `id`=?")
		mock.ExpectQuery(query).
			WithArgs(tt.req).
			WillReturnRows(rows)

		t.Run(tt.name, func(t *testing.T) {
			_, err := daos.FindRoleByID(tt.req)
			assert.Equal(t, err, tt.err)

		})
	}
}
