package daos

import (
	"context"
	"fmt"

	"github.com/volatiletech/sqlboiler/v4/queries/qm"
	"github.com/wednesday-solutions/go-template/graphql_models"
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

// FindPaginationStatus returns if next and previoius pages exist
func GetSubjectsPageInfo(startID int, endID int) (*graphql_models.PageInfo, error) {
	contextExecutor := getContextExecutor(nil)
	startCursor := EncodeCursor(startID)
	endCursor := EncodeCursor(endID)
	pageInfo := &graphql_models.PageInfo{
		StartCursor: startCursor,
		EndCursor:   endCursor,
	}
	findNextPageQuery := []qm.QueryMod{qm.Where(fmt.Sprintf("ID > %v", endID))}
	hasNextPage, err := models.Subjects(findNextPageQuery...).Exists(context.Background(), contextExecutor)
	pageInfo.HasNextPage = hasNextPage
	if err != nil {
		return pageInfo, err
	}
	findPreviousPageQuery := []qm.QueryMod{qm.Where(fmt.Sprintf("ID < %v", startID))}
	hasPreviousPage, err := models.Subjects(findPreviousPageQuery...).Exists(context.Background(), contextExecutor)
	pageInfo.HasPreviousPage = hasPreviousPage
	if err != nil {
		return pageInfo, err
	}
	return pageInfo, nil
}
