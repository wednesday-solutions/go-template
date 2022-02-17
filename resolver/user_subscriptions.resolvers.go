package resolver

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"fmt"

	"github.com/wednesday-solutions/go-template/graphql_models"
	"github.com/wednesday-solutions/go-template/pkg/utl"
)

func (r *subscriptionResolver) UserNotification(ctx context.Context) (<-chan *graphql_models.User, error) {
	id := utl.RandomSequence(5)
	event := make(chan *graphql_models.User, 1)
	go func() {
		<-ctx.Done()
		r.Lock()
		delete(r.Observers, id)
		r.Unlock()
	}()
	r.Lock()
	r.Observers[id] = event
	r.Unlock()
	fmt.Print("Subscribed to user creation updates!")
	return event, nil
}

// Subscription returns graphql_models.SubscriptionResolver implementation.
func (r *Resolver) Subscription() graphql_models.SubscriptionResolver {
	return &subscriptionResolver{r}
}

type subscriptionResolver struct{ *Resolver }
