package main

import (
	"context"
	"fmt"
	"math"
	"net"
	"os"

	"github.com/grahamjleach/stream-source/service"
	server "github.com/grahamjleach/stream-source/interfaces/server/grpc"
	"github.com/grahamjleach/stream-source/interfaces/logger/zap"
	"github.com/grahamjleach/stream-source/interfaces/store"
	"github.com/grahamjleach/stream-source/interfaces/store/json"
	"google.golang.org/grpc"
)

var (
	version   string
	namespace = "stream_client_schedule"
)

func main() {
	var logger *zap.Logger
	{
		var err error
		loggerConfig := &zap.LoggerConfig{
			LogLevel:      os.Getenv("LOG_LEVEL"),
			CaptureStdLog: true,
			Context: map[string]interface{}{
				"application": namespace,
				"version":     version,
			},
		}
		if logger, err = zap.NewJSONLogger(loggerConfig); err != nil {
			panic(err)
		}
	}

	var scheduleStore store.ScheduleStore
	{
		var err error
		scheduleJSONPath := os.Getenv("SCHEDULE_JSON_PATH")
		logger.Info("initializing schedule store", map[string]string{"json_path": scheduleJSONPath})
		if scheduleStore, err = json.New(scheduleJSONPath); err != nil {
			panic(err)
		}
	}

	logger.Info("initializing schedule service", nil)
	service, err := service.NewScheduleService(scheduleStore)
	if err != nil {
		panic(err)
	}

	logger.Info("initializing schedule server handler", nil)
	handler := server.NewScheduleHandler(service)

	logger.Info("initializing schedule server", nil)
	grpcServer := grpc.NewServer(
		grpc.MaxConcurrentStreams(math.MaxUint32),
		grpc.UnaryInterceptor(NewLoggerMiddleware(logger)),
	)

	logger.Info("registering schedule server handler", nil)
	server.RegisterScheduleServer(grpcServer, handler)

	var tcpListener net.Listener
	{
		var err error
		grpcBindPort := os.Getenv("GRPC_BIND_PORT")
		logger.Info("intializing tcp listener", map[string]string{"port": grpcBindPort})
		if tcpListener, err = net.Listen("tcp", fmt.Sprintf(":%s", grpcBindPort)); err != nil {
			panic(err)
		}
	}

	logger.Info("starting server", nil)
	if err := grpcServer.Serve(tcpListener); err != nil {
		panic(err)
	}
}

func NewLoggerMiddleware(logger *zap.Logger) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, request interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (response interface{}, err error) {

		logger.Info(info.FullMethod, nil)

		return handler(ctx, request)
	}
}
