package daos

import (
	"context"
	"fmt"

	"github.com/volatiletech/sqlboiler/v4/queries/qm"
	"github.com/wednesday-solutions/go-template/models"
	resultwrapper "github.com/wednesday-solutions/go-template/pkg/utl/result_wrapper"
)

// FindSubjectById ...
func FindUserSubjectById(userID *int, subjectID *int) (*models.UserSubject, error) {
	var queryMods []qm.QueryMod
	if userID != nil {
		queryMods = append(queryMods, qm.Where("user_id = ?", *userID))
	}
	if subjectID != nil {
		queryMods = append(queryMods, qm.Where("subject_id = ?", *subjectID))
	}
	if len(queryMods) == 0 {
		return nil, fmt.Errorf("either user_id or subject_id must be set")
	}
	contextExecutor := getContextExecutor(nil)
	return models.UserSubjects(queryMods...).One(context.Background(), contextExecutor)
}

// FindAllSubjectsWithCount ...
func FindAllUserSubjectsWithCount(queryMods []qm.QueryMod) (models.UserSubjectSlice, int64, error) {
	contextExecutor := getContextExecutor(nil)
	userSubjects, err := models.UserSubjects(queryMods...).All(context.Background(), contextExecutor)
	if err != nil {
		return models.UserSubjectSlice{}, 0, err
	}
	queryMods = append(queryMods, qm.Offset(0))
	count, err := models.UserSubjects(queryMods...).Count(context.Background(), contextExecutor)
	return userSubjects, count, err
}

func GetUserForUserSubject(us *models.UserSubject) (*models.User, error) {
	var queryMods []qm.QueryMod
	contextExecutor := getContextExecutor(nil)
	user, err := us.User(queryMods...).One(context.Background(), contextExecutor)
	if err != nil {
		return nil, resultwrapper.ResolverSQLError(err, "data")
	}
	return user, nil
}
func GetSubjectForUserSubject(us *models.UserSubject) (*models.Subject, error) {
	var queryMods []qm.QueryMod
	contextExecutor := getContextExecutor(nil)
	subject, err := us.Subject(queryMods...).One(context.Background(), contextExecutor)
	if err != nil {
		return nil, resultwrapper.ResolverSQLError(err, "data")
	}
	return subject, nil
}
