package graph

import "gin-gorm-gql/graph/model"

func mapItemsFromInput(itemsInput []*model.ItemInput) []*model.Item {
	var items []*model.Item
	for _, itemInput := range itemsInput {
		if itemInput.ID == nil {
			items = append(items, &model.Item{
				ProductCode: itemInput.ProductCode,
				ProductName: itemInput.ProductName,
				Quantity:    itemInput.Quantity,
			})
		} else {
			items = append(items, &model.Item{
				ID:          *itemInput.ID,
				ProductCode: itemInput.ProductCode,
				ProductName: itemInput.ProductName,
				Quantity:    itemInput.Quantity,
			})
		}
	}
	return items
}
