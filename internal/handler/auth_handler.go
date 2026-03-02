package handler

import (
	"fmt"
	"net/http"
	"time"

	"tung-tung-money-api/internal/middleware"
	"tung-tung-money-api/internal/model"
	"tung-tung-money-api/internal/router"
	"tung-tung-money-api/internal/service"

	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
)

type AuthHandler struct {
	userService *service.UserService
	jwtSecret   string
}

func NewAuthHandler(userService *service.UserService, jwtSecret string) *AuthHandler {
	return &AuthHandler{
		userService: userService,
		jwtSecret:   jwtSecret,
	}
}

func (h *AuthHandler) RegisterRoutes() []router.Route {
	return []router.Route{
		{Method: "POST", Path: "/auth/register", Handler: h.Register},
		{Method: "POST", Path: "/auth/login", Handler: h.Login},
	}
}

// Register godoc
// @Summary      Register a new user
// @Description  Create a new user account
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        request body model.RegisterRequest true "Register request"
// @Success      201 {object} model.User
// @Failure      400 {object} map[string]string
// @Router       /auth/register [post]
func (h *AuthHandler) Register(c echo.Context) error {
	var req model.RegisterRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid request body"})
	}

	if err := c.Validate(req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}

	user, err := h.userService.Register(c.Request().Context(), &req)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusCreated, user)
}

// Login godoc
// @Summary      Login
// @Description  Authenticate user and return JWT token
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        request body model.LoginRequest true "Login request"
// @Success      200 {object} map[string]string
// @Failure      401 {object} map[string]string
// @Router       /auth/login [post]
func (h *AuthHandler) Login(c echo.Context) error {
	var req model.LoginRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid request body"})
	}

	if err := c.Validate(req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}

	user, err := h.userService.Login(c.Request().Context(), &req)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, map[string]string{"error": err.Error()})
	}

	claims := &middleware.JWTClaims{
		UserID: fmt.Sprintf("%d", user.ID),
		Email:  user.Email,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenStr, err := token.SignedString([]byte(h.jwtSecret))
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "failed to generate token"})
	}

	return c.JSON(http.StatusOK, map[string]string{"token": tokenStr})
}
