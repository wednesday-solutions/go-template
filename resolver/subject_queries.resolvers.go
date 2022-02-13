package resolver

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"fmt"

	"github.com/labstack/echo"
	"github.com/volatiletech/sqlboiler/v4/queries/qm"
	"github.com/wednesday-solutions/go-template/daos"
	"github.com/wednesday-solutions/go-template/graphql_models"
	"github.com/wednesday-solutions/go-template/pkg/utl/convert"
	resultwrapper "github.com/wednesday-solutions/go-template/pkg/utl/result_wrapper"
)

func (r *queryResolver) Subject(ctx context.Context, id int) (*graphql_models.Subject, error) {
	subject, err := daos.FindSubjectById(id)
	if err != nil {
		return nil, resultwrapper.ResolverSQLError(err, "data")
	}
	return convert.SubjectToGraphQlSubject(subject), nil
}

func (r *queryResolver) Subjects(ctx context.Context, pagination graphql_models.Pagination) (*graphql_models.SubjectPayload, error) {
	var queryMods []qm.QueryMod
	if pagination.Limit != 0 {
		queryMods = append(queryMods, qm.Limit(pagination.Limit), qm.Offset((pagination.Page-1)*pagination.Limit))
	}

	subjects, count, err := daos.FindAllSubjectsWithCount(queryMods)
	if err != nil {
		return nil, resultwrapper.ResolverSQLError(err, "data")
	}
	return &graphql_models.SubjectPayload{
		Subjects: convert.SubjectsToGraphQlSubjects(subjects),
		Total:    int(count),
	}, nil
}

func (r *queryResolver) SubjectConnection(ctx context.Context, forward *graphql_models.ForwardSubjectsInput, backward *graphql_models.BackwardSubjectsInput) (*graphql_models.SubjectConnection, error) {
	queryMods := []qm.QueryMod{}
	if forward != nil {
		if forward.First > 0 {
			var offsetId int
			fmt.Printf("forward.After = %v \n", forward.After)
			if forward.After != nil {
				offset, err := daos.DecodeCursor(*forward.After)
				if err != nil {
					panic(fmt.Sprintf("Invalid Cursor!, got %v", offset))
				}
				offsetId = offset
			}
			fmt.Printf("decoded offsetId, %v \n", offsetId)

			limit := qm.Limit(forward.First)
			offset := qm.Where(fmt.Sprintf("ID > %v", offsetId))
			queryMods = append(queryMods, limit, offset)
		}
	}
	subjects, _, err := daos.FindAllSubjectsWithCount(queryMods)

	if err != nil {
		return nil, resultwrapper.ResolverSQLError(err, "data")
	}

	startID := subjects[0].ID
	endID := subjects[len(subjects)-1].ID

	pageInfo, err := daos.GetSubjectsPageInfo(startID, endID)

	if err != nil {
		return nil, resultwrapper.InternalServerErrorFromMessage(echo.New().AcquireContext(), err.Error())
	}

	edges := convert.ToConnectionResult(subjects)

	return &graphql_models.SubjectConnection{
		Edges:    edges,
		PageInfo: pageInfo,
	}, nil
}
