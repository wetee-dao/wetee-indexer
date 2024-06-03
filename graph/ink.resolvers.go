package graph

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.
// Code generated by github.com/99designs/gqlgen version v0.17.47

import (
	"context"
	"encoding/json"

	"github.com/wetee-dao/indexer/graph/model"
	"github.com/wetee-dao/indexer/store"
)

// UploadContractAbi is the resolver for the upload_contract_abi field.
func (r *mutationResolver) UploadContractAbi(ctx context.Context, project string, abi string) (bool, error) {
	abiMap := make(map[string]interface{})
	json.Unmarshal([]byte(abi), &abiMap)
	source := abiMap["source"].(map[string]interface{})
	hash := source["hash"].(string)

	cacheAbi, err := store.GetInkCodeBlock(hash)
	if err != nil {
		return false, err
	}
	if cacheAbi == "" {
		store.SetInkCodeBlock(hash, abi)
		store.AddToList("ink_abi", project, []byte(hash))
	}

	return true, nil
}

// ListContractAbi is the resolver for the list_contract_abi field.
func (r *queryResolver) ListContractAbi(ctx context.Context, project string, page int, pageSize int) ([]string, error) {
	list, err := store.GetList("ink_abi", project, page, pageSize)
	if err != nil {
		return nil, err
	}

	slist := make([]string, 0, len(list))
	for _, v := range list {
		slist = append(slist, string(v))
	}

	return slist, nil
}

// ListContract is the resolver for the list_contract field.
func (r *queryResolver) ListContract(ctx context.Context, project string, page int, pageSize int) ([]*model.Contract, error) {
	list, err := store.GetList("contract", project, page, pageSize)
	if err != nil {
		return nil, err
	}

	slist := make([]*model.Contract, 0, len(list))
	codes := map[string]string{}
	for _, v := range list {
		c := model.Contract{}
		json.Unmarshal(v, &c)
		slist = append(slist, &c)
		if cacheAbi, ok := codes[c.CodeHash]; !ok {
			abi, err := store.GetInkCodeBlock("0x" + c.CodeHash)
			if err == nil {
				codes[c.CodeHash] = abi
				c.Abi = abi
			}
		} else {
			c.Abi = cacheAbi
		}
	}

	return slist, nil
}

// Mutation returns MutationResolver implementation.
func (r *Resolver) Mutation() MutationResolver { return &mutationResolver{r} }

type mutationResolver struct{ *Resolver }
