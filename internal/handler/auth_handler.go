package handler

import (
	"errors"
	"net/http"

	"golang-test/internal/dto"
	"golang-test/internal/middleware"
	"golang-test/internal/service"
	"golang-test/pkg/response"

	"github.com/gin-gonic/gin"
)

type AuthHandler struct {
	authService service.AuthService
}

func NewAuthHandler(authService service.AuthService) *AuthHandler {
	return &AuthHandler{authService: authService}
}

func (h *AuthHandler) Register(c *gin.Context) {
	var req dto.RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.ValidationError(c, err.Error())
		return
	}

	res, err := h.authService.Register(req)
	if err != nil {
		if errors.Is(err, service.ErrEmailAlreadyExists) {
			response.Error(c, http.StatusConflict, err.Error(), nil)
			return
		}
		response.Error(c, http.StatusInternalServerError, "Failed to register user", nil)
		return
	}

	response.Success(c, http.StatusCreated, "User registered successfully", res)
}

func (h *AuthHandler) Login(c *gin.Context) {
	var req dto.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.ValidationError(c, err.Error())
		return
	}

	res, err := h.authService.Login(req)
	if err != nil {
		if errors.Is(err, service.ErrInvalidCredentials) {
			response.Error(c, http.StatusUnauthorized, err.Error(), nil)
			return
		}
		response.Error(c, http.StatusInternalServerError, "Failed to login", nil)
		return
	}

	response.Success(c, http.StatusOK, "Login successful", res)
}

func (h *AuthHandler) RefreshToken(c *gin.Context) {
	var req dto.RefreshTokenRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.ValidationError(c, err.Error())
		return
	}

	res, err := h.authService.RefreshToken(req)
	if err != nil {
		if errors.Is(err, service.ErrInvalidToken) || errors.Is(err, service.ErrUserNotFound) {
			response.Error(c, http.StatusUnauthorized, err.Error(), nil)
			return
		}
		response.Error(c, http.StatusInternalServerError, "Failed to refresh token", nil)
		return
	}

	response.Success(c, http.StatusOK, "Token refreshed successfully", res)
}

func (h *AuthHandler) Logout(c *gin.Context) {
	var req dto.RefreshTokenRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.ValidationError(c, err.Error())
		return
	}

	if err := h.authService.Logout(req.RefreshToken); err != nil {
		response.Error(c, http.StatusInternalServerError, "Failed to logout", nil)
		return
	}

	response.Success(c, http.StatusOK, "Logout successful", nil)
}

func (h *AuthHandler) Me(c *gin.Context) {
	userID, ok := middleware.GetUserID(c)
	if !ok {
		response.Error(c, http.StatusUnauthorized, "Unauthorized", nil)
		return
	}

	user, err := h.authService.GetProfile(userID)
	if err != nil {
		if errors.Is(err, service.ErrUserNotFound) {
			response.Error(c, http.StatusNotFound, err.Error(), nil)
			return
		}
		response.Error(c, http.StatusInternalServerError, "Failed to get profile", nil)
		return
	}

	response.Success(c, http.StatusOK, "Profile retrieved successfully", user)
}

func (h *AuthHandler) GoogleLogin(c *gin.Context) {
	url := h.authService.GetGoogleAuthURL()
	c.Redirect(http.StatusTemporaryRedirect, url)
}

func (h *AuthHandler) GoogleCallback(c *gin.Context) {
	code := c.Query("code")
	if code == "" {
		response.Error(c, http.StatusBadRequest, "Authorization code is required", nil)
		return
	}

	res, err := h.authService.HandleGoogleCallback(c.Request.Context(), code)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, "Failed to authenticate with Google: "+err.Error(), nil)
		return
	}

	response.Success(c, http.StatusOK, "Google authentication successful", res)
}
