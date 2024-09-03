package cart

import (
	model "cart-order-service/repository/models"
)

// cartStore is an interface that defines the methods required for managing a shopping cart.
type cartStore interface {
	// GetCartByUserID retrieves the cart for a given user.
	GetCartByUserID(bReq model.GetCartRequest) (*[]model.Cart, error)
}

// cart is a struct that holds the store for managing a shopping cart.
type cart struct {
	store cartStore
}

// NewCart is a constructor function that returns a new cart instance.
func NewCart(store cartStore) *cart {
	return &cart{store}
}

// GetCartByUserID is a method that retrieves the cart for a given user and returns a response with the total items.
func (c *cart) GetCartByUserID(bReq model.GetCartRequest) (*[]model.Cart, error) {
	result, err := c.store.GetCartByUserID(bReq)
	if err != nil {
		return nil, err
	}

	if len(*result) == 0 {
		return nil, nil
	}

	return result, nil
}
