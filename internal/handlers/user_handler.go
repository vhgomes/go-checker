package handlers

import (
	"go-checker/internal/repository"
	"go-checker/internal/utils"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

var validate = validator.New()

type UserHandler struct {
	userRepo *repository.UserRepo
}

func NewUserHandler(userRepo *repository.UserRepo) *UserHandler {
	return &UserHandler{userRepo: userRepo}
}

type registerRequest struct {
	Email    string `json:"email"    validate:"required,email,max=254"`
	Password string `json:"password" validate:"required,min=8,max=72"`
	Name     string `json:"name"     validate:"required,min=2,max=100"`
}

type loginRequest struct {
	Email    string `json:"email"    validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

func (h *UserHandler) RegisterUser(c *gin.Context) {
	var body registerRequest

	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "requisição inválida"})
		return
	}

	if err := validate.Struct(body); err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"error": formatValidationError(err)})
		return
	}

	if err := h.userRepo.CreateUser(c.Request.Context(), body.Email, body.Password, body.Name); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"msg": "usuário criado com sucesso"})
}

func (h *UserHandler) Login(c *gin.Context) {
	var body loginRequest

	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "requisição inválida"})
		return
	}

	if err := validate.Struct(body); err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"error": formatValidationError(err)})
		return
	}

	user, err := h.userRepo.Login(c.Request.Context(), body.Email, body.Password)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	token, err := utils.GenerateToken(user.ID, user.Email)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "falha ao gerar token"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"msg":   "login realizado com sucesso",
		"token": token,
	})
}

func formatValidationError(err error) string {
	var ve validator.ValidationErrors
	if ok := (err.(validator.ValidationErrors)); ok != nil {
		ve = ok
	}
	if len(ve) == 0 {
		return "dados inválidos"
	}
	fe := ve[0]
	switch fe.Tag() {
	case "required":
		return "campo '" + fe.Field() + "' é obrigatório"
	case "email":
		return "email inválido"
	case "min":
		return "campo '" + fe.Field() + "' deve ter pelo menos " + fe.Param() + " caracteres"
	case "max":
		return "campo '" + fe.Field() + "' deve ter no máximo " + fe.Param() + " caracteres"
	default:
		return "campo '" + fe.Field() + "' inválido"
	}
}
