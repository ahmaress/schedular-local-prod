package router

import (
	"github.com/gin-gonic/gin"

	"scheduler-service/internal/app"
	"scheduler-service/internal/config"
	"scheduler-service/internal/handlers"
	"scheduler-service/internal/repository/postgres"
	"scheduler-service/internal/service"
)

func Build(appInstance *app.App, cfg *config.Config) *gin.Engine {
	r := gin.Default()

	// OAuth2 callback (must be before auth middleware)
	r.GET("/oauth2callback", appInstance.GoogleOAuth2CallbackHandler)

	// Reuse existing auth middleware logic
	r.Use(app.AuthMiddlewareFromEnv())

	api := r.Group("/api")
	{
		availRepo := postgres.NewAvailabilityRepo()
		bookingRepo := postgres.NewBookingRepo()
		availService := service.NewAvailabilityService(appInstance.DB, availRepo, bookingRepo)
		bookingService := service.NewBookingService(appInstance.DB, bookingRepo, availService)

		availHandlers := &handlers.AvailabilityHandlers{DB: appInstance.DB, AvailSv: availService, BookSv: bookingService}

		users := api.Group("/users")
		{
			users.POST("/:id/availability", availHandlers.SetAvailability)
			users.PUT("/:id/availability/:rule_id", availHandlers.UpdateAvailability)
			users.GET("/:id/availability", availHandlers.ListAvailability)
			users.GET("/:id/slots", availHandlers.GetSlots)
			users.POST("/:id/bookings", availHandlers.CreateBooking)
			users.GET("/:id/bookings", availHandlers.ListBookings)
		}

		api.DELETE("/bookings/:id", availHandlers.CancelBooking)

		// Google Calendar integration routes - delegate to existing methods to avoid logic change
		calendar := api.Group("/calendar")
		{
			calendar.GET("/auth", appInstance.GoogleAuthHandler)
			calendar.GET("/events", appInstance.GetGoogleCalendarEvents)
			calendar.GET("/calendars", appInstance.GetGoogleCalendarList)
			calendar.POST("/refresh-token", appInstance.RefreshGoogleToken)
		}
	}

	return r
}
