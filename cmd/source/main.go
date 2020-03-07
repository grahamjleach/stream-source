package main

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/grahamjleach/stream-source/interfaces/logger/zap"
	server "github.com/grahamjleach/stream-source/interfaces/server/grpc"
	"github.com/grahamjleach/stream-source/interfaces/stream"
	"github.com/grahamjleach/stream-source/usecases"
	"github.com/stunndard/goicy/config"
	"github.com/stunndard/goicy/metadata"
	"google.golang.org/grpc"
)

var (
	version   string
	namespace = "stream_client_source"
)

func main() {
	backgroundContext, cancel := context.WithCancel(context.Background())
	defer cancel()

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

	var conn *grpc.ClientConn
	{
		var err error
		scheduleTarget := os.Getenv("SCHEDULE_TARGET")
		logger.Info("connecting to schedule", map[string]string{"target": scheduleTarget})
		if conn, err = grpc.Dial(scheduleTarget, grpc.WithInsecure()); err != nil {
			panic(err)
		}
	}

	logger.Info("initializing schedule client", nil)
	client := server.NewScheduleClient(conn)

	logger.Info("initializing schedule feed", nil)
	feed := NewScheduleFeed(backgroundContext, client)

	logger.Info("loading stream configuration", nil)
	if err := config.LoadConfig(os.Getenv("INI_FILE")); err != nil {
		panic(err)
	}

	logger.Info(
		"initializing stream client",
		map[string]string{
			"stream_format":        config.Cfg.StreamFormat,
			"stream_server_host":   config.Cfg.Host,
			"stream_server_port":   fmt.Sprintf("%d", config.Cfg.Port),
			"stream_server_hmount": config.Cfg.Mount,
			"ffmpeg_path":          config.Cfg.FFMPEGPath,
		},
	)
	streamClient, err := stream.NewClient(&config.Cfg)
	if err != nil {
		panic(err)
	}

	logger.Info(
		"consuming schedule feed",
		map[string]string{
			"current_offset": usecases.GetCurrentOffset().String(),
		},
	)

	go func() {
		for {
			item := feed.GetNextScheduleItem()

			metadata.SendMetadata(item.GetTitle())

			logger.Info(
				"streaming episode",
				map[string]string{
					"program_name":   item.GetProgram().Name,
					"host_name":      item.GetProgram().HostName,
					"episode_name":   item.GetEpisode().Name,
					"air_date":       item.GetEpisode().AirDate.Format(time.RFC3339),
					"recording_path": item.GetEpisode().RecordingPath,
				},
			)

			go func() {
				if err := streamClient.StreamFFMPEG(item.GetEpisode().RecordingPath); err != nil {
					panic(err)
				}
			}()
		}
	}()

	feed.LoadItems()
}
