package graph

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.
// Code generated by github.com/99designs/gqlgen version v0.17.47

import (
	"context"

	"github.com/wetee-dao/indexer/store"
)

// UploadContract is the resolver for the upload_contract field.
func (r *mutationResolver) UploadContract(ctx context.Context, project string, abi string) (bool, error) {
	return true, store.AddToList("ink_contract", []byte(project), []byte(abi))
}

// ListContract is the resolver for the list_contract field.
func (r *queryResolver) ListContract(ctx context.Context, project string, page int, pageSize int) ([]string, error) {
	list, err := store.GetList("ink_contract", []byte(project), page, pageSize)
	if err != nil {
		return nil, err
	}

	slist := make([]string, 0, len(list))
	for _, v := range list {
		slist = append(slist, string(v))
	}

	return slist, nil
}

// Mutation returns MutationResolver implementation.
func (r *Resolver) Mutation() MutationResolver { return &mutationResolver{r} }

// Query returns QueryResolver implementation.
func (r *Resolver) Query() QueryResolver { return &queryResolver{r} }

type mutationResolver struct{ *Resolver }
type queryResolver struct{ *Resolver }
