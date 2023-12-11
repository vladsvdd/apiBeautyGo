package handler

import (
	"github.com/gin-gonic/gin"
	"go_api/models"
	"net/http"
)

// signUp отвечает за обработку запроса на регистрацию пользователя
// @Summary SignUp
// @Tags auth
// @Description create account
// @ID create-account
// @Accept  json
// @Produce  json
// @Param input body models.UserSignUpInput true "account info"
// @Success 200 {integer} integer 1
// @Failure 400,404 {object} errorResponse
// @Failure 500 {object} errorResponse
// @Failure default {object} errorResponse
// @Router /go_api/auth/sign-up [post]
func (h *Handler) signUp(c *gin.Context) {
	var input models.UserSignUpInput

	if err := c.ShouldBindJSON(&input); err != nil {
		newErrorDataResponse(c, http.StatusBadRequest, GetErrorsMsg(err))
		return
	}

	id, err := h.services.Authorization.CreateUser(input)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusCreated, map[string]interface{}{
		"id": id,
	})
}

type signInInput struct {
	Email    string `json:"email" binding:"required"`
	Password string `json:"password" binding:"required"`
}

// signIn отвечает за обработку запроса на вход пользователя (аутентификация)
// @Summary SignIn
// @Tags auth
// @Description login
// @ID login
// @Accept  json
// @Produce  json
// @Param input body signInInput true "credentials"
// @Success 200 {string} string "token"
// @Failure 400,404 {object} errorResponse
// @Failure 500 {object} errorResponse
// @Failure default {object} errorResponse
// @Router /go_api/auth/sign-in [post]
func (h *Handler) signIn(c *gin.Context) {
	op := "pkg/handler/auth.go signIn(c *gin.Context) -> "
	var input signInInput

	if err := c.ShouldBindJSON(&input); err != nil {
		newErrorDataResponse(c, http.StatusBadRequest, GetErrorsMsg(err))
		return
	}

	token, err := h.services.Authorization.GenerateToken(input.Email, input.Password)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, op+err.Error())
		return
	}

	c.JSON(http.StatusOK, map[string]interface{}{
		"token": token,
	})
}
