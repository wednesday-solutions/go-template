package resolver_test

import (
	"context"
	"database/sql/driver"
	"fmt"
	"go-template/daos"
	fm "go-template/gqlmodels"
	"go-template/internal/config"
	"go-template/models"
	"go-template/pkg/utl/convert"
	"go-template/pkg/utl/throttle"
	"go-template/resolver"
	"go-template/testutls"
	"log"
	"regexp"
	"testing"
	"time"

	"github.com/agiledragon/gomonkey/v2"
	"github.com/joho/godotenv"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"github.com/volatiletech/sqlboiler/v4/boil"
)

const (
	ErrorFromCreateAuthor = "Error Creating Author"
	ErrorFindingAuthor    = "Error Finding Author"
	ErrorDeleteAuthor     = "Error Delete Author"
	ErrorUpdateAuthor     = "Error Update Author"
)

func TestCreateAuthor(
	t *testing.T,
) {
	cases := []struct {
		name     string
		req      fm.AuthorCreateInput
		wantResp *fm.Author
		wantErr  bool
	}{
		{
			name:    ErrorFromCreateAuthor,
			req:     fm.AuthorCreateInput{},
			wantErr: true,
		},
		{
			name:    ErrorFromThrottleCheck,
			req:     fm.AuthorCreateInput{},
			wantErr: true,
		},
		{
			name:    ErrorFromConfig,
			req:     fm.AuthorCreateInput{},
			wantErr: true,
		},
		{
			name: SuccessCase,
			req: fm.AuthorCreateInput{
				FirstName: testutls.MockAuthor().FirstName.String,
				LastName:  testutls.MockAuthor().LastName.String,
				Email:     testutls.MockAuthor().Email.String,
			},
			wantResp: &fm.Author{
				ID:        fmt.Sprint(testutls.MockAuthor().ID),
				Email:     convert.NullDotStringToPointerString(testutls.MockAuthor().Email),
				FirstName: convert.NullDotStringToPointerString(testutls.MockAuthor().FirstName),
				LastName:  convert.NullDotStringToPointerString(testutls.MockAuthor().LastName),
				DeletedAt: convert.NullDotTimeToPointerInt(testutls.MockAuthor().DeletedAt),
				UpdatedAt: convert.NullDotTimeToPointerInt(testutls.MockAuthor().UpdatedAt),
			},
			wantErr: false,
		},
	}

	resolver1 := resolver.Resolver{}
	for _, tt := range cases {
		t.Run(
			tt.name,
			func(t *testing.T) {

				if tt.name == ErrorFromThrottleCheck {
					patch := gomonkey.ApplyFunc(throttle.Check, func(ctx context.Context, limit int, dur time.Duration) error {
						return fmt.Errorf("Internal error")
					})
					defer patch.Reset()
				}

				if tt.name == ErrorFromConfig {
					patch := gomonkey.ApplyFunc(config.Load, func() (*config.Configuration, error) {
						return nil, fmt.Errorf("error in loading config")
					})
					defer patch.Reset()

				}

				err := config.LoadEnvWithFilePrefix(convert.StringToPointerString("./../"))
				if err != nil {
					log.Fatal(err)
				}
				mock, db, _ := testutls.SetupMockDB(t)
				oldDB := boil.GetDB()
				defer func() {
					db.Close()
					boil.SetDB(oldDB)
				}()
				boil.SetDB(db)

				if tt.name == ErrorFromCreateAuthor {
					// insert new Author
					mock.ExpectQuery(regexp.QuoteMeta(`INSERT INTO "authors"`)).
						WithArgs().
						WillReturnError(fmt.Errorf(""))
				}
				// insert new Author
				rows := sqlmock.NewRows([]string{
					"id",
				}).
					AddRow(
						testutls.MockAuthor().ID,
					)
				mock.ExpectQuery(regexp.QuoteMeta(`INSERT INTO "authors"`)).
					WithArgs(
						testutls.MockAuthor().FirstName,
						testutls.MockAuthor().LastName,
						"",
						testutls.MockAuthor().Email,
						AnyTime{},
						AnyTime{},
					).
					WillReturnRows(rows)

				c := context.Background()
				response, err := resolver1.Mutation().
					CreateAuthor(c, tt.req)
				if tt.wantResp != nil {
					assert.Equal(t, tt.wantResp, response)
				}
				assert.Equal(t, tt.wantErr, err != nil)
			},
		)
	}
}

func TestUpdateAuthor(
	t *testing.T,
) {
	cases := []struct {
		name     string
		req      *fm.AuthorUpdateInput
		wantResp *fm.Author
		wantErr  bool
	}{
		{
			name:    ErrorFindingAuthor,
			req:     &fm.AuthorUpdateInput{},
			wantErr: true,
		},
		{
			name: ErrorUpdateAuthor,
			req: &fm.AuthorUpdateInput{
				FirstName: &testutls.MockAuthor().FirstName.String,
				LastName:  &testutls.MockAuthor().LastName.String,
				Email:     &testutls.MockAuthor().Email.String,
			},
			wantErr: true,
		},
		{
			name: SuccessCase,
			req: &fm.AuthorUpdateInput{
				FirstName: &testutls.MockAuthor().FirstName.String,
				LastName:  &testutls.MockAuthor().LastName.String,
				Email:     &testutls.MockAuthor().Email.String,
			},
			wantResp: &fm.Author{
				ID:        "0",
				FirstName: &testutls.MockAuthor().FirstName.String,
				LastName:  &testutls.MockAuthor().LastName.String,
				Email:     &testutls.MockAuthor().LastName.String,
			},
			wantErr: false,
		},
	}

	resolver1 := resolver.Resolver{}
	for _, tt := range cases {
		t.Run(
			tt.name,
			func(t *testing.T) {

				if tt.name == ErrorUpdateAuthor {

					patch := gomonkey.ApplyFunc(daos.UpdateAuthor,
						func(author models.Author, ctx context.Context) (models.Author, error) {
							return author, fmt.Errorf("error for update Author")
						})
					defer patch.Reset()
				}
				err := config.LoadEnvWithFilePrefix(convert.StringToPointerString("./../"))
				if err != nil {
					log.Fatal(err)
				}
				mock, db, _ := testutls.SetupMockDB(t)
				oldDB := boil.GetDB()
				defer func() {
					db.Close()
					boil.SetDB(oldDB)
				}()
				boil.SetDB(db)

				if tt.name == ErrorFindingAuthor {
					mock.ExpectQuery(regexp.QuoteMeta(`UPDATE "authors"`)).WithArgs().WillReturnError(fmt.Errorf(""))
				}

				rows := sqlmock.NewRows([]string{"first_name"}).AddRow(testutls.MockAuthor().FirstName)
				mock.ExpectQuery(regexp.QuoteMeta(`select * from "authors"`)).WithArgs(0).WillReturnRows(rows)

				// update Authors with new information
				result := driver.Result(driver.RowsAffected(1))
				mock.ExpectExec(regexp.QuoteMeta(`UPDATE "authors"`)).WillReturnResult(result)

				c := context.Background()
				ctx := context.WithValue(c, testutls.AuthorKey, testutls.MockAuthor())
				response, err := resolver1.Mutation().UpdateAuthor(ctx, *tt.req)
				if tt.wantResp != nil &&
					response != nil {
					assert.Equal(t, tt.wantResp, response)
				}
				assert.Equal(t, tt.wantErr, err != nil)
			},
		)
	}
}

func TestDeleteAuthor(
	t *testing.T,
) {
	cases := []struct {
		name     string
		wantResp *fm.AuthorDeletePayload
		wantErr  bool
	}{
		{
			name:    ErrorFindingAuthor,
			wantErr: true,
		},
		{
			name:    ErrorDeleteAuthor,
			wantErr: true,
		},
		{
			name: SuccessCase,
			wantResp: &fm.AuthorDeletePayload{
				ID: "0",
			},
			wantErr: false,
		},
	}

	resolver1 := resolver.Resolver{}
	for _, tt := range cases {
		t.Run(
			tt.name,
			func(t *testing.T) {
				if tt.name == ErrorDeleteAuthor {

					patch := gomonkey.ApplyFunc(daos.DeleteAuthor,
						func(author models.Author, ctx context.Context) (int64, error) {
							return 0, fmt.Errorf("error for delete author")
						})
					defer patch.Reset()
				}

				err := godotenv.Load(
					"../.env.local",
				)
				if err != nil {
					fmt.Print("error loading .env file")
				}
				db, mock, err := sqlmock.New()
				if err != nil {
					t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
				}
				oldDB := boil.GetDB()
				defer func() {
					db.Close()
					boil.SetDB(oldDB)
				}()
				boil.SetDB(db)

				if tt.name == ErrorFindingAuthor {
					mock.ExpectQuery(regexp.QuoteMeta(`select * from "authors" where "id"=$1`)).
						WithArgs().
						WillReturnError(fmt.Errorf(""))
				}
				// get Author by id
				rows := sqlmock.NewRows([]string{"id"}).
					AddRow(1)
				mock.ExpectQuery(regexp.QuoteMeta(`select * from "authors" where "id"=$1`)).
					WithArgs().
					WillReturnRows(rows)
				// delete Author
				result := driver.Result(driver.RowsAffected(1))
				mock.ExpectExec(regexp.QuoteMeta(`DELETE FROM "authors" WHERE "id"=$1`)).
					WillReturnResult(result)

				c := context.Background()
				ctx := context.WithValue(c, testutls.AuthorKey, testutls.MockAuthor())
				response, err := resolver1.Mutation().
					DeleteAuthor(ctx, fm.AuthorDeleteInput{ID: "1"})
				if tt.wantResp != nil {
					assert.Equal(t, tt.wantResp, response)
				}
				assert.Equal(t, tt.wantErr, err != nil)
			},
		)
	}
}
