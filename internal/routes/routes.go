package routes

import (
	"cinema_backend_system/internal/handlers"
	"cinema_backend_system/internal/middleware"
	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

func SetupUserRoutes(e *echo.Echo, userHandler *handlers.UserHandler, db *gorm.DB) {
	protected := e.Group("/api/v1/user/")
	protected.Use(middleware.AuthMiddleware(db))
	{
		protected.GET("profile", userHandler.Profile)
		protected.PATCH("update", userHandler.Profile)
	}
}

func SetupAuthRoutes(e *echo.Echo, authHandler *handlers.AuthHandler, db *gorm.DB) {
	e.POST("api/v1/auth/login", authHandler.Login)
	e.POST("api/v1/auth/register", authHandler.Register)
	authMiddlewareGroup := e.Group("api/v1/auth/")
	authMiddlewareGroup.Use(middleware.AuthMiddleware(db))
	{
		authMiddlewareGroup.GET("logout", authHandler.Logout)
		authMiddlewareGroup.POST("reset/password", authHandler.ResetPassword)
	}
}

func SetupAdminMovieRoutes(e *echo.Echo, adminMovieHandler *handlers.AdminMovieHandler, db *gorm.DB) {
	movieGroup := e.Group("/api/v1/movie/")
	movieGroup.Use(middleware.AuthMiddleware(db), middleware.AdminMiddleware())
	{
		movieGroup.POST("create", adminMovieHandler.Create)
		movieGroup.PATCH("update", adminMovieHandler.Update)
		movieGroup.DELETE("delete", adminMovieHandler.Delete)
		movieGroup.GET("", adminMovieHandler.Index)
		movieGroup.GET("coming-soon", adminMovieHandler.IndexComingSoon)
		movieGroup.GET("show", adminMovieHandler.Show)
	}
}

func SetupAdminPremiereRoutes(e *echo.Echo, adminPremiereHandler *handlers.AdminPremiereHandler, db *gorm.DB) {
	premiereGroup := e.Group("/api/v1/premiere/")
	premiereGroup.Use(middleware.AuthMiddleware(db), middleware.AdminMiddleware())
	{
		premiereGroup.POST("create", adminPremiereHandler.Create)
		premiereGroup.PATCH("update", adminPremiereHandler.Update)
		premiereGroup.DELETE("delete", adminPremiereHandler.Delete)
		premiereGroup.GET("", adminPremiereHandler.Index)
		premiereGroup.GET("show", adminPremiereHandler.Show)
	}
}

func SetupUserPaymentMethodRoutes(e *echo.Echo, userPaymentMethodHandler *handlers.UserPaymentMethodHandler, db *gorm.DB) {
	premiereGroup := e.Group("/api/v1/payment_method/")
	premiereGroup.Use(middleware.AuthMiddleware(db))
	{
		premiereGroup.POST("create", userPaymentMethodHandler.Create)
		premiereGroup.PATCH("update", userPaymentMethodHandler.Update)
		premiereGroup.DELETE("delete", userPaymentMethodHandler.Delete)
		premiereGroup.GET("", userPaymentMethodHandler.Index)
		premiereGroup.GET("show", userPaymentMethodHandler.Show)
	}
}

func SetupUserPremiereRoutes(e *echo.Echo, userPremiereHandler *handlers.UserPremiereHandler, db *gorm.DB) {
	premiereGroup := e.Group("/api/v1/user/premiere/")
	premiereGroup.Use(middleware.AuthMiddleware(db))
	{
		premiereGroup.GET("seat-map", userPremiereHandler.SeatMap)
	}
}

func SetupUserOrderRoutes(e *echo.Echo, userOrderHandler *handlers.UserOrderHandler, db *gorm.DB) {
	premiereGroup := e.Group("/api/v1/order/")
	premiereGroup.Use(middleware.AuthMiddleware(db))
	{
		premiereGroup.POST("create", userOrderHandler.Create)
		premiereGroup.POST("paid", userOrderHandler.Paid)
		premiereGroup.POST("refund", userOrderHandler.Refund)
		premiereGroup.PATCH("update", userOrderHandler.Update)
		premiereGroup.DELETE("delete", userOrderHandler.Delete)
		premiereGroup.GET("", userOrderHandler.Index)
		premiereGroup.GET("summary", userOrderHandler.Summary)
		premiereGroup.GET("show", userOrderHandler.Show)
	}
}

func SetupUserReviewRoutes(e *echo.Echo, userReviewHandler *handlers.UserReviewHandler, db *gorm.DB) {
	premiereGroup := e.Group("/api/v1/review/")
	premiereGroup.Use(middleware.AuthMiddleware(db))
	{
		premiereGroup.POST("create", userReviewHandler.Create)
		premiereGroup.PATCH("update", userReviewHandler.Update)
		//premiereGroup.DELETE("delete", userOrderHandler.Delete)
		//premiereGroup.GET("", userOrderHandler.Index)
		//premiereGroup.GET("show", userOrderHandler.Show)
	}
}

func SetupAdminReviewRoutes(e *echo.Echo, AdminReviewHandler *handlers.AdminReviewHandler, db *gorm.DB) {
	premiereGroup := e.Group("/api/v1/admin/review/")
	premiereGroup.Use(middleware.AuthMiddleware(db), middleware.AdminMiddleware())
	{
		premiereGroup.POST("approve", AdminReviewHandler.Approve)
	}
}
