package resolver

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"fmt"
	"go-template/gqlmodels"
	"go-template/pkg/utl"
)

// UserNotification is the resolver for the userNotification field.
func (r *subscriptionResolver) UserNotification(ctx context.Context) (<-chan *gqlmodels.User, error) {
	id := utl.RandomSequence(5)
	event := make(chan *gqlmodels.User, 1)
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

// Subscription returns gqlmodels.SubscriptionResolver implementation.
func (r *Resolver) Subscription() gqlmodels.SubscriptionResolver { return &subscriptionResolver{r} }

type subscriptionResolver struct{ *Resolver }
