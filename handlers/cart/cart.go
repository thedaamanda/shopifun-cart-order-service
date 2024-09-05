package cart

import (
	"cart-order-service/helper"
	model "cart-order-service/repository/models"
	"encoding/json"
	"net/http"

	"github.com/google/uuid"
)

// cartDto is an interface that defines the methods that our Handler struct depends on.
type cartDto interface {
	GetCartByUserID(bReq model.GetCartRequest) (*[]model.Cart, error)
	AddCart(bReq model.Cart) (*uuid.UUID, error)
	UpdateQty(bReq model.Cart) (string, error)
	DeleteCart(bReq model.DeleteCartRequest) (string, error)
}

// Handler is a struct that holds a cartDto.
type Handler struct {
	cart cartDto
}

// NewHandler is a constructor function that returns a new Handler.
func NewHandler(cart cartDto) *Handler {
	return &Handler{cart}
}

// GetCartByUserID is a handler function to get a cart by user id.
// It first extracts the user id from the URL path, then decodes the request body into a GetCartRequest model.
// It then calls the GetCartByUserID method of the cartDto and sends the helper back to the client.
func (h *Handler) GetCartByUserID(w http.ResponseWriter, r *http.Request) {
	userID := r.PathValue("user_id")
	if userID == "" {
		helper.HandleResponse(w, http.StatusBadRequest, "User ID is required")
		return
	}

	uid, err := uuid.Parse(userID)
	if err != nil {
		helper.HandleResponse(w, http.StatusBadRequest, err.Error())
		return
	}

	var bReq model.GetCartRequest
	if err := json.NewDecoder(r.Body).Decode(&bReq); err != nil {
		helper.HandleResponse(w, http.StatusBadRequest, err.Error())
		return
	}

	var pidSlice []uuid.UUID
	pidSlice = append(pidSlice, bReq.ProductID...)

	bReq = model.GetCartRequest{
		UserID:    uid,
		ProductID: pidSlice,
	}

	bResp, err := h.cart.GetCartByUserID(bReq)
	if err != nil {
		helper.HandleResponse(w, http.StatusInternalServerError, err.Error())
		return
	}

	helper.HandleResponse(w, http.StatusOK, bResp)
}

func (h *Handler) AddCart(w http.ResponseWriter, r *http.Request) {
	var bReq model.Cart
	if err := json.NewDecoder(r.Body).Decode(&bReq); err != nil {
		helper.HandleResponse(w, http.StatusBadRequest, err.Error())
		return
	}

	if bReq.Qty == 0 {
		helper.HandleResponse(w, http.StatusBadRequest, "Qty must be greater than 0")
		return
	}

	bResp, err := h.cart.AddCart(bReq)
	if err != nil {
		helper.HandleResponse(w, http.StatusInternalServerError, err.Error())
		return
	}

	helper.HandleResponse(w, http.StatusOK, bResp)
}

func (h *Handler) UpdateCart(w http.ResponseWriter, r *http.Request) {
	userID := r.PathValue("user_id")

	uid, err := uuid.Parse(userID)
	if err != nil {
		helper.HandleResponse(w, http.StatusBadRequest, err.Error())
		return
	}

	var bReq model.Cart
	if err := json.NewDecoder(r.Body).Decode(&bReq); err != nil {
		helper.HandleResponse(w, http.StatusBadRequest, err.Error())
		return
	}
	bReq.UserID = uid

	bResp, err := h.cart.UpdateQty(bReq)
	if err != nil {
		helper.HandleResponse(w, http.StatusInternalServerError, err.Error())
		return
	}

	helper.HandleResponse(w, http.StatusOK, bResp)
}

func (h *Handler) DeleteCart(w http.ResponseWriter, r *http.Request) {
	userID := r.PathValue("user_id")

	uid, err := uuid.Parse(userID)
	if err != nil {
		helper.HandleResponse(w, http.StatusBadRequest, err.Error())
		return
	}

	var bReq model.DeleteCartRequest
	if err := json.NewDecoder(r.Body).Decode(&bReq); err != nil {
		helper.HandleResponse(w, http.StatusBadRequest, err.Error())
		return
	}
	bReq.UserID = uid

	bResp, err := h.cart.DeleteCart(bReq)
	if err != nil {
		helper.HandleResponse(w, http.StatusInternalServerError, err.Error())
		return
	}

	helper.HandleResponse(w, http.StatusOK, bResp)
}
