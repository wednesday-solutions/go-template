package helper

import (
	"context"
	"github.com/99designs/gqlgen/graphql"
)

func GetInputFromContext(ctx context.Context, key string) map[string]interface{} {
	requestContext := graphql.GetOperationContext(ctx)
	//return requestContext.Variables[key]
	return requestContext.Variables
}
