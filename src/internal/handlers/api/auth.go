package handler_api

import (
	_ "libs/src/docs"
	"libs/src/internal/dto"
	services "libs/src/internal/usecase"
	usecase_errors "libs/src/internal/usecase/errors"
	"libs/src/settings"
	"net/http"

	"github.com/gin-gonic/gin"
)

// @Summary User register
// @Description Register a new user behind the selected fields
// @Tags Auth
// @Accept json
// @Produce json
// @Param user body dto.RegisterRequest true "Data to register"
// @Success 200 {object} dto.RegisterResponse
// @Failure 400 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /accounts/auth/register [post]
func Register(c *gin.Context) {
	app := c.MustGet("app").(*settings.App)
	user := c.MustGet("user").(dto.UserDTO)

	var registerData dto.RegisterRequest

	if err := c.ShouldBindJSON(&registerData); err != nil {
		c.Error(usecase_errors.BadRequestError{Msg: err.Error()})
		return
	}

	service := services.NewAuthService(app)

	err := service.RegisterUser(c.Request.Context(), user, registerData)
	if err != nil {
		c.Error(err)
		return
	}
	c.JSON(http.StatusOK, dto.RegisterResponse{Message: "success", Status: true})
}

// @Summary User confirm registration
// @Description Confirm users email
// @Tags Auth
// @Accept json
// @Produce json
// @Param token path string true "Token to confirm account"
// @Success 200 {object} dto.MessageResponse
// @Failure 400 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /accounts/auth/confirm-account [get]
func ConfirmAccount(c *gin.Context) {
	app := c.MustGet("app").(*settings.App)
	user := c.MustGet("user").(dto.UserDTO)
	session_id := c.Param("token")

	service := services.NewAuthService(app)
	sess, err := service.ConfirmAccount(c.Request.Context(), user, session_id)
	if err != nil {
		c.Error(err)
		return
	}
	c.SetCookie("sessionID", sess, int(app.Config.AuthConfig.AuthSessionTTL), "/", "", true, true)
	c.JSON(http.StatusOK, dto.MessageResponse{Message: "success"})
}

// @Summary Login
// @Description Login to account
// @Tags Auth
// @Accept json
// @Produce json
// @Param user body dto.LoginRequest true "Data"
// @Success 200 {object} dto.MessageResponse
// @Failure 400 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /accounts/auth/login [post]
func Login(c *gin.Context) {
	app := c.MustGet("app").(*settings.App)
	user := c.MustGet("user").(dto.UserDTO)

	var loginData dto.LoginRequest

	if err := c.ShouldBindJSON(&loginData); err != nil {
		c.Error(usecase_errors.BadRequestError{Msg: err.Error()})
		return
	}

	service := services.NewAuthService(app)

	sess, err := service.Login(c.Request.Context(), user, loginData)
	if err != nil {
		c.Error(err)
		return
	}
	c.SetCookie("sessionID", sess, int(app.Config.AuthConfig.AuthSessionTTL), "/", "", true, true)
	c.JSON(http.StatusOK, dto.MessageResponse{Message: "success"})
}

// @Summary Logout
// @Description Logout the session
// @Tags Auth
// @Accept json
// @Produce json
// @Success 200 {object} dto.MessageResponse
// @Failure 400 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /accounts/auth/logout [delete]
func Logout(c *gin.Context) {
	app := c.MustGet("app").(*settings.App)
	cookie, _ := c.Cookie("sessionID")

	services.NewAuthService(app).Logout(cookie)

	c.SetCookie("sessionID", "", -1, "/", "", true, true)
	c.JSON(200, dto.MessageResponse{Message: "success"})
}
