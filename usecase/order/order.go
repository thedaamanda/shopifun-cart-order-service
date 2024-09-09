package order

import (
	model "cart-order-service/repository/models"

	"github.com/google/uuid"
)

type orderStore interface {
	CreateOrder(bReq model.Order) (*uuid.UUID, *string, error)
	CreateOrderItemsLogs(bReq model.OrderItemsLogs) (*string, error)
	UpdateOrder(bReq model.UpdateRequest) (*string, error)
}

type order struct {
	store orderStore
}

func NewOrder(store orderStore) *order {
	return &order{store}
}

func (o *order) CreateOrder(bReq model.Order) (*uuid.UUID, error) {
	orderID, refCode, err := o.store.CreateOrder(bReq)
	if err != nil {
		return nil, err
	}

	_, err = o.store.CreateOrderItemsLogs(model.OrderItemsLogs{
		OrderID:    *orderID,
		RefCode:    *refCode,
		FromStatus: "",
		ToStatus:   model.OrderStatusPending,
		Notes:      "Order created",
	})
	if err != nil {
		return nil, err
	}

	return orderID, nil
}

func (o *order) UpdatePayment(bReq model.UpdateRequest) (*string, error) {
	refCode, err := o.store.UpdateOrder(bReq)
	if err != nil {
		return nil, err
	}

	_, err = o.store.CreateOrderItemsLogs(model.OrderItemsLogs{
		OrderID:    bReq.OrderID,
		RefCode:    *refCode,
		FromStatus: model.OrderStatusPending,
		ToStatus:   model.OrderStatusPaid,
		Notes:      "Payment success",
	})
	if err != nil {
		return nil, err
	}

	updateOK := "Payment Success"
	return &updateOK, nil
}
