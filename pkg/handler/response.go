package handler

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/sirupsen/logrus"
)

type errorResponse struct {
	Message string `json:"message"`
}

type ErrorMsg struct {
	Field   string `json:"field"`
	Message string `json:"message"`
}

func newErrorResponse(c *gin.Context, statusCode int, message string) {
	logrus.Errorf(message)
	c.AbortWithStatusJSON(statusCode, gin.H{"errors": errorResponse{message}})
}

// Формируем json с данными по ошибкам
func newErrorDataResponse(c *gin.Context, statusCode int, out []ErrorMsg) {
	jsonData, err := json.Marshal(out)
	if err != nil {
		logrus.Errorf(fmt.Sprintln("Ошибка при маршалинге в JSON: ", err, out))
	} else {
		logrus.Errorf(fmt.Sprintln(string(jsonData)))
	}
	c.AbortWithStatusJSON(statusCode, gin.H{"errors": out})
}

func GetErrorMsg(fe validator.FieldError) string {
	switch fe.Tag() {
	case "required":
		return "This field is required"
	case "lte":
		return "Should be less than " + fe.Param()
	case "gte":
		return "Should be greater than " + fe.Param()
	case "email":
		return "Not valid email " + fe.Param()
	case "eqfield":
		return "Поля не сопадают " + fe.Param()
	}
	return "Unknown error"
}

func GetErrorsMsg(err error) []ErrorMsg {
	var ve validator.ValidationErrors
	if errors.As(err, &ve) {
		out := make([]ErrorMsg, len(ve))
		for i, fe := range ve {
			out[i] = ErrorMsg{Field: fe.Field(), Message: GetErrorMsg(fe)}
		}
		return out
	}
	return nil
}

type statusResponse struct {
	Status string `json:"status"`
}
