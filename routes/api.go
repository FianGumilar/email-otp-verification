package routes

import (
	"net/http"

	"github.com/FianGumilar/email-otp-verification/config"
	"github.com/FianGumilar/email-otp-verification/services"
	"github.com/FianGumilar/email-otp-verification/services/userService"

	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
)

// RoutesApi
func RoutesApi(e echo.Echo, usecaseSvc services.UsecaseService) {

	public := e.Group("/public")

	userSvc := userService.NewUserService(usecaseSvc)
	userGroup := public.Group("/user")
	userGroup.POST("/register", userSvc.Register)
	userGroup.POST("/validate-otp", userSvc.ValidateOtp)

	private := e.Group("/private")
	private.Use(middleware.JWT([]byte(config.GetEnv("JWT_KEY"))))
	private.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowCredentials: true,
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{http.MethodGet, http.MethodPut, http.MethodPost, http.MethodDelete},
	}))

}
