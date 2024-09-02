package main

import (
	"cart-order-service/config"
	cartHandler "cart-order-service/handlers/cart"
	"cart-order-service/repository/cart"
	"cart-order-service/routes"
	cartUsecase "cart-order-service/usecase/cart"
	"database/sql"
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

	routes := setupRoutes(sqlDb)
	routes.Run(cfg.AppPort)
}

func setupRoutes(db *sql.DB) *routes.Routes {
	cartRepository := cart.NewStore(db)
	cartUseCase := cartUsecase.NewCart(cartRepository)
	cartHandler := cartHandler.NewHandler(cartUseCase)

	return &routes.Routes{
		Cart: cartHandler,
	}
}
