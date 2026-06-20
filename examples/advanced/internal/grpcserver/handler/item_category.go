package handler

import (
	"context"
	"log"

	"gitlab.com/iktdev-boilerplate/go/protoc-gen-go-mapper/examples/advanced/internal/gen/item_categoriespb"
	"gitlab.com/iktdev-boilerplate/go/protoc-gen-go-mapper/examples/advanced/internal/postgres"
	"gitlab.com/iktdev-boilerplate/go/protoc-gen-go-mapper/examples/advanced/pkg/utils"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type ItemCategoriesHandler struct {
	item_categoriespb.UnimplementedItemCategoryServiceServer
	pg *postgres.Postgres
}

var _ item_categoriespb.ItemCategoryServiceServer = &ItemCategoriesHandler{}

type OptsFucIC func(*ItemCategoriesHandler) (err error)

func WithhDatabaseIC(pg *postgres.Postgres) OptsFucIC {
	return func(ich *ItemCategoriesHandler) (err error) {
		ich.pg = pg
		return
	}
}

func NewItemCategories(
	opts ...OptsFucIC,
) *ItemCategoriesHandler {
	itc := &ItemCategoriesHandler{}
	for _, opt := range opts {
		if err := opt(itc); err != nil {
			log.Fatal("ItemCategoriesHandler fatal connection database")
		}
	}
	return itc
}

func (h *ItemCategoriesHandler) Create(ctx context.Context, req *item_categoriespb.CreateItemCategoryRequest) (res *item_categoriespb.ItemCategoryResponse, err error) {
	return nil, status.Error(codes.Unimplemented, "method Create not implemented")
}

func (h *ItemCategoriesHandler) Update(ctx context.Context, req *item_categoriespb.UpdateItemCategoryRequest) (res *item_categoriespb.ItemCategoryResponse, err error) {
	return nil, status.Error(codes.Unimplemented, "method Update not implemented")
}

func (h *ItemCategoriesHandler) Delete(ctx context.Context, req *item_categoriespb.DeleteItemCategoryRequest) (res *item_categoriespb.ItemCategoryResponse, err error) {
	return nil, status.Error(codes.Unimplemented, "method Delete not implemented")
}

func (h *ItemCategoriesHandler) Restore(ctx context.Context, req *item_categoriespb.RestoreItemCategoryRequest) (res *item_categoriespb.ItemCategoryResponse, err error) {
	return nil, status.Error(codes.Unimplemented, "method Restore not implemented")
}

func (h *ItemCategoriesHandler) Get(ctx context.Context, req *item_categoriespb.GetItemCategoryRequest) (res *item_categoriespb.ItemCategoryResponse, err error) {
	return nil, status.Error(codes.Unimplemented, "method Get not implemented")
}

func (h *ItemCategoriesHandler) List(ctx context.Context, req *item_categoriespb.ListItemCategoryRequest) (res *item_categoriespb.ListItemCategoryResponse, err error) {
	responses, err := h.pg.Q.ListsItemCategories(ctx, utils.StringToPgUUID(req.TenantId))
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	// Use generated response helper
	return item_categoriespb.ToProtoListItemCategoryResponseFromRows(responses, int64(len(responses)), req.Page, req.Limit), nil
}

func (h *ItemCategoriesHandler) Tree(ctx context.Context, req *item_categoriespb.ListItemCategoryTreeRequest) (res *item_categoriespb.ListItemCategoryTreeResponse, err error) {
	rows, err := h.pg.Q.ListsItemCategoriesTree(ctx, utils.StringToPgUUID(req.TenantId))
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	// Use generated response helper
	return item_categoriespb.ToProtoListItemCategoryTreeResponseFromRows(rows), nil
}
