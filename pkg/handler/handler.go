package handler

import (
	"github.com/gin-gonic/gin"      // Импорт пакета Gin для создания маршрутов и обработки HTTP-запросов.
	"github.com/swaggo/files"       // swagger embed files
	"github.com/swaggo/gin-swagger" // gin-swagger middleware
	_ "go_api/docs"                 // важно указать директорию к файлам docs для работы swagger
	"go_api/pkg/service"            // Импорт пакета service для обработки бизнес-логики.
)

// Handler представляет обработчик запросов
type Handler struct {
	services *service.Service
}

// NewHandler создает новый экземпляр обработчика с переданным сервисом.
func NewHandler(services *service.Service) *Handler {
	return &Handler{services: services}
}

// InitRoutes Метод инициализирует все маршруты для приложения и возвращает экземпляр маршрутизатора Gin.
func (h *Handler) InitRoutes() *gin.Engine {
	router := gin.New() // Создание нового экземпляра маршрутизатора Gin.

	router.GET("/go_api/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	auth := router.Group("/go_api/auth") // Группировка маршрутов для аутентификации.
	{
		auth.POST("/sign-up", h.signUp) // Обработка POST-запроса для регистрации пользователя.
		auth.POST("/sign-in", h.signIn) // Обработка POST-запроса для входа пользователя.
	}

	api := router.Group("/go_api/api/v1") // Группировка маршрутов для API.
	{
		meetings := api.Group("/...")
		{
			meetings.POST("/", h.createMeeting)
			meetings.GET("/", h.userIdentity, h.getAllMeeting)
		}

		apiProtected := api.Group("/", h.userIdentity)
		{
			companyDomains := apiProtected.Group("/..._domain")
			{
				companyDomains.POST("/", h.createDomains)
				companyDomains.GET("/", h.getAllDomains)
				companyDomains.GET("/:id", h.getDomainById)
				companyDomains.PUT("/:id", h.updateDomain)
				companyDomains.DELETE("/:id", h.deleteDomain)
			}
		}
	}

	return router // Возвращение настроенного маршрутизатора.
}
