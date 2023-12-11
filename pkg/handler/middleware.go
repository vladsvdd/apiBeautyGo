package handler

import (
	"errors"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
	"strings"
)

const (
	authorizationHeader = "Authorization" // Заголовок авторизации
	userContext         = "userId"        // Контекст пользователя
)

// userIdentity проверяет идентификационные данные пользователя в заголовке запроса
// Устанавливаем id пользователя в контекст
func (h *Handler) userIdentity(c *gin.Context) {
	header := c.GetHeader(authorizationHeader) // Получаем заголовок авторизации из запроса
	if header == "" {
		newErrorResponse(c, http.StatusUnauthorized, "empty auth header") // Отправляем ошибку, если заголовок пустой
		return
	}
	headerParts := strings.Split(header, " ") // Разделяем заголовок на части
	if len(headerParts) != 2 {
		newErrorResponse(c, http.StatusUnauthorized, "invalid auth header") // Отправляем ошибку, если заголовок недействителен
		return
	}

	userId, err := h.services.Authorization.ParseToken(headerParts[1]) // Извлекаем идентификатор пользователя из токена
	if err != nil {
		newErrorResponse(c, http.StatusUnauthorized, err.Error()) // Отправляем ошибку, если возникает ошибка при извлечении идентификатора
		return
	}

	c.Set(userContext, userId) // Устанавливаем идентификатор пользователя в контексте запроса
}

// getUserId извлекает идентификатор пользователя из контекста запроса, в виде числа
func getUserId(c *gin.Context) (int64, error) {
	id, ok := c.Get(userContext) // Получаем идентификатор пользователя из контекста запроса
	if !ok {
		newErrorResponse(c, http.StatusInternalServerError, "user is not found") // Отправляем ошибку, если пользователь не найден
		return 0, errors.New("user id not found")
	}

	idInt, ok := id.(int64) // Преобразуем идентификатор в целое число
	if !ok {
		newErrorResponse(c, http.StatusInternalServerError, "user id is of invalid type") // Отправляем ошибку, если тип идентификатора недействителен
		return 0, errors.New("user id is of invalid type")
	}

	return idInt, nil // Возвращаем идентификатор пользователя
}

// getPage извлекает идентификатор номера страницы из контекста запроса, в виде числа
func getPage(c *gin.Context) (int, error) {
	page := 1
	pageStr := c.Query("page")
	if pageStr != "" {
		p, err := strconv.Atoi(pageStr)
		if err != nil {
			newErrorResponse(c, http.StatusBadRequest, "invalid page value")
			return 0, errors.New("invalid page value")
		}

		if p > 0 {
			page = p
		}
	}

	return page, nil
}

// getLimit извлекает идентификатор лимита записей на странице из контекста запроса, в виде числа
func getLimit(c *gin.Context) (int, error) {
	limit := 1000
	limitStr := c.Query("limit")
	if limitStr != "" {
		l, err := strconv.Atoi(limitStr)
		if err != nil {
			newErrorResponse(c, http.StatusBadRequest, "invalid limit value")
			return 0, errors.New("invalid limit value")
		}
		if l > 1000 {
			limit = 1000
		} else {
			limit = l
		}
	}

	return limit, nil
}
