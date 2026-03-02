package handler

import (
	"net/http"
	"strconv"

	custommw "tung-tung-money-api/internal/middleware"
	"tung-tung-money-api/internal/model"
	"tung-tung-money-api/internal/router"
	"tung-tung-money-api/internal/service"

	"github.com/labstack/echo/v4"
)

type UserHandler struct {
	userService *service.UserService
	jwtSecret   string
}

func NewUserHandler(userService *service.UserService, jwtSecret string) *UserHandler {
	return &UserHandler{
		userService: userService,
		jwtSecret:   jwtSecret,
	}
}

func (h *UserHandler) RegisterRoutes() []router.Route {
	const usersPath = "/api/users"
	const userByIDPath = "/api/users/:id"

	jwtMw := custommw.JWTAuth(h.jwtSecret)
	return []router.Route{
		{Method: "GET", Path: usersPath, Handler: h.GetAll, Middlewares: []echo.MiddlewareFunc{jwtMw}},
		{Method: "GET", Path: userByIDPath, Handler: h.GetByID, Middlewares: []echo.MiddlewareFunc{jwtMw}},
		{Method: "PUT", Path: userByIDPath, Handler: h.Update, Middlewares: []echo.MiddlewareFunc{jwtMw}},
		{Method: "DELETE", Path: userByIDPath, Handler: h.Delete, Middlewares: []echo.MiddlewareFunc{jwtMw}},
	}
}

// GetAll godoc
// @Summary      List all users
// @Description  Get all users
// @Tags         users
// @Produce      json
// @Security     BearerAuth
// @Success      200 {array} model.User
// @Failure      500 {object} map[string]string
// @Router       /api/users [get]
func (h *UserHandler) GetAll(c echo.Context) error {
	users, err := h.userService.GetAll(c.Request().Context())
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}
	return c.JSON(http.StatusOK, users)
}

// GetByID godoc
// @Summary      Get user by ID
// @Description  Get a single user by ID
// @Tags         users
// @Produce      json
// @Security     BearerAuth
// @Param        id path int true "User ID"
// @Success      200 {object} model.User
// @Failure      400 {object} map[string]string
// @Failure      404 {object} map[string]string
// @Router       /api/users/{id} [get]
func (h *UserHandler) GetByID(c echo.Context) error {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid id"})
	}

	user, err := h.userService.GetByID(c.Request().Context(), uint(id))
	if err != nil {
		return c.JSON(http.StatusNotFound, map[string]string{"error": "user not found"})
	}

	return c.JSON(http.StatusOK, user)
}

// Update godoc
// @Summary      Update user
// @Description  Update user by ID
// @Tags         users
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        id path int true "User ID"
// @Param        request body model.UpdateUserRequest true "Update request"
// @Success      200 {object} map[string]string
// @Failure      400 {object} map[string]string
// @Router       /api/users/{id} [put]
func (h *UserHandler) Update(c echo.Context) error {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid id"})
	}

	var req model.UpdateUserRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid request body"})
	}

	if err := h.userService.Update(c.Request().Context(), uint(id), &req); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, map[string]string{"message": "user updated"})
}

// Delete godoc
// @Summary      Delete user
// @Description  Delete user by ID
// @Tags         users
// @Produce      json
// @Security     BearerAuth
// @Param        id path int true "User ID"
// @Success      200 {object} map[string]string
// @Failure      400 {object} map[string]string
// @Router       /api/users/{id} [delete]
func (h *UserHandler) Delete(c echo.Context) error {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid id"})
	}

	if err := h.userService.Delete(c.Request().Context(), uint(id)); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, map[string]string{"message": "user deleted"})
}
