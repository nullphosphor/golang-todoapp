package main

import (
	"context"
	"log"
	"os/signal"
	"syscall"
	"time"

	core_logger "github.com/nullphosphor/golang-todoapp/internal/core/logger"
	core_pgx_pool "github.com/nullphosphor/golang-todoapp/internal/core/repository/postgres/pool/pgx"
	core_http_middleware "github.com/nullphosphor/golang-todoapp/internal/core/transport/http/middleware"
	core_http_server "github.com/nullphosphor/golang-todoapp/internal/core/transport/http/server"
	tasks_postgres_repository "github.com/nullphosphor/golang-todoapp/internal/features/tasks/repository/postgres"
	tasks_service "github.com/nullphosphor/golang-todoapp/internal/features/tasks/service"
	tasks_transport_http "github.com/nullphosphor/golang-todoapp/internal/features/tasks/transport/http"
	users_postgres_repository "github.com/nullphosphor/golang-todoapp/internal/features/users/repository/postgres"
	users_service "github.com/nullphosphor/golang-todoapp/internal/features/users/service"
	users_transport_http "github.com/nullphosphor/golang-todoapp/internal/features/users/transport/http"
	"go.uber.org/zap"
)

var (
	timeZone = time.UTC
)

func main() {
	time.Local = timeZone

	ctx, cancel := signal.NotifyContext(
		context.Background(),
		syscall.SIGINT, syscall.SIGTERM,
	)
	defer cancel()
	logger, err := core_logger.NewLogger(core_logger.NewConfigMust())
	if err != nil {
		log.Fatalf("failed to init application router: %v", err)
	}
	defer logger.Close()

	logger.Debug("application time zone", zap.Any("zone", timeZone))

	logger.Debug("initializing postgres connection pool")
	pool, err := core_pgx_pool.NewPool(
		ctx, core_pgx_pool.NewConfigMust(),
	)
	if err != nil {
		logger.Fatal("failed to init postgres connection pool", zap.Error(err))
	}
	defer pool.Close()

	logger.Debug("initializing feature", zap.String("feature", "users"))
	usersRepository := users_postgres_repository.NewUsersRepository(pool)
	usersService := users_service.NewUsersService(usersRepository)
	usersTrasportHTTP := users_transport_http.NewUsersHTTPhandler(usersService)

	logger.Debug("initializing feature", zap.String("feature", "tasks"))
	tasksRepository := tasks_postgres_repository.NewTasksRepository(pool)
	tasksService := tasks_service.NewTasksService(tasksRepository)
	tasksTrasportHTTP := tasks_transport_http.NewTasksHTTPHandler(tasksService)

	logger.Debug("initializing HTTP server")
	httpServer := core_http_server.NewHTTPServer(
		core_http_server.NewConfigMust(),
		logger,
		core_http_middleware.RequestID(),
		core_http_middleware.Logger(logger),
		core_http_middleware.Trace(),
		core_http_middleware.Panic(),
	)

	apiVersionRouter := core_http_server.NewAPIVersionRouter(core_http_server.ApiVersion1)
	apiVersionRouter.RegisterRoutes(usersTrasportHTTP.Routes()...)
	apiVersionRouter.RegisterRoutes(tasksTrasportHTTP.Routes()...)
	httpServer.RegisterAPIRouters(apiVersionRouter)

	if err := httpServer.Run(ctx); err != nil {
		logger.Error("start HTTP server", zap.Error(err))
	}
}
