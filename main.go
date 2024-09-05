package main

import (
	"cart-order-service/config"
	cartHandler "cart-order-service/handlers/cart"
	"cart-order-service/repository/cart"
	"cart-order-service/repository/order"
	"cart-order-service/routes"
	cartUsecase "cart-order-service/usecase/cart"
	"database/sql"

	orderHandler "cart-order-service/handlers/order"
	orderUseCase "cart-order-service/usecase/order"

	"github.com/go-playground/validator"
)

func main() {
	cfg, err := config.LoadConfig()
	if err != nil {
		return
	}

	sqlDb, err := config.ConnectToDatabase(config.Connection{
		Host:     cfg.DBHost,
		Port:     cfg.DBPort,
		User:     cfg.DBUser,
		Password: cfg.DBPassword,
		DBName:   cfg.DBName,
	})
	if err != nil {
		return
	}
	defer sqlDb.Close()

	validator := validator.New()

	routes := setupRoutes(sqlDb, validator)
	routes.Run(cfg.AppPort)
}

func setupRoutes(db *sql.DB, validator *validator.Validate) *routes.Routes {
	cartRepository := cart.NewStore(db)
	cartUseCase := cartUsecase.NewCart(cartRepository)
	cartHandler := cartHandler.NewHandler(cartUseCase)

	orderRepository := order.NewStore(db)
	orderUseCase := orderUseCase.NewOrder(orderRepository)
	orderHandler := orderHandler.NewHandler(orderUseCase, validator)

	return &routes.Routes{
		Cart:  cartHandler,
		Order: orderHandler,
	}
}
