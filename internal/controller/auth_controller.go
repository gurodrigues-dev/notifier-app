package controller

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/gurodrigues-dev/notifier-app/internal/entity"
	"github.com/gurodrigues-dev/notifier-app/internal/infra"
	"github.com/gurodrigues-dev/notifier-app/internal/infra/contracts"
	"github.com/gurodrigues-dev/notifier-app/internal/usecase"
	"github.com/gurodrigues-dev/notifier-app/pkg/stringcommon"
	"github.com/jinzhu/gorm"
)

type AuthController struct {
	logger contracts.Logger
}

func NewAuthController(
	logger contracts.Logger,
) *AuthController {
	return &AuthController{
		logger: logger,
	}
}

// CreateToken godoc
// @Summary Create a new token
// @Description Creates a new authentication token based on the provided user data.
// @Tags auth
// @Accept json
// @Produce json
// @Param token body entity.Token true "Token request body"
// @Success 201 {object} entity.Token "Token created successfully"
// @Failure 400 {object} map[string]string "Invalid request body"
// @Failure 500 {object} map[string]string "Internal server error"
// @Router /token [post]
func (ac *AuthController) CreateToken(httpContext *gin.Context) {
	var requestParams entity.Token
	if err := httpContext.BindJSON(&requestParams); err != nil {
		httpContext.JSON(http.StatusBadRequest, gin.H{"error": "invalid body content"})
		return
	}

	validate := validator.New()
	if err := validate.Struct(requestParams); err != nil {
		httpContext.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	usecase := usecase.NewCreateTokenUsecase(
		infra.App.Repositories.AuthRepository,
		ac.logger,
	)

	token, err := usecase.CreateToken(&requestParams)
	if err != nil {
		httpContext.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	httpContext.JSON(http.StatusCreated, token)
	return
}

// GetToken godoc
// @Summary Get token by user
// @Description Retrieves the authentication token for a specified user.
// @Tags auth
// @Produce json
// @Param user path string true "User identifier"
// @Success 200 {object} entity.Token "Token retrieved successfully"
// @Failure 400 {object} map[string]string "User parameter is required"
// @Failure 404 {object} map[string]string "Token not found"
// @Failure 500 {object} map[string]string "Internal server error"
// @Router /token/{user} [get]
func (ac *AuthController) GetToken(httpContext *gin.Context) {
	user := httpContext.Param("user")
	if stringcommon.Empty(user) {
		httpContext.JSON(http.StatusBadRequest, gin.H{"error": "user is required"})
		return
	}

	usecase := usecase.NewGetTokenByUserUsecase(
		infra.App.Repositories.AuthRepository,
		ac.logger,
	)

	token, err := usecase.GetTokenByUser(user)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		httpContext.JSON(http.StatusNotFound, gin.H{"error": "token not found"})
		return
	}

	if err != nil {
		httpContext.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		return
	}

	httpContext.JSON(http.StatusOK, token)
}

// DeleteToken godoc
// @Summary Delete token by user
// @Description Deletes the authentication token for a specified user.
// @Tags auth
// @Produce json
// @Param user path string true "User identifier"
// @Success 200 {object} map[string]string "Token deleted successfully"
// @Failure 400 {object} map[string]string "User parameter is required"
// @Failure 404 {object} map[string]string "Token not found"
// @Failure 500 {object} map[string]string "Internal server error"
// @Router /token/{user} [delete]
func (ac *AuthController) DeleteToken(httpContext *gin.Context) {
	user := httpContext.Param("user")
	if stringcommon.Empty(user) {
		httpContext.JSON(http.StatusBadRequest, gin.H{"error": "user is required"})
		return
	}

	usecase := usecase.NewDeleteTokenUsecase(
		infra.App.Repositories.AuthRepository,
		ac.logger,
	)

	err := usecase.DeleteToken(user)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		httpContext.JSON(http.StatusNotFound, gin.H{"error": "token not found"})
		return
	}

	if err != nil {
		httpContext.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		return
	}

	httpContext.JSON(http.StatusOK, gin.H{"message": "token deleted ssuccessfully"})
}
