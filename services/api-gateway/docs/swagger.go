package docs

// @title Shiroyama API Gateway
// @version 1.0
// @description REST API Gateway для микросервисной системы управления проектами Shiroyama. API Gateway предоставляет единую точку входа для всех клиентских приложений и управляет взаимодействием с внутренними gRPC сервисами.

// @contact.name Shiroyama Support
// @contact.email support@shiroyama.dev

// @host localhost:8080
// @BasePath /

// @securityDefinitions.apikey Bearer
// @in header
// @name Authorization
// @description Введите токен в формате: Bearer YOUR_JWT_TOKEN

// @tag.name auth
// @tag.description Операции аутентификации и авторизации

// @tag.name users
// @tag.description Управление пользователями

// @tag.name teams
// @tag.description Управление командами и участниками

// @tag.name boards
// @tag.description Управление досками и списками

// @tag.name tasks
// @tag.description Управление задачами

// @tag.name labels
// @tag.description Управление метками

// @tag.name comments
// @tag.description Комментарии и реакции

// @tag.name activity
// @tag.description Журнал активности

// @tag.name health
// @tag.description Проверки состояния системы
