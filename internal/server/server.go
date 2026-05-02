package server

import (
	"fmt"
	"net/http"

	"github.com/Lzrb0x/smartBookingGoApi/internal/config"
	"github.com/Lzrb0x/smartBookingGoApi/internal/database"
	"github.com/Lzrb0x/smartBookingGoApi/internal/handlers"
	"github.com/Lzrb0x/smartBookingGoApi/internal/repositories"
	"github.com/Lzrb0x/smartBookingGoApi/internal/services"
	"github.com/Lzrb0x/smartBookingGoApi/internal/swagger"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
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

	r.GET("/swagger.json", func(c *gin.Context) {
		c.Data(http.StatusOK, "application/json; charset=utf-8", swagger.Spec())
	})
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler, ginSwagger.URL("/swagger.json")))

	// Repositories
	userRepo := repositories.NewUserRepository(db)
	barbershopRepo := repositories.NewBarbershopRepository(db)
	ownerRepo := repositories.NewOwnerRepository(db)
	employeeRepo := repositories.NewEmployeeRepository(db)
	serviceRepo := repositories.NewServiceRepository(db)
	barbershopServiceRepo := repositories.NewBarbershopServiceRepository(db)
	serviceEmployeeRepo := repositories.NewServiceEmployeeRepository(db)
	employeeWorkingHourRepo := repositories.NewEmployeeWorkingHourRepository(db)
	employeeWorkingHourOverrideRepo := repositories.NewEmployeeWorkingHourOverrideRepository(db)
	bookingRepo := repositories.NewBookingRepository(db)
	refreshTokenRepo := repositories.NewRefreshTokenRepository(db)

	// Services
	authService := services.NewAuthService(
		userRepo,
		refreshTokenRepo,
		cfg.JWTSecret,
		cfg.JWTAccessTTL,
		cfg.JWTRefreshTTL,
	)

	// Handlers
	userHandler := handlers.NewUserHandler(userRepo)
	barbershopHandler := handlers.NewBarbershopHandler(barbershopRepo)
	ownerHandler := handlers.NewOwnerHandler(ownerRepo)
	employeeHandler := handlers.NewEmployeeHandler(employeeRepo, userRepo)
	serviceHandler := handlers.NewServiceHandler(serviceRepo)
	barbershopServiceHandler := handlers.NewBarbershopServiceHandler(barbershopServiceRepo)
	serviceEmployeeHandler := handlers.NewServiceEmployeeHandler(serviceEmployeeRepo)
	employeeWorkingHourHandler := handlers.NewEmployeeWorkingHourHandler(employeeWorkingHourRepo, employeeRepo)
	employeeWorkingHourOverrideHandler := handlers.NewEmployeeWorkingHourOverrideHandler(employeeWorkingHourOverrideRepo, employeeRepo)
	authHandler := handlers.NewAuthHandler(authService)
	bookingHandler := handlers.NewBookingHandler(
		bookingRepo,
		employeeRepo,
		barbershopServiceRepo,
		serviceEmployeeRepo,
		employeeWorkingHourRepo,
		employeeWorkingHourOverrideRepo,
	)

	// Routes
	v1 := r.Group("/api/v1")
	{
		auth := v1.Group("/auth")
		{
			auth.POST("/register", authHandler.Register)
			auth.POST("/login", authHandler.Login)
			auth.POST("/refresh", authHandler.Refresh)
			auth.POST("/logout", authHandler.Logout)
		}

		users := v1.Group("/users")
		{
			users.GET("", userHandler.GetAll)
			users.GET("/:id", userHandler.GetByID)
			users.GET("/:id/dashboard", bookingHandler.GetUserDashboard)
			users.GET("/:id/bookings/recent", bookingHandler.GetRecentByCustomer)
			users.GET("/:id/bookings/current", bookingHandler.GetCurrentByCustomer)
			users.GET("/:id/barbershops/recent", bookingHandler.GetRecentBarbershopsByCustomer)
			users.POST("", userHandler.Create)
			users.PUT("/:id", userHandler.Update)
			users.DELETE("/:id", userHandler.Delete)
		}

		owners := v1.Group("/owners")
		{
			owners.POST("", ownerHandler.Create)
		}

		barbershops := v1.Group("/barbershops")
		{
			barbershops.GET("", barbershopHandler.GetAll)
			barbershops.GET("/:id", barbershopHandler.GetByID)
			barbershops.POST("", barbershopHandler.Create)
			barbershops.PUT("/:id", barbershopHandler.Update)
			barbershops.DELETE("/:id", barbershopHandler.Delete)

			employees := barbershops.Group("/:id/employees")
			{
				employees.GET("", employeeHandler.GetAll)
				employees.POST("", employeeHandler.Create)
				employees.DELETE("/:employeeId", employeeHandler.Delete)
				employees.GET("/:employeeId/availability", bookingHandler.GetAvailability)

				workingHours := employees.Group("/:employeeId/working-hours")
				{
					workingHours.GET("", employeeWorkingHourHandler.GetAll)
					workingHours.GET("/:workingHourId", employeeWorkingHourHandler.GetByID)
					workingHours.POST("", employeeWorkingHourHandler.Create)
					workingHours.PUT("/:workingHourId", employeeWorkingHourHandler.Update)
					workingHours.DELETE("/:workingHourId", employeeWorkingHourHandler.Delete)
				}

				workingHourOverrides := employees.Group("/:employeeId/working-hour-overrides")
				{
					workingHourOverrides.GET("", employeeWorkingHourOverrideHandler.GetAll)
					workingHourOverrides.GET("/:overrideId", employeeWorkingHourOverrideHandler.GetByID)
					workingHourOverrides.POST("", employeeWorkingHourOverrideHandler.Create)
					workingHourOverrides.PUT("/:overrideId", employeeWorkingHourOverrideHandler.Update)
					workingHourOverrides.DELETE("/:overrideId", employeeWorkingHourOverrideHandler.Delete)
				}
			}

			services := barbershops.Group("/:id/services")
			{
				services.GET("", barbershopServiceHandler.GetAll)
				services.GET("/:service_id", barbershopServiceHandler.GetByID)
				services.POST("", barbershopServiceHandler.Create)
				services.PUT("/:service_id", barbershopServiceHandler.Update)
				services.DELETE("/:service_id", barbershopServiceHandler.Delete)
			}

			serviceEmployees := barbershops.Group("/:id/service-employees")
			{
				serviceEmployees.POST("/assign", serviceEmployeeHandler.AssignService)
				serviceEmployees.DELETE("/unassign", serviceEmployeeHandler.UnassignService)
				serviceEmployees.GET("/employees/:serviceId", serviceEmployeeHandler.GetEmployeesByService)
				serviceEmployees.GET("/services/:employeeId", serviceEmployeeHandler.GetServicesByEmployee)
				serviceEmployees.GET("/check/:employeeId/:serviceId", serviceEmployeeHandler.IsAssigned)
			}

			bookings := barbershops.Group("/:id/bookings")
			{
				bookings.GET("", bookingHandler.GetAllByBarbershop)
				bookings.GET("/employees/:employeeId", bookingHandler.GetAllByEmployee)
				bookings.POST("", bookingHandler.Create)
				bookings.DELETE("/:bookingId", authMiddleware(cfg.JWTSecret), bookingHandler.Cancel)
			}
		}

		services := v1.Group("/services")
		{
			services.GET("", serviceHandler.GetAll)
			services.GET("/:id", serviceHandler.GetByID)
			services.POST("", serviceHandler.Create)
			services.PUT("/:id", serviceHandler.Update)
			services.DELETE("/:id", serviceHandler.Delete)
		}
	}

	return r
}
