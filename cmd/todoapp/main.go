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
	os.Exit(start())
}

type configs struct {
	app    *core_config.Config
	logger core_logger.Config
	pool   core_pgx_pool.Config
	http   core_http_server.Config
}

func loadConfigs() (configs, error) {
	var (
		cfgs configs
		err  error
	)

	if cfgs.app, err = core_config.NewConfig(); err != nil {
		return configs{}, fmt.Errorf("app config: %w", err)
	}

	if cfgs.logger, err = core_logger.NewConfig(); err != nil {
		return configs{}, fmt.Errorf("logger config: %w", err)
	}

	if cfgs.pool, err = core_pgx_pool.NewConfig(); err != nil {
		return configs{}, fmt.Errorf("postgres config: %w", err)
	}

	if cfgs.http, err = core_http_server.NewConfig(); err != nil {
		return configs{}, fmt.Errorf("HTTP server config: %w", err)
	}

	return cfgs, nil
}

func start() int {
	cfgs, err := loadConfigs()
	if err != nil {
		fmt.Fprintln(os.Stderr, "failed to load config:", err)

		return 1
	}

	logger, err := core_logger.NewLogger(cfgs.logger)
	if err != nil {
		fmt.Fprintln(os.Stderr, "failed to init application logger:", err)

		return 1
	}
	defer logger.Close()

	if err := run(cfgs, logger); err != nil {
		logger.Error("application stopped with error", zap.Error(err))

		return 1
	}

	return 0
}

func run(cfgs configs, logger *core_logger.Logger) error {
	ctx, cancel := signal.NotifyContext(
		context.Background(),
		syscall.SIGINT, syscall.SIGTERM,
	)
	defer cancel()

	logger.Debug("application time zone", zap.Stringer("zone", cfgs.app.TimeZone))

	logger.Debug("initializing postgres connection pool")
	pool, err := core_pgx_pool.NewPool(ctx, cfgs.pool)
	if err != nil {
		return fmt.Errorf("init postgres connection pool: %w", err)
	}
	defer pool.Close()

	logger.Debug("initializing feature", zap.String("feature", "users"))
	usersRepository := users_postgres_repository.NewUsersRepository(pool)
	usersService := users_service.NewUsersService(usersRepository)
	usersTransportHTTP := users_transport_http.NewUsersHTTPHandler(usersService)

	logger.Debug("initializing feature", zap.String("feature", "tasks"))
	tasksRepository := tasks_postgres_repository.NewTasksRepository(pool)
	tasksService := tasks_service.NewTasksService(tasksRepository)
	tasksTransportHTTP := tasks_transport_http.NewTasksHTTPHandler(tasksService)

	logger.Debug("initializing feature", zap.String("feature", "statistics"))
	statisticsRepository := statistics_postgres_repository.NewStatisticsRepository(pool)
	statisticsService := statistics_service.NewStatisticsService(statisticsRepository)
	statisticsTransportHTTP := statistics_transport_http.NewStatisticsHTTPHandler(statisticsService, cfgs.app.TimeZone)

	logger.Debug("initializing feature", zap.String("feature", "web"))
	webRepository := web_fs_repository.NewWebRepository(public.FS)
	webService := web_service.NewWebService(webRepository)
	webTransportHTTP := web_transport_http.NewWebHTTPHandler(webService)

	logger.Debug("initializing HTTP server")
	httpConfig := cfgs.http
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

	apiVersionRouterV1 := core_http_server.NewAPIVersionRouter(core_http_server.APIVersion1)
	apiVersionRouterV1.RegisterRoutes(usersTransportHTTP.Routes()...)
	apiVersionRouterV1.RegisterRoutes(tasksTransportHTTP.Routes()...)
	apiVersionRouterV1.RegisterRoutes(statisticsTransportHTTP.Routes()...)

	httpServer.RegisterAPIRouters(
		apiVersionRouterV1,
	)
	httpServer.RegisterRoutes(webTransportHTTP.Routes()...)
	if httpConfig.SwaggerEnabled {
		httpServer.RegisterSwagger()
	}
	httpServer.RegisterHealthChecks(pool.Ping)

	if err := httpServer.Run(ctx); err != nil {
		return fmt.Errorf("run HTTP server: %w", err)
	}

	return nil
}
