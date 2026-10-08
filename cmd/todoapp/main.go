package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	_ "time/tzdata"

	core_config "github.com/MmGrand/TodoApp/internal/core/config"
	core_logger "github.com/MmGrand/TodoApp/internal/core/logger"
	core_pgx_pool "github.com/MmGrand/TodoApp/internal/core/repository/postgres/pool/pgx"
	core_http_middleware "github.com/MmGrand/TodoApp/internal/core/transport/http/middleware"
	core_http_server "github.com/MmGrand/TodoApp/internal/core/transport/http/server"
	statistics_postgres_repository "github.com/MmGrand/TodoApp/internal/features/statistics/repository/postgres"
	statistics_service "github.com/MmGrand/TodoApp/internal/features/statistics/service"
	statistics_transport_http "github.com/MmGrand/TodoApp/internal/features/statistics/transport/http"
	tasks_postgres_repository "github.com/MmGrand/TodoApp/internal/features/tasks/repository/postgres"
	tasks_service "github.com/MmGrand/TodoApp/internal/features/tasks/service"
	tasks_transport_http "github.com/MmGrand/TodoApp/internal/features/tasks/transport/http"
	users_postgres_repository "github.com/MmGrand/TodoApp/internal/features/users/repository/postgres"
	users_service "github.com/MmGrand/TodoApp/internal/features/users/service"
	users_transport_http "github.com/MmGrand/TodoApp/internal/features/users/transport/http"
	web_fs_repository "github.com/MmGrand/TodoApp/internal/features/web/repository/file_system"
	web_service "github.com/MmGrand/TodoApp/internal/features/web/service"
	web_transport_http "github.com/MmGrand/TodoApp/internal/features/web/transport"
	"github.com/MmGrand/TodoApp/public"
	"go.uber.org/zap"
)

// @title 			Golang Todo API
// @version 		1.0
// @description 	Todo Application REST-API scheme
// @BasePath 		/api/v1
func main() {
	cfg := core_config.NewConfigMust()

	logger, err := core_logger.NewLogger(core_logger.NewConfigMust())
	if err != nil {
		fmt.Println("failed to init application logger: ", err)
		os.Exit(1)
	}

	err = run(cfg, logger)
	logger.Close()

	if err != nil {
		os.Exit(1)
	}
}

func run(cfg *core_config.Config, logger *core_logger.Logger) error {
	ctx, cancel := signal.NotifyContext(
		context.Background(),
		syscall.SIGINT, syscall.SIGTERM,
	)
	defer cancel()

	logger.Debug("application time zone", zap.Stringer("zone", cfg.TimeZone))

	logger.Debug("initializing postgres connection pool")
	pool, err := core_pgx_pool.NewPool(
		ctx,
		core_pgx_pool.NewConfigMust(),
	)
	if err != nil {
		logger.Error("failed to init postgres connection pool", zap.Error(err))
		return err
	}
	defer pool.Close()

	logger.Debug("initializing feature", zap.String("feature", "users"))
	usersRepository := users_postgres_repository.NewUsersRepository(pool)
	usersService := users_service.NewUserService(usersRepository)
	usersTransportHTTP := users_transport_http.NewUsersHTTPHandler(usersService)

	logger.Debug("initializing feature", zap.String("feature", "tasks"))
	tasksRepository := tasks_postgres_repository.NewTasksRepository(pool)
	tasksService := tasks_service.NewTasksService(tasksRepository)
	tasksTransportHTTP := tasks_transport_http.NewTasksHTTPHandler(tasksService)

	logger.Debug("initializing feature", zap.String("feature", "statistics"))
	statisticsRepository := statistics_postgres_repository.NewStatisticsRepository(pool)
	statisticsService := statistics_service.NewStatisticsService(statisticsRepository)
	statisticsTransportHTTP := statistics_transport_http.NewStatisticsHTTPHandler(statisticsService, cfg.TimeZone)

	logger.Debug("initializing feature", zap.String("feature", "web"))
	webRepository := web_fs_repository.NewWebRepository(public.FS)
	webService := web_service.NewWebService(webRepository)
	webTransportHTTP := web_transport_http.NewWebHTTPHandler(webService)

	logger.Debug("initializing HTTP server")
	httpConfig := core_http_server.NewConfigMust()
	httpServer := core_http_server.NewHTTPServer(
		httpConfig,
		logger,
		core_http_middleware.Panic(logger),
		core_http_middleware.SecurityHeaders(),
		core_http_middleware.CORS(httpConfig.AllowedOrigins),
		core_http_middleware.RequestID(),
		core_http_middleware.Logger(logger),
		core_http_middleware.Trace(),
		core_http_middleware.BodyLimit(httpConfig.MaxBodyBytes),
	)

	apiVersionRouterV1 := core_http_server.NewApiVersionRouter(core_http_server.ApiVersion1)
	apiVersionRouterV1.RegisterRoutes(usersTransportHTTP.Routes()...)
	apiVersionRouterV1.RegisterRoutes(tasksTransportHTTP.Routes()...)
	apiVersionRouterV1.RegisterRoutes(statisticsTransportHTTP.Routes()...)

	httpServer.RegisterApiRouters(
		apiVersionRouterV1,
	)
	httpServer.RegisterRoutes(webTransportHTTP.Routes()...)
	if httpConfig.SwaggerEnabled {
		httpServer.RegisterSwagger()
	}
	httpServer.RegisterHealthCheck(pool.Ping)

	if err := httpServer.Run(ctx); err != nil {
		logger.Error("HTTP server run error", zap.Error(err))
		return err
	}

	return nil
}
