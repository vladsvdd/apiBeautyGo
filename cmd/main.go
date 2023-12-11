package main

import (
	"context"
	"github.com/joho/godotenv"   // Пакет godotenv загружает переменные среды из файла .env.
	_ "github.com/lib/pq"        // Драйвер PostgreSQL.
	"github.com/sirupsen/logrus" // Пакет logrus для логирования.
	"github.com/spf13/viper"     // Пакет viper для работы с конфигурационными файлами.
	"go_api"
	"go_api/pkg/handler"    // Пакет handler для обработки HTTP-запросов.
	"go_api/pkg/repository" // Пакет repository для работы с базой данных.
	"go_api/pkg/service"    // Пакет service для бизнес-логики.
	"os"                    // Пакет os для работы с операционной системой.
	"os/signal"
	"syscall"
)

// @title Beauty App API
// @version 1.0
// @description API Server for Beauty Application
// @host localhost:1111
// @BasePath /
// @securityDefinitions.apikey ApiKeyAuth
// @in header
// @name Authorization
func main() {
	// Настройка логгера
	setupLogger()

	// Инициализация конфигурационного файла
	if err := initConfig(); err != nil {
		logrus.Fatalf("Ошибка инициализации конфиг файла %s", err.Error())
	}

	// Загрузка переменных среды из файла .env
	if err := godotenv.Load(); err != nil {
		logrus.Fatalf("Ошибка при загрузке файла .env: %s", err)
	}

	// Инициализация подключения к базе данных PostgreSQL
	db, err := repository.NewPostgresDB(repository.Config{
		Host:     viper.GetString("db.host"),
		Port:     viper.GetString("db.port"),
		Username: viper.GetString("db.username"),
		DBName:   viper.GetString("db.dbname"),
		SSLMode:  viper.GetString("db.sslmode"),
		Password: os.Getenv("DB_PASSWORD"),
	})

	if err != nil {
		logrus.Fatalf("ошибка инициализации подключения к базе данных: %s", err.Error())
	}

	// Инициализация репозиториев, сервисов и обработчиков
	repos := repository.NewRepository(db)    // ↑3)Работа с БД
	services := service.NewService(repos)    // ↑2)Бизнес логика
	handlers := handler.NewHandler(services) // ↑1)Работа с HTTP

	// Создание экземпляра сервера и запуск сервера в отдельной горутине
	srv := new(go_api.Server)
	go func() {
		if err = srv.Run(viper.GetString("port"), handlers.InitRoutes()); err != nil {
			logrus.Fatalf("ОШИБКА при запуске http-сервера: %s", err.Error())
		}
	}()
	logrus.Println("Beauty API начато")

	// Ожидание сигналов SIGTERM или SIGINT для завершения работы
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGTERM, syscall.SIGINT)
	<-quit

	// Завершение работы сервера
	logrus.Print("Beauty API завершено")
	if err = srv.Shutdown(context.Background()); err != nil {
		logrus.Errorf("ошибка при завершении работы сервера: %s", err.Error())
	}

	// Закрытие соединения с базой данных
	if err = db.Close(); err != nil {
		logrus.Errorf("ошибка при закрытии соединения с базой данных: %s", err.Error())
	}
}

// Метод initConfig используется для инициализации конфигурационного файла
func initConfig() error {
	viper.AddConfigPath("configs")
	viper.SetConfigName("config")
	return viper.ReadInConfig()
}

// setupLogger создает логгер и настраивает форматирование журнала с использованием logrus
func setupLogger() {
	logFile, err := os.OpenFile("logs.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		logrus.Fatalf("Не удалось открыть файл журнала: %s", err.Error())
	}
	logrus.SetOutput(logFile)
	logrus.SetFormatter(&logrus.JSONFormatter{})
}
