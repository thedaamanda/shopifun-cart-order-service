package order

import (
	model "cart-order-service/repository/models"
	"database/sql"

	"github.com/google/uuid"
)

type store struct {
	db *sql.DB
}

func NewStore(db *sql.DB) *store {
	return &store{db}
}

// CreateOrder is a method that creates a new order and returns the order ID.
// It returns an error if any occurs during the creation process.
func (o *store) CreateOrder(bReq model.Order) (*uuid.UUID, *string, error) {
	tx, err := o.db.Begin()
	if err != nil {
		return nil, nil, err
	}

	queryCreate := `
		INSERT INTO orders (
			user_id,
			payment_type_id,
			order_number,
			total_price,
			product_order,
			status,
			is_paid,
			ref_code,
			created_at
		) VALUES (
			$1, $2, $3, $4, $5, $6, $7, $8, NOW()
		) RETURNING id, ref_code
	`

	var orderID uuid.UUID
	var refCode string
	if err := tx.QueryRow(
		queryCreate,
		bReq.UserID,
		bReq.PaymentTypeID,
		bReq.OrderNumber,
		bReq.TotalPrice,
		bReq.ProductOrder,
		bReq.Status,
		bReq.IsPaid,
		bReq.RefCode,
	).Scan(&orderID, &refCode); err != nil {
		tx.Rollback()
		return nil, nil, err
	}

	if err := tx.Commit(); err != nil {
		tx.Rollback()
		return nil, nil, err
	}

	return &orderID, &refCode, nil
}

// createOrderItemsLogs is a method that creates a new order items log.
// It returns an error if any occurs during the creation process.
func (o *store) CreateOrderItemsLogs(bReq model.OrderItemsLogs) (*string, error) {
	tx, err := o.db.Begin()
	if err != nil {
		return nil, err
	}

	queryCreate := `
		INSERT INTO order_status_logs (
			order_id,
			ref_code,
			from_status,
			to_status,
			notes,
			created_at
		) VALUES (
			$1, $2, $3, $4, $5, NOW()
		) RETURNING ref_code
	`

	var refCode string
	if err := tx.QueryRow(
		queryCreate,
		bReq.OrderID,
		bReq.RefCode,
		bReq.FromStatus,
		bReq.ToStatus,
		bReq.Notes,
	).Scan(&refCode); err != nil {
		tx.Rollback()
		return nil, err
	}

	if err := tx.Commit(); err != nil {
		tx.Rollback()
		return nil, err
	}

	return &refCode, nil
}

func (o *store) UpdateOrder(bReq model.UpdateRequest) (*string, error) {
	tx, err := o.db.Begin()
	if err != nil {
		return nil, err
	}

	queryUpdate := `
		UPDATE orders SET
			status = $1,
			is_paid = $2, 
			updated_at = NOW()
		WHERE id = $3 RETURNING ref_code
	`

	var refCode string
	if err := tx.QueryRow(
		queryUpdate,
		bReq.Status,
		bReq.IsPaid,
		bReq.OrderID,
	).Scan(&refCode); err != nil {
		tx.Rollback()
		return nil, err
	}

	if err := tx.Commit(); err != nil {
		tx.Rollback()
		return nil, err
	}

	return &refCode, nil
}
