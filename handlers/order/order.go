package order

import (
	"cart-order-service/helper"
	model "cart-order-service/repository/models"
	"encoding/json"
	"net/http"

	"github.com/go-playground/validator"
	"github.com/google/uuid"
)

type orderDto interface {
	CreateOrder(bReq model.Order) (*uuid.UUID, error)
	UpdatePayment(bReq model.UpdateRequest) (*string, error)
}

type Handler struct {
	order     orderDto
	validator *validator.Validate
}

func NewHandler(order orderDto, validator *validator.Validate) *Handler {
	return &Handler{order, validator}
}

func (h *Handler) CreateOrder(w http.ResponseWriter, r *http.Request) {
	var bReq model.Order
	if err := json.NewDecoder(r.Body).Decode(&bReq); err != nil {
		helper.HandleResponse(w, http.StatusBadRequest, err.Error())
		return
	}

	bReq.RefCode = helper.GenerateRefCode()

	if bReq.ProductOrder == nil {
		bReq.ProductOrder = json.RawMessage("[]")
	}

	if err := h.validator.Struct(bReq); err != nil {
		helper.HandleResponse(w, http.StatusBadRequest, err.Error())
		return
	}

	bRes, err := h.order.CreateOrder(bReq)
	if err != nil {
		helper.HandleResponse(w, http.StatusInternalServerError, err.Error())
		return
	}

	helper.HandleResponse(w, http.StatusCreated, bRes)
}

func (h *Handler) UpdateOrder(w http.ResponseWriter, r *http.Request) {
	var bReq model.UpdateRequest
	if err := json.NewDecoder(r.Body).Decode(&bReq); err != nil {
		helper.HandleResponse(w, http.StatusBadRequest, err.Error())
		return
	}

	if err := h.validator.Struct(&bReq); err != nil {
		helper.HandleResponse(w, http.StatusBadRequest, err.Error())
		return
	}

	// payment success
	message, err := h.order.UpdatePayment(bReq)
	if err != nil {
		helper.HandleResponse(w, http.StatusInternalServerError, err.Error())
		return
	}

	helper.HandleResponse(w, http.StatusOK, message)
}
