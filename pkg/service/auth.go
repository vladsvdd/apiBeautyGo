package service

import (
	"database/sql"
	"errors"
	"fmt"
	"github.com/sirupsen/logrus"
	"go_api/models"

	//"github.com/dgrijalva/jwt-go"
	"github.com/golang-jwt/jwt/v5"
	"go_api/pkg/repository"
	"golang.org/x/crypto/bcrypt"
	"os"
	"time"
)

// tokenTTL представляет продолжительность времени жизни токена.
const (
	tokenTTL = 12 * time.Hour
)

// tokenClaims содержит информацию, которая будет включена в токен.
type tokenClaims struct {
	UserId               *int64 `json:"user_id"` // Идентификатор пользователя
	jwt.RegisteredClaims        // Поля, обязательные для JWT
}

// AuthService представляет собой службу аутентификации.
type AuthService struct {
	repo repository.Authorization // Репозиторий аутентификации
}

// NewAuthService создает новый экземпляр сервиса аутентификации с заданным репозиторием.
func NewAuthService(repo repository.Authorization) *AuthService {
	return &AuthService{repo: repo}
}

// CreateUser принимает пользователя в качестве входного параметра и создает нового пользователя в базе данных.
// Пароль пользователя хешируется с помощью функции generatePasswordHash.
// Возвращает идентификатор созданного пользователя и ошибку, если что-то пошло не так при создании пользователя.
func (s *AuthService) CreateUser(user models.UserSignUpInput) (int64, error) {
	// Хеширование пароля пользователя
	hashedPassword, err := generatePasswordHash(user.Password)
	if err != nil {
		return 0, err
	}

	// Замена пароля пользователя на хешированный
	user.Password = hashedPassword

	// Создание нового пользователя в базе данных
	return s.repo.CreateUser(user)
}

// generatePasswordHash принимает пароль в виде строки и возвращает его хешированное значение с использованием bcrypt.
// Возвращает хешированный пароль и ошибку, если что-то пошло не так при хешировании.
func generatePasswordHash(password string) (string, error) {
	// Хеширование пароля с помощью bcrypt и установка стандартной стоимости
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)

	// Возвращение хешированного пароля в виде строки и ошибки, если таковая имеется
	return string(hashedPassword), err
}

// comparePasswords сравнивает хешированный пароль с переданным паролем для проверки подлинности.
func comparePasswords(hashedPassword string, password string) bool {
	// Сравнение хешированного пароля с переданным паролем
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))

	// Возвращение true, если пароли совпадают, иначе - false
	return err == nil
}

// GenerateToken генерирует JWT на основе предоставленной электронной почты и пароля.
func (s *AuthService) GenerateToken(username, password string) (string, error) {
	op := "pkg/service/auth.go GenerateToken(username, password string) -> "
	// Получение пользователя из базы данных
	user, err := s.repo.GetUser(username)
	if err != nil {
		logrus.Println(op + err.Error())
		if err == sql.ErrNoRows {
			return "", errors.New(fmt.Sprint("Неверный email или пароль"))
		}
		return "", errors.New(fmt.Sprint("Неизвестная ошибка при входе"))
	}

	// Сравнение хешированных паролей для проверки подлинности
	if comparePasswords(*user.Password, password) == false {
		return "", errors.New(fmt.Sprint("Неверный пароль"))
	}

	// Создание нового токена с заданными утверждениями, включая идентификатор пользователя и срок действия
	token := jwt.NewWithClaims(jwt.SigningMethodHS256,
		&tokenClaims{user.Id,
			jwt.RegisteredClaims{
				ExpiresAt: jwt.NewNumericDate(time.Now().Add(tokenTTL)),
				IssuedAt:  jwt.NewNumericDate(time.Now()),
			},
		})

	// Возврат подписанного токена в виде строки
	return token.SignedString([]byte(os.Getenv("SIGNED_KEY")))
}

// ParseToken принимает строку токена и возвращает идентификатор пользователя, связанный с этим токеном.
// Этот метод извлекает информацию из токена и проверяет его подлинность, используя указанный ключ.
// Возвращает идентификатор пользователя и ошибку в случае неудачи.
func (s *AuthService) ParseToken(accessToken string) (int64, error) {
	// Попытка разобрать токен с утверждениями, хранящимися в tokenClaims
	token, err := jwt.ParseWithClaims(accessToken, &tokenClaims{}, func(token *jwt.Token) (interface{}, error) {
		// Проверяем, является ли метод подписи HMAC
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("invalid signing method")
		}

		// Возвращаем подписанный ключ из переменной окружения SIGNED_KEY в качестве ключа для проверки токена
		return []byte(os.Getenv("SIGNED_KEY")), nil
	})

	// Если возникла ошибка при разборе токена, возвращаем 0 и ошибку
	if err != nil {
		return 0, err
	}

	// Проверяем, что утверждения токена являются экземпляром tokenClaims
	claims, ok := token.Claims.(*tokenClaims)
	if !ok {
		return 0, errors.New("token claims are not of type *tokenClaims")
	}

	// Возвращаем идентификатор пользователя из утверждений токена
	return *claims.UserId, nil
}
