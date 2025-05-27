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
