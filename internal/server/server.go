package server

import (
	"fmt"
	"net/http"

	"github.com/Lzrb0x/smartBookingGoApi/internal/config"
	"github.com/Lzrb0x/smartBookingGoApi/internal/database"
	"github.com/Lzrb0x/smartBookingGoApi/internal/handlers"
	"github.com/Lzrb0x/smartBookingGoApi/internal/repositories"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func NewServer(cfg *config.Config, db *database.DB) *http.Server {
	return &http.Server{
		Addr:         fmt.Sprintf(":%s", cfg.Port),
		Handler:      registerRouters(cfg, db),
		ReadTimeout:  cfg.ServerReadTimeout,
		WriteTimeout: cfg.ServerWriteTimeout,
	}
}

func registerRouters(cfg *config.Config, db *database.DB) http.Handler {
	r := gin.Default()

	r.Use(cors.New(cors.Config{
		AllowOrigins:     cfg.CORSAllowOrigins,
		AllowMethods:     cfg.CORSAllowMethods,
		AllowHeaders:     cfg.CORSAllowHeaders,
		AllowCredentials: cfg.CORSAllowCredentials,
	}))

	// Repositories
	userRepo := repositories.NewUserRepository(db)
	barbershopRepo := repositories.NewBarbershopRepository(db)

	// Handlers
	userHandler := handlers.NewUserHandler(userRepo)
	barbershopHandler := handlers.NewBarbershopHandler(barbershopRepo)

	// Routes
	v1 := r.Group("/api/v1")
	{
		users := v1.Group("/users")
		{
			users.GET("", userHandler.GetAll)
			users.GET("/:id", userHandler.GetByID)
			users.POST("", userHandler.Create)
			users.PUT("/:id", userHandler.Update)
			users.DELETE("/:id", userHandler.Delete)
		}

		barbershops := v1.Group("/barbershops")
		{
			barbershops.GET("", barbershopHandler.GetAll)
			barbershops.GET("/:id", barbershopHandler.GetByID)
			barbershops.POST("", barbershopHandler.Create)
			barbershops.PUT("/:id", barbershopHandler.Update)
			barbershops.DELETE("/:id", barbershopHandler.Delete)
		}
	}

	return r
}
