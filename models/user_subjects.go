// Code generated by SQLBoiler 3.7.1 (https://github.com/volatiletech/sqlboiler). DO NOT EDIT.
// This file is meant to be re-generated in place and/or deleted at any time.

package models

import (
	"context"
	"database/sql"
	"fmt"
	"reflect"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/friendsofgo/errors"
	"github.com/volatiletech/null"
	"github.com/volatiletech/sqlboiler/boil"
	"github.com/volatiletech/sqlboiler/queries"
	"github.com/volatiletech/sqlboiler/queries/qm"
	"github.com/volatiletech/sqlboiler/queries/qmhelper"
	"github.com/volatiletech/sqlboiler/strmangle"
)

// UserSubject is an object representing the database table.
type UserSubject struct {
	ID        int       `boil:"id" json:"id" toml:"id" yaml:"id"`
	SubjectID null.Int  `boil:"subject_id" json:"subject_id,omitempty" toml:"subject_id" yaml:"subject_id,omitempty"`
	UserID    null.Int  `boil:"user_id" json:"user_id,omitempty" toml:"user_id" yaml:"user_id,omitempty"`
	CreatedAt null.Time `boil:"created_at" json:"created_at,omitempty" toml:"created_at" yaml:"created_at,omitempty"`
	UpdatedAt null.Time `boil:"updated_at" json:"updated_at,omitempty" toml:"updated_at" yaml:"updated_at,omitempty"`
	DeletedAt null.Time `boil:"deleted_at" json:"deleted_at,omitempty" toml:"deleted_at" yaml:"deleted_at,omitempty"`

	R *userSubjectR `boil:"-" json:"-" toml:"-" yaml:"-"`
	L userSubjectL  `boil:"-" json:"-" toml:"-" yaml:"-"`
}

var UserSubjectColumns = struct {
	ID        string
	SubjectID string
	UserID    string
	CreatedAt string
	UpdatedAt string
	DeletedAt string
}{
	ID:        "id",
	SubjectID: "subject_id",
	UserID:    "user_id",
	CreatedAt: "created_at",
	UpdatedAt: "updated_at",
	DeletedAt: "deleted_at",
}

// Generated where

type whereHelpernull_Int struct{ field string }

func (w whereHelpernull_Int) EQ(x null.Int) qm.QueryMod {
	return qmhelper.WhereNullEQ(w.field, false, x)
}
func (w whereHelpernull_Int) NEQ(x null.Int) qm.QueryMod {
	return qmhelper.WhereNullEQ(w.field, true, x)
}
func (w whereHelpernull_Int) IsNull() qm.QueryMod    { return qmhelper.WhereIsNull(w.field) }
func (w whereHelpernull_Int) IsNotNull() qm.QueryMod { return qmhelper.WhereIsNotNull(w.field) }
func (w whereHelpernull_Int) LT(x null.Int) qm.QueryMod {
	return qmhelper.Where(w.field, qmhelper.LT, x)
}
func (w whereHelpernull_Int) LTE(x null.Int) qm.QueryMod {
	return qmhelper.Where(w.field, qmhelper.LTE, x)
}
func (w whereHelpernull_Int) GT(x null.Int) qm.QueryMod {
	return qmhelper.Where(w.field, qmhelper.GT, x)
}
func (w whereHelpernull_Int) GTE(x null.Int) qm.QueryMod {
	return qmhelper.Where(w.field, qmhelper.GTE, x)
}

var UserSubjectWhere = struct {
	ID        whereHelperint
	SubjectID whereHelpernull_Int
	UserID    whereHelpernull_Int
	CreatedAt whereHelpernull_Time
	UpdatedAt whereHelpernull_Time
	DeletedAt whereHelpernull_Time
}{
	ID:        whereHelperint{field: "\"user_subjects\".\"id\""},
	SubjectID: whereHelpernull_Int{field: "\"user_subjects\".\"subject_id\""},
	UserID:    whereHelpernull_Int{field: "\"user_subjects\".\"user_id\""},
	CreatedAt: whereHelpernull_Time{field: "\"user_subjects\".\"created_at\""},
	UpdatedAt: whereHelpernull_Time{field: "\"user_subjects\".\"updated_at\""},
	DeletedAt: whereHelpernull_Time{field: "\"user_subjects\".\"deleted_at\""},
}

// UserSubjectRels is where relationship names are stored.
var UserSubjectRels = struct {
	Subject string
	User    string
}{
	Subject: "Subject",
	User:    "User",
}

// userSubjectR is where relationships are stored.
type userSubjectR struct {
	Subject *Subject
	User    *User
}

// NewStruct creates a new relationship struct
func (*userSubjectR) NewStruct() *userSubjectR {
	return &userSubjectR{}
}

// userSubjectL is where Load methods for each relationship are stored.
type userSubjectL struct{}

var (
	userSubjectAllColumns            = []string{"id", "subject_id", "user_id", "created_at", "updated_at", "deleted_at"}
	userSubjectColumnsWithoutDefault = []string{"subject_id", "user_id", "created_at", "deleted_at"}
	userSubjectColumnsWithDefault    = []string{"id", "updated_at"}
	userSubjectPrimaryKeyColumns     = []string{"id"}
)

type (
	// UserSubjectSlice is an alias for a slice of pointers to UserSubject.
	// This should generally be used opposed to []UserSubject.
	UserSubjectSlice []*UserSubject

	userSubjectQuery struct {
		*queries.Query
	}
)

// Cache for insert, update and upsert
var (
	userSubjectType                 = reflect.TypeOf(&UserSubject{})
	userSubjectMapping              = queries.MakeStructMapping(userSubjectType)
	userSubjectPrimaryKeyMapping, _ = queries.BindMapping(userSubjectType, userSubjectMapping, userSubjectPrimaryKeyColumns)
	userSubjectInsertCacheMut       sync.RWMutex
	userSubjectInsertCache          = make(map[string]insertCache)
	userSubjectUpdateCacheMut       sync.RWMutex
	userSubjectUpdateCache          = make(map[string]updateCache)
	userSubjectUpsertCacheMut       sync.RWMutex
	userSubjectUpsertCache          = make(map[string]insertCache)
)

var (
	// Force time package dependency for automated UpdatedAt/CreatedAt.
	_ = time.Second
	// Force qmhelper dependency for where clause generation (which doesn't
	// always happen)
	_ = qmhelper.Where
)

// One returns a single userSubject record from the query.
func (q userSubjectQuery) One(ctx context.Context, exec boil.ContextExecutor) (*UserSubject, error) {
	o := &UserSubject{}

	queries.SetLimit(q.Query, 1)

	err := q.Bind(ctx, exec, o)
	if err != nil {
		if errors.Cause(err) == sql.ErrNoRows {
			return nil, sql.ErrNoRows
		}
		return nil, errors.Wrap(err, "models: failed to execute a one query for user_subjects")
	}

	return o, nil
}

// All returns all UserSubject records from the query.
func (q userSubjectQuery) All(ctx context.Context, exec boil.ContextExecutor) (UserSubjectSlice, error) {
	var o []*UserSubject

	err := q.Bind(ctx, exec, &o)
	if err != nil {
		return nil, errors.Wrap(err, "models: failed to assign all query results to UserSubject slice")
	}

	return o, nil
}

// Count returns the count of all UserSubject records in the query.
func (q userSubjectQuery) Count(ctx context.Context, exec boil.ContextExecutor) (int64, error) {
	var count int64

	queries.SetSelect(q.Query, nil)
	queries.SetCount(q.Query)

	err := q.Query.QueryRowContext(ctx, exec).Scan(&count)
	if err != nil {
		return 0, errors.Wrap(err, "models: failed to count user_subjects rows")
	}

	return count, nil
}

// Exists checks if the row exists in the table.
func (q userSubjectQuery) Exists(ctx context.Context, exec boil.ContextExecutor) (bool, error) {
	var count int64

	queries.SetSelect(q.Query, nil)
	queries.SetCount(q.Query)
	queries.SetLimit(q.Query, 1)

	err := q.Query.QueryRowContext(ctx, exec).Scan(&count)
	if err != nil {
		return false, errors.Wrap(err, "models: failed to check if user_subjects exists")
	}

	return count > 0, nil
}

// Subject pointed to by the foreign key.
func (o *UserSubject) Subject(mods ...qm.QueryMod) subjectQuery {
	queryMods := []qm.QueryMod{
		qm.Where("\"id\" = ?", o.SubjectID),
	}

	queryMods = append(queryMods, mods...)

	query := Subjects(queryMods...)
	queries.SetFrom(query.Query, "\"subjects\"")

	return query
}

// User pointed to by the foreign key.
func (o *UserSubject) User(mods ...qm.QueryMod) userQuery {
	queryMods := []qm.QueryMod{
		qm.Where("\"id\" = ?", o.UserID),
	}

	queryMods = append(queryMods, mods...)

	query := Users(queryMods...)
	queries.SetFrom(query.Query, "\"users\"")

	return query
}

// LoadSubject allows an eager lookup of values, cached into the
// loaded structs of the objects. This is for an N-1 relationship.
func (userSubjectL) LoadSubject(ctx context.Context, e boil.ContextExecutor, singular bool, maybeUserSubject interface{}, mods queries.Applicator) error {
	var slice []*UserSubject
	var object *UserSubject

	if singular {
		object = maybeUserSubject.(*UserSubject)
	} else {
		slice = *maybeUserSubject.(*[]*UserSubject)
	}

	args := make([]interface{}, 0, 1)
	if singular {
		if object.R == nil {
			object.R = &userSubjectR{}
		}
		if !queries.IsNil(object.SubjectID) {
			args = append(args, object.SubjectID)
		}

	} else {
	Outer:
		for _, obj := range slice {
			if obj.R == nil {
				obj.R = &userSubjectR{}
			}

			for _, a := range args {
				if queries.Equal(a, obj.SubjectID) {
					continue Outer
				}
			}

			if !queries.IsNil(obj.SubjectID) {
				args = append(args, obj.SubjectID)
			}

		}
	}

	if len(args) == 0 {
		return nil
	}

	query := NewQuery(qm.From(`subjects`), qm.WhereIn(`subjects.id in ?`, args...))
	if mods != nil {
		mods.Apply(query)
	}

	results, err := query.QueryContext(ctx, e)
	if err != nil {
		return errors.Wrap(err, "failed to eager load Subject")
	}

	var resultSlice []*Subject
	if err = queries.Bind(results, &resultSlice); err != nil {
		return errors.Wrap(err, "failed to bind eager loaded slice Subject")
	}

	if err = results.Close(); err != nil {
		return errors.Wrap(err, "failed to close results of eager load for subjects")
	}
	if err = results.Err(); err != nil {
		return errors.Wrap(err, "error occurred during iteration of eager loaded relations for subjects")
	}

	if len(resultSlice) == 0 {
		return nil
	}

	if singular {
		foreign := resultSlice[0]
		object.R.Subject = foreign
		if foreign.R == nil {
			foreign.R = &subjectR{}
		}
		foreign.R.UserSubjects = append(foreign.R.UserSubjects, object)
		return nil
	}

	for _, local := range slice {
		for _, foreign := range resultSlice {
			if queries.Equal(local.SubjectID, foreign.ID) {
				local.R.Subject = foreign
				if foreign.R == nil {
					foreign.R = &subjectR{}
				}
				foreign.R.UserSubjects = append(foreign.R.UserSubjects, local)
				break
			}
		}
	}

	return nil
}

// LoadUser allows an eager lookup of values, cached into the
// loaded structs of the objects. This is for an N-1 relationship.
func (userSubjectL) LoadUser(ctx context.Context, e boil.ContextExecutor, singular bool, maybeUserSubject interface{}, mods queries.Applicator) error {
	var slice []*UserSubject
	var object *UserSubject

	if singular {
		object = maybeUserSubject.(*UserSubject)
	} else {
		slice = *maybeUserSubject.(*[]*UserSubject)
	}

	args := make([]interface{}, 0, 1)
	if singular {
		if object.R == nil {
			object.R = &userSubjectR{}
		}
		if !queries.IsNil(object.UserID) {
			args = append(args, object.UserID)
		}

	} else {
	Outer:
		for _, obj := range slice {
			if obj.R == nil {
				obj.R = &userSubjectR{}
			}

			for _, a := range args {
				if queries.Equal(a, obj.UserID) {
					continue Outer
				}
			}

			if !queries.IsNil(obj.UserID) {
				args = append(args, obj.UserID)
			}

		}
	}

	if len(args) == 0 {
		return nil
	}

	query := NewQuery(qm.From(`users`), qm.WhereIn(`users.id in ?`, args...))
	if mods != nil {
		mods.Apply(query)
	}

	results, err := query.QueryContext(ctx, e)
	if err != nil {
		return errors.Wrap(err, "failed to eager load User")
	}

	var resultSlice []*User
	if err = queries.Bind(results, &resultSlice); err != nil {
		return errors.Wrap(err, "failed to bind eager loaded slice User")
	}

	if err = results.Close(); err != nil {
		return errors.Wrap(err, "failed to close results of eager load for users")
	}
	if err = results.Err(); err != nil {
		return errors.Wrap(err, "error occurred during iteration of eager loaded relations for users")
	}

	if len(resultSlice) == 0 {
		return nil
	}

	if singular {
		foreign := resultSlice[0]
		object.R.User = foreign
		if foreign.R == nil {
			foreign.R = &userR{}
		}
		foreign.R.UserSubjects = append(foreign.R.UserSubjects, object)
		return nil
	}

	for _, local := range slice {
		for _, foreign := range resultSlice {
			if queries.Equal(local.UserID, foreign.ID) {
				local.R.User = foreign
				if foreign.R == nil {
					foreign.R = &userR{}
				}
				foreign.R.UserSubjects = append(foreign.R.UserSubjects, local)
				break
			}
		}
	}

	return nil
}

// SetSubject of the userSubject to the related item.
// Sets o.R.Subject to related.
// Adds o to related.R.UserSubjects.
func (o *UserSubject) SetSubject(ctx context.Context, exec boil.ContextExecutor, insert bool, related *Subject) error {
	var err error
	if insert {
		if err = related.Insert(ctx, exec, boil.Infer()); err != nil {
			return errors.Wrap(err, "failed to insert into foreign table")
		}
	}

	updateQuery := fmt.Sprintf(
		"UPDATE \"user_subjects\" SET %s WHERE %s",
		strmangle.SetParamNames("\"", "\"", 1, []string{"subject_id"}),
		strmangle.WhereClause("\"", "\"", 2, userSubjectPrimaryKeyColumns),
	)
	values := []interface{}{related.ID, o.ID}

	if boil.IsDebug(ctx) {
		writer := boil.DebugWriterFrom(ctx)
		fmt.Fprintln(writer, updateQuery)
		fmt.Fprintln(writer, values)
	}
	if _, err = exec.ExecContext(ctx, updateQuery, values...); err != nil {
		return errors.Wrap(err, "failed to update local table")
	}

	queries.Assign(&o.SubjectID, related.ID)
	if o.R == nil {
		o.R = &userSubjectR{
			Subject: related,
		}
	} else {
		o.R.Subject = related
	}

	if related.R == nil {
		related.R = &subjectR{
			UserSubjects: UserSubjectSlice{o},
		}
	} else {
		related.R.UserSubjects = append(related.R.UserSubjects, o)
	}

	return nil
}

// RemoveSubject relationship.
// Sets o.R.Subject to nil.
// Removes o from all passed in related items' relationships struct (Optional).
func (o *UserSubject) RemoveSubject(ctx context.Context, exec boil.ContextExecutor, related *Subject) error {
	var err error

	queries.SetScanner(&o.SubjectID, nil)
	if _, err = o.Update(ctx, exec, boil.Whitelist("subject_id")); err != nil {
		return errors.Wrap(err, "failed to update local table")
	}

	if o.R != nil {
		o.R.Subject = nil
	}
	if related == nil || related.R == nil {
		return nil
	}

	for i, ri := range related.R.UserSubjects {
		if queries.Equal(o.SubjectID, ri.SubjectID) {
			continue
		}

		ln := len(related.R.UserSubjects)
		if ln > 1 && i < ln-1 {
			related.R.UserSubjects[i] = related.R.UserSubjects[ln-1]
		}
		related.R.UserSubjects = related.R.UserSubjects[:ln-1]
		break
	}
	return nil
}

// SetUser of the userSubject to the related item.
// Sets o.R.User to related.
// Adds o to related.R.UserSubjects.
func (o *UserSubject) SetUser(ctx context.Context, exec boil.ContextExecutor, insert bool, related *User) error {
	var err error
	if insert {
		if err = related.Insert(ctx, exec, boil.Infer()); err != nil {
			return errors.Wrap(err, "failed to insert into foreign table")
		}
	}

	updateQuery := fmt.Sprintf(
		"UPDATE \"user_subjects\" SET %s WHERE %s",
		strmangle.SetParamNames("\"", "\"", 1, []string{"user_id"}),
		strmangle.WhereClause("\"", "\"", 2, userSubjectPrimaryKeyColumns),
	)
	values := []interface{}{related.ID, o.ID}

	if boil.IsDebug(ctx) {
		writer := boil.DebugWriterFrom(ctx)
		fmt.Fprintln(writer, updateQuery)
		fmt.Fprintln(writer, values)
	}
	if _, err = exec.ExecContext(ctx, updateQuery, values...); err != nil {
		return errors.Wrap(err, "failed to update local table")
	}

	queries.Assign(&o.UserID, related.ID)
	if o.R == nil {
		o.R = &userSubjectR{
			User: related,
		}
	} else {
		o.R.User = related
	}

	if related.R == nil {
		related.R = &userR{
			UserSubjects: UserSubjectSlice{o},
		}
	} else {
		related.R.UserSubjects = append(related.R.UserSubjects, o)
	}

	return nil
}

// RemoveUser relationship.
// Sets o.R.User to nil.
// Removes o from all passed in related items' relationships struct (Optional).
func (o *UserSubject) RemoveUser(ctx context.Context, exec boil.ContextExecutor, related *User) error {
	var err error

	queries.SetScanner(&o.UserID, nil)
	if _, err = o.Update(ctx, exec, boil.Whitelist("user_id")); err != nil {
		return errors.Wrap(err, "failed to update local table")
	}

	if o.R != nil {
		o.R.User = nil
	}
	if related == nil || related.R == nil {
		return nil
	}

	for i, ri := range related.R.UserSubjects {
		if queries.Equal(o.UserID, ri.UserID) {
			continue
		}

		ln := len(related.R.UserSubjects)
		if ln > 1 && i < ln-1 {
			related.R.UserSubjects[i] = related.R.UserSubjects[ln-1]
		}
		related.R.UserSubjects = related.R.UserSubjects[:ln-1]
		break
	}
	return nil
}

// UserSubjects retrieves all the records using an executor.
func UserSubjects(mods ...qm.QueryMod) userSubjectQuery {
	mods = append(mods, qm.From("\"user_subjects\""))
	return userSubjectQuery{NewQuery(mods...)}
}

// FindUserSubject retrieves a single record by ID with an executor.
// If selectCols is empty Find will return all columns.
func FindUserSubject(ctx context.Context, exec boil.ContextExecutor, iD int, selectCols ...string) (*UserSubject, error) {
	userSubjectObj := &UserSubject{}

	sel := "*"
	if len(selectCols) > 0 {
		sel = strings.Join(strmangle.IdentQuoteSlice(dialect.LQ, dialect.RQ, selectCols), ",")
	}
	query := fmt.Sprintf(
		"select %s from \"user_subjects\" where \"id\"=$1", sel,
	)

	q := queries.Raw(query, iD)

	err := q.Bind(ctx, exec, userSubjectObj)
	if err != nil {
		if errors.Cause(err) == sql.ErrNoRows {
			return nil, sql.ErrNoRows
		}
		return nil, errors.Wrap(err, "models: unable to select from user_subjects")
	}

	return userSubjectObj, nil
}

// Insert a single record using an executor.
// See boil.Columns.InsertColumnSet documentation to understand column list inference for inserts.
func (o *UserSubject) Insert(ctx context.Context, exec boil.ContextExecutor, columns boil.Columns) error {
	if o == nil {
		return errors.New("models: no user_subjects provided for insertion")
	}

	var err error
	if !boil.TimestampsAreSkipped(ctx) {
		currTime := time.Now().In(boil.GetLocation())

		if queries.MustTime(o.CreatedAt).IsZero() {
			queries.SetScanner(&o.CreatedAt, currTime)
		}
		if queries.MustTime(o.UpdatedAt).IsZero() {
			queries.SetScanner(&o.UpdatedAt, currTime)
		}
	}

	nzDefaults := queries.NonZeroDefaultSet(userSubjectColumnsWithDefault, o)

	key := makeCacheKey(columns, nzDefaults)
	userSubjectInsertCacheMut.RLock()
	cache, cached := userSubjectInsertCache[key]
	userSubjectInsertCacheMut.RUnlock()

	if !cached {
		wl, returnColumns := columns.InsertColumnSet(
			userSubjectAllColumns,
			userSubjectColumnsWithDefault,
			userSubjectColumnsWithoutDefault,
			nzDefaults,
		)

		cache.valueMapping, err = queries.BindMapping(userSubjectType, userSubjectMapping, wl)
		if err != nil {
			return err
		}
		cache.retMapping, err = queries.BindMapping(userSubjectType, userSubjectMapping, returnColumns)
		if err != nil {
			return err
		}
		if len(wl) != 0 {
			cache.query = fmt.Sprintf("INSERT INTO \"user_subjects\" (\"%s\") %%sVALUES (%s)%%s", strings.Join(wl, "\",\""), strmangle.Placeholders(dialect.UseIndexPlaceholders, len(wl), 1, 1))
		} else {
			cache.query = "INSERT INTO \"user_subjects\" %sDEFAULT VALUES%s"
		}

		var queryOutput, queryReturning string

		if len(cache.retMapping) != 0 {
			queryReturning = fmt.Sprintf(" RETURNING \"%s\"", strings.Join(returnColumns, "\",\""))
		}

		cache.query = fmt.Sprintf(cache.query, queryOutput, queryReturning)
	}

	value := reflect.Indirect(reflect.ValueOf(o))
	vals := queries.ValuesFromMapping(value, cache.valueMapping)

	if boil.IsDebug(ctx) {
		writer := boil.DebugWriterFrom(ctx)
		fmt.Fprintln(writer, cache.query)
		fmt.Fprintln(writer, vals)
	}

	if len(cache.retMapping) != 0 {
		err = exec.QueryRowContext(ctx, cache.query, vals...).Scan(queries.PtrsFromMapping(value, cache.retMapping)...)
	} else {
		_, err = exec.ExecContext(ctx, cache.query, vals...)
	}

	if err != nil {
		return errors.Wrap(err, "models: unable to insert into user_subjects")
	}

	if !cached {
		userSubjectInsertCacheMut.Lock()
		userSubjectInsertCache[key] = cache
		userSubjectInsertCacheMut.Unlock()
	}

	return nil
}

// Update uses an executor to update the UserSubject.
// See boil.Columns.UpdateColumnSet documentation to understand column list inference for updates.
// Update does not automatically update the record in case of default values. Use .Reload() to refresh the records.
func (o *UserSubject) Update(ctx context.Context, exec boil.ContextExecutor, columns boil.Columns) (int64, error) {
	if !boil.TimestampsAreSkipped(ctx) {
		currTime := time.Now().In(boil.GetLocation())

		queries.SetScanner(&o.UpdatedAt, currTime)
	}

	var err error
	key := makeCacheKey(columns, nil)
	userSubjectUpdateCacheMut.RLock()
	cache, cached := userSubjectUpdateCache[key]
	userSubjectUpdateCacheMut.RUnlock()

	if !cached {
		wl := columns.UpdateColumnSet(
			userSubjectAllColumns,
			userSubjectPrimaryKeyColumns,
		)

		if !columns.IsWhitelist() {
			wl = strmangle.SetComplement(wl, []string{"created_at"})
		}
		if len(wl) == 0 {
			return 0, errors.New("models: unable to update user_subjects, could not build whitelist")
		}

		cache.query = fmt.Sprintf("UPDATE \"user_subjects\" SET %s WHERE %s",
			strmangle.SetParamNames("\"", "\"", 1, wl),
			strmangle.WhereClause("\"", "\"", len(wl)+1, userSubjectPrimaryKeyColumns),
		)
		cache.valueMapping, err = queries.BindMapping(userSubjectType, userSubjectMapping, append(wl, userSubjectPrimaryKeyColumns...))
		if err != nil {
			return 0, err
		}
	}

	values := queries.ValuesFromMapping(reflect.Indirect(reflect.ValueOf(o)), cache.valueMapping)

	if boil.IsDebug(ctx) {
		writer := boil.DebugWriterFrom(ctx)
		fmt.Fprintln(writer, cache.query)
		fmt.Fprintln(writer, values)
	}
	var result sql.Result
	result, err = exec.ExecContext(ctx, cache.query, values...)
	if err != nil {
		return 0, errors.Wrap(err, "models: unable to update user_subjects row")
	}

	rowsAff, err := result.RowsAffected()
	if err != nil {
		return 0, errors.Wrap(err, "models: failed to get rows affected by update for user_subjects")
	}

	if !cached {
		userSubjectUpdateCacheMut.Lock()
		userSubjectUpdateCache[key] = cache
		userSubjectUpdateCacheMut.Unlock()
	}

	return rowsAff, nil
}

// UpdateAll updates all rows with the specified column values.
func (q userSubjectQuery) UpdateAll(ctx context.Context, exec boil.ContextExecutor, cols M) (int64, error) {
	queries.SetUpdate(q.Query, cols)

	result, err := q.Query.ExecContext(ctx, exec)
	if err != nil {
		return 0, errors.Wrap(err, "models: unable to update all for user_subjects")
	}

	rowsAff, err := result.RowsAffected()
	if err != nil {
		return 0, errors.Wrap(err, "models: unable to retrieve rows affected for user_subjects")
	}

	return rowsAff, nil
}

// UpdateAll updates all rows with the specified column values, using an executor.
func (o UserSubjectSlice) UpdateAll(ctx context.Context, exec boil.ContextExecutor, cols M) (int64, error) {
	ln := int64(len(o))
	if ln == 0 {
		return 0, nil
	}

	if len(cols) == 0 {
		return 0, errors.New("models: update all requires at least one column argument")
	}

	colNames := make([]string, len(cols))
	args := make([]interface{}, len(cols))

	i := 0
	for name, value := range cols {
		colNames[i] = name
		args[i] = value
		i++
	}

	// Append all of the primary key values for each column
	for _, obj := range o {
		pkeyArgs := queries.ValuesFromMapping(reflect.Indirect(reflect.ValueOf(obj)), userSubjectPrimaryKeyMapping)
		args = append(args, pkeyArgs...)
	}

	sql := fmt.Sprintf("UPDATE \"user_subjects\" SET %s WHERE %s",
		strmangle.SetParamNames("\"", "\"", 1, colNames),
		strmangle.WhereClauseRepeated(string(dialect.LQ), string(dialect.RQ), len(colNames)+1, userSubjectPrimaryKeyColumns, len(o)))

	if boil.IsDebug(ctx) {
		writer := boil.DebugWriterFrom(ctx)
		fmt.Fprintln(writer, sql)
		fmt.Fprintln(writer, args...)
	}
	result, err := exec.ExecContext(ctx, sql, args...)
	if err != nil {
		return 0, errors.Wrap(err, "models: unable to update all in userSubject slice")
	}

	rowsAff, err := result.RowsAffected()
	if err != nil {
		return 0, errors.Wrap(err, "models: unable to retrieve rows affected all in update all userSubject")
	}
	return rowsAff, nil
}

// Upsert attempts an insert using an executor, and does an update or ignore on conflict.
// See boil.Columns documentation for how to properly use updateColumns and insertColumns.
func (o *UserSubject) Upsert(ctx context.Context, exec boil.ContextExecutor, updateOnConflict bool, conflictColumns []string, updateColumns, insertColumns boil.Columns) error {
	if o == nil {
		return errors.New("models: no user_subjects provided for upsert")
	}
	if !boil.TimestampsAreSkipped(ctx) {
		currTime := time.Now().In(boil.GetLocation())

		if queries.MustTime(o.CreatedAt).IsZero() {
			queries.SetScanner(&o.CreatedAt, currTime)
		}
		queries.SetScanner(&o.UpdatedAt, currTime)
	}

	nzDefaults := queries.NonZeroDefaultSet(userSubjectColumnsWithDefault, o)

	// Build cache key in-line uglily - mysql vs psql problems
	buf := strmangle.GetBuffer()
	if updateOnConflict {
		buf.WriteByte('t')
	} else {
		buf.WriteByte('f')
	}
	buf.WriteByte('.')
	for _, c := range conflictColumns {
		buf.WriteString(c)
	}
	buf.WriteByte('.')
	buf.WriteString(strconv.Itoa(updateColumns.Kind))
	for _, c := range updateColumns.Cols {
		buf.WriteString(c)
	}
	buf.WriteByte('.')
	buf.WriteString(strconv.Itoa(insertColumns.Kind))
	for _, c := range insertColumns.Cols {
		buf.WriteString(c)
	}
	buf.WriteByte('.')
	for _, c := range nzDefaults {
		buf.WriteString(c)
	}
	key := buf.String()
	strmangle.PutBuffer(buf)

	userSubjectUpsertCacheMut.RLock()
	cache, cached := userSubjectUpsertCache[key]
	userSubjectUpsertCacheMut.RUnlock()

	var err error

	if !cached {
		insert, ret := insertColumns.InsertColumnSet(
			userSubjectAllColumns,
			userSubjectColumnsWithDefault,
			userSubjectColumnsWithoutDefault,
			nzDefaults,
		)
		update := updateColumns.UpdateColumnSet(
			userSubjectAllColumns,
			userSubjectPrimaryKeyColumns,
		)

		if updateOnConflict && len(update) == 0 {
			return errors.New("models: unable to upsert user_subjects, could not build update column list")
		}

		conflict := conflictColumns
		if len(conflict) == 0 {
			conflict = make([]string, len(userSubjectPrimaryKeyColumns))
			copy(conflict, userSubjectPrimaryKeyColumns)
		}
		cache.query = buildUpsertQueryPostgres(dialect, "\"user_subjects\"", updateOnConflict, ret, update, conflict, insert)

		cache.valueMapping, err = queries.BindMapping(userSubjectType, userSubjectMapping, insert)
		if err != nil {
			return err
		}
		if len(ret) != 0 {
			cache.retMapping, err = queries.BindMapping(userSubjectType, userSubjectMapping, ret)
			if err != nil {
				return err
			}
		}
	}

	value := reflect.Indirect(reflect.ValueOf(o))
	vals := queries.ValuesFromMapping(value, cache.valueMapping)
	var returns []interface{}
	if len(cache.retMapping) != 0 {
		returns = queries.PtrsFromMapping(value, cache.retMapping)
	}

	if boil.IsDebug(ctx) {
		writer := boil.DebugWriterFrom(ctx)
		fmt.Fprintln(writer, cache.query)
		fmt.Fprintln(writer, vals)
	}
	if len(cache.retMapping) != 0 {
		err = exec.QueryRowContext(ctx, cache.query, vals...).Scan(returns...)
		if err == sql.ErrNoRows {
			err = nil // Postgres doesn't return anything when there's no update
		}
	} else {
		_, err = exec.ExecContext(ctx, cache.query, vals...)
	}
	if err != nil {
		return errors.Wrap(err, "models: unable to upsert user_subjects")
	}

	if !cached {
		userSubjectUpsertCacheMut.Lock()
		userSubjectUpsertCache[key] = cache
		userSubjectUpsertCacheMut.Unlock()
	}

	return nil
}

// Delete deletes a single UserSubject record with an executor.
// Delete will match against the primary key column to find the record to delete.
func (o *UserSubject) Delete(ctx context.Context, exec boil.ContextExecutor) (int64, error) {
	if o == nil {
		return 0, errors.New("models: no UserSubject provided for delete")
	}

	args := queries.ValuesFromMapping(reflect.Indirect(reflect.ValueOf(o)), userSubjectPrimaryKeyMapping)
	sql := "DELETE FROM \"user_subjects\" WHERE \"id\"=$1"

	if boil.IsDebug(ctx) {
		writer := boil.DebugWriterFrom(ctx)
		fmt.Fprintln(writer, sql)
		fmt.Fprintln(writer, args...)
	}
	result, err := exec.ExecContext(ctx, sql, args...)
	if err != nil {
		return 0, errors.Wrap(err, "models: unable to delete from user_subjects")
	}

	rowsAff, err := result.RowsAffected()
	if err != nil {
		return 0, errors.Wrap(err, "models: failed to get rows affected by delete for user_subjects")
	}

	return rowsAff, nil
}

// DeleteAll deletes all matching rows.
func (q userSubjectQuery) DeleteAll(ctx context.Context, exec boil.ContextExecutor) (int64, error) {
	if q.Query == nil {
		return 0, errors.New("models: no userSubjectQuery provided for delete all")
	}

	queries.SetDelete(q.Query)

	result, err := q.Query.ExecContext(ctx, exec)
	if err != nil {
		return 0, errors.Wrap(err, "models: unable to delete all from user_subjects")
	}

	rowsAff, err := result.RowsAffected()
	if err != nil {
		return 0, errors.Wrap(err, "models: failed to get rows affected by deleteall for user_subjects")
	}

	return rowsAff, nil
}

// DeleteAll deletes all rows in the slice, using an executor.
func (o UserSubjectSlice) DeleteAll(ctx context.Context, exec boil.ContextExecutor) (int64, error) {
	if len(o) == 0 {
		return 0, nil
	}

	var args []interface{}
	for _, obj := range o {
		pkeyArgs := queries.ValuesFromMapping(reflect.Indirect(reflect.ValueOf(obj)), userSubjectPrimaryKeyMapping)
		args = append(args, pkeyArgs...)
	}

	sql := "DELETE FROM \"user_subjects\" WHERE " +
		strmangle.WhereClauseRepeated(string(dialect.LQ), string(dialect.RQ), 1, userSubjectPrimaryKeyColumns, len(o))

	if boil.IsDebug(ctx) {
		writer := boil.DebugWriterFrom(ctx)
		fmt.Fprintln(writer, sql)
		fmt.Fprintln(writer, args)
	}
	result, err := exec.ExecContext(ctx, sql, args...)
	if err != nil {
		return 0, errors.Wrap(err, "models: unable to delete all from userSubject slice")
	}

	rowsAff, err := result.RowsAffected()
	if err != nil {
		return 0, errors.Wrap(err, "models: failed to get rows affected by deleteall for user_subjects")
	}

	return rowsAff, nil
}

// Reload refetches the object from the database
// using the primary keys with an executor.
func (o *UserSubject) Reload(ctx context.Context, exec boil.ContextExecutor) error {
	ret, err := FindUserSubject(ctx, exec, o.ID)
	if err != nil {
		return err
	}

	*o = *ret
	return nil
}

// ReloadAll refetches every row with matching primary key column values
// and overwrites the original object slice with the newly updated slice.
func (o *UserSubjectSlice) ReloadAll(ctx context.Context, exec boil.ContextExecutor) error {
	if o == nil || len(*o) == 0 {
		return nil
	}

	slice := UserSubjectSlice{}
	var args []interface{}
	for _, obj := range *o {
		pkeyArgs := queries.ValuesFromMapping(reflect.Indirect(reflect.ValueOf(obj)), userSubjectPrimaryKeyMapping)
		args = append(args, pkeyArgs...)
	}

	sql := "SELECT \"user_subjects\".* FROM \"user_subjects\" WHERE " +
		strmangle.WhereClauseRepeated(string(dialect.LQ), string(dialect.RQ), 1, userSubjectPrimaryKeyColumns, len(*o))

	q := queries.Raw(sql, args...)

	err := q.Bind(ctx, exec, &slice)
	if err != nil {
		return errors.Wrap(err, "models: unable to reload all in UserSubjectSlice")
	}

	*o = slice

	return nil
}

// UserSubjectExists checks if the UserSubject row exists.
func UserSubjectExists(ctx context.Context, exec boil.ContextExecutor, iD int) (bool, error) {
	var exists bool
	sql := "select exists(select 1 from \"user_subjects\" where \"id\"=$1 limit 1)"

	if boil.IsDebug(ctx) {
		writer := boil.DebugWriterFrom(ctx)
		fmt.Fprintln(writer, sql)
		fmt.Fprintln(writer, iD)
	}
	row := exec.QueryRowContext(ctx, sql, iD)

	err := row.Scan(&exists)
	if err != nil {
		return false, errors.Wrap(err, "models: unable to check if user_subjects exists")
	}

	return exists, nil
}
