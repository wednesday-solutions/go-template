package daos

import (
	"context"

	"github.com/volatiletech/sqlboiler/queries/qm"
	"github.com/wednesday-solutions/go-template/models"
)

// FindSubjectById ...
func FindSubjectById(subjectID int) (*models.Subject, error) {
	contextExecutor := getContextExecutor(nil)
	return models.FindSubject(context.Background(), contextExecutor, subjectID)
}

// FindAllSubjectsWithCount ...
func FindAllSubjectsWithCount(queryMods []qm.QueryMod) (models.SubjectSlice, int64, error) {
	contextExecutor := getContextExecutor(nil)
	subjects, err := models.Subjects(queryMods...).All(context.Background(), contextExecutor)
	if err != nil {
		return models.SubjectSlice{}, 0, err
	}
	if len(subjects) == 0 {
		return subjects, 0, nil
	}
	queryMods = append(queryMods, qm.Offset(0))
	count, err := models.Subjects(queryMods...).Count(context.Background(), contextExecutor)
	return subjects, count, err
}
