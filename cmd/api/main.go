package main

import (
	"context"
	"net/http"

	"ecommerce-app/configs"
	"ecommerce-app/internal/app/api/router"
	"ecommerce-app/internal/infra/db"
	"ecommerce-app/internal/pkg/logger"
)

func main() {
	cfg := configs.Load()
	
	ctx := context.Background()

	pool,err := db.NewPostgresPool(ctx, cfg)
	if err != nil {
		logger.Fatal("Failed to connect to database: %v", err)
	}
	
	defer pool.Close()

	r:= router.NewRouter(pool)

	logger.Info("Server started on port %s", cfg.ServerPort)
	
	http.ListenAndServe(cfg.ServerPort, r)

}
