package cart

import (
	model "cart-order-service/repository/models"
	"database/sql"
	"errors"
	"fmt"
	"strings"

	"github.com/google/uuid"
)

type store struct {
	db *sql.DB
}

// NewStore is a constructor function that returns a new store instance.
func NewStore(db *sql.DB) *store {
	return &store{db}
}

// GetCartByUserID is a method that retrieves the cart for a given user.
// It returns a slice of cart and an error if any occurs during the retrieval process.
func (s *store) GetCartByUserID(bReq model.GetCartRequest) (*[]model.Cart, error) {
	querySelect := `
		SELECT
			*
		FROM cart_items
	`

	var queryConditions []string

	if bReq.UserID != uuid.Nil {
		queryConditions = append(queryConditions, fmt.Sprintf("user_id = '%s'", bReq.UserID))
	}

	if len(bReq.ProductID) > 0 {
		var productIDs []string
		for _, pid := range bReq.ProductID {
			productIDs = append(productIDs, fmt.Sprintf("'%s'", pid))
		}
		queryConditions = append(queryConditions, fmt.Sprintf("product_id IN (%s)", strings.Join(productIDs, ",")))
	}

	if len(queryConditions) > 0 {
		querySelect += " WHERE " + strings.Join(queryConditions, " AND ")
	} else {
		querySelect += " WHERE deleted_at IS NULL"
	}

	querySelect += " AND deleted_at IS NULL"

	rows, err := s.db.Query(querySelect)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var carts []model.Cart
	for rows.Next() {
		var cart model.Cart
		if err := rows.Scan(
			&cart.ID,
			&cart.UserID,
			&cart.ProductID,
			&cart.Qty,
			&cart.CreatedAt,
			&cart.UpdatedAt,
			&cart.DeletedAt,
		); err != nil {
			return nil, err
		}
		carts = append(carts, cart)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return &carts, nil
}

func (s *store) AddCart(bReq model.Cart) (*uuid.UUID, error) {
	tx, err := s.db.Begin()
	if err != nil {
		return nil, err
	}

	var id uuid.UUID
	queryCreate := `
		INSERT INTO cart_items (
			user_id,
			product_id,
			qty,
			created_at
		) VALUES (
			$1,
			$2,
			$3,
			NOW()
		) RETURNING id
	`
	if err := tx.QueryRow(
		queryCreate,
		bReq.UserID,
		bReq.ProductID,
		bReq.Qty,
	).Scan(&id); err != nil {
		tx.Rollback()
		return nil, err
	}

	if err := tx.Commit(); err != nil {
		tx.Rollback()
		return nil, err
	}

	return &id, nil
}

func (s *store) UpdateQty(userID, productID uuid.UUID, qty int) error {
	tx, err := s.db.Begin()
	if err != nil {
		return err
	}

	queryLock := `
		SELECT 1
		FROM cart_items
		WHERE user_id = $1
		FOR UPDATE
	`
	if _, err := tx.Exec(queryLock, userID); err != nil {
		tx.Rollback()
		return errors.New("failed to lock data")
	}

	queryUpdate := `
		UPDATE cart_items
		SET qty = $1
		WHERE user_id = $2 AND product_id = $3
	`
	if _, err := tx.Exec(queryUpdate, qty, userID, productID); err != nil {
		tx.Rollback()
		return errors.New("failed to update data")
	}

	if err := tx.Commit(); err != nil {
		tx.Rollback()
		return err
	}

	return nil
}

func (s *store) DeleteProduct(bReq model.DeleteCartRequest) error {
	tx, err := s.db.Begin()
	if err != nil {
		return err
	}

	queryLock := `
		SELECT 1
		FROM cart_items
		WHERE user_id = $1
		FOR UPDATE
	`
	if _, err := tx.Exec(queryLock, bReq.UserID); err != nil {
		tx.Rollback()
		return errors.New("failed to lock data")
	}

	queryUpdate := `
		UPDATE cart_items
		SET deleted_at = NOW()
		WHERE user_id = $1 AND product_id = $2
	`
	if _, err := tx.Exec(queryUpdate, bReq.UserID, bReq.ProductID); err != nil {
		tx.Rollback()
		return errors.New("failed to delete data")
	}

	if err := tx.Commit(); err != nil {
		tx.Rollback()
		return err
	}

	return nil
}
