package graph

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"gin-gorm-gql/graph/generated"
	"gin-gorm-gql/graph/model"
)

func (r *mutationResolver) CreateOrder(ctx context.Context, input model.OrderInput) (*model.Order, error) {
	order := model.Order{
		CustomerName: input.CustomerName,
		OrderAmount:  input.OrderAmount,
		Items:        mapItemsFromInput(input.Items),
	}
	r.DB.Create(&order)

	return &order, nil
}

func (r *mutationResolver) UpdateOrder(ctx context.Context, orderID int, input model.OrderInput) (*model.Order, error) {
	updatedOrder := model.Order{
		ID:           orderID,
		CustomerName: input.CustomerName,
		OrderAmount:  input.OrderAmount,
		Items:        mapItemsFromInput(input.Items),
	}
	r.DB.Save(&updatedOrder)

	return &updatedOrder, nil
}

func (r *mutationResolver) DeleteOrder(ctx context.Context, orderID int) (bool, error) {
	r.DB.Where("id = ?", orderID).Delete(&model.Item{})
	r.DB.Where("id = ?", orderID).Delete(&model.Order{})
	return true, nil
}

func (r *queryResolver) Orders(ctx context.Context) ([]*model.Order, error) {
	var orders []*model.Order
	r.DB.Preload("Items").Find(&orders)

	return orders, nil
}

// Mutation returns generated.MutationResolver implementation.
func (r *Resolver) Mutation() generated.MutationResolver { return &mutationResolver{r} }

// Query returns generated.QueryResolver implementation.
func (r *Resolver) Query() generated.QueryResolver { return &queryResolver{r} }

type mutationResolver struct{ *Resolver }
type queryResolver struct{ *Resolver }
