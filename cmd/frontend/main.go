package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/go-chi/chi"
	chimiddleware "github.com/go-chi/chi/middleware"
	//"github.com/go-chi/jwtauth"
	"github.com/grahamjleach/stream-source/interfaces/logger/zap"
	transport "github.com/grahamjleach/stream-source/interfaces/server/grpc"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
)

var (
	version   string
	namespace = "stream_client_frontend"
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
	client := transport.NewScheduleClient(conn)

	logger.Info("initializing schedule frontend controller", nil)
	controller := newScheduleFrontendController(client)

	logger.Info("initializing schedule frontend router", nil)
	router := chi.NewRouter()
	router.Use(chimiddleware.Recoverer)
	router.Get("/health", statusOK)
	router.Post("/programs", controller.CreateProgram)
	router.Post("/programs/{program_id}/episodes", controller.CreateProgramEpisode)
	router.Post("/programs/{program_id}/schedule", controller.CreateProgramScheduleEntry)

	{
		httpBindPort := os.Getenv("HTTP_BIND_PORT")
		logger.Info("listening", map[string]string{"port": httpBindPort})
		if err := http.ListenAndServe(fmt.Sprintf(":%s", httpBindPort), router); err != nil {
			panic(err)
		}
	}
}

func statusOK(w http.ResponseWriter, _ *http.Request) {
	w.WriteHeader(http.StatusOK)
}

type ScheduleFrontendController struct {
	client transport.ScheduleClient
}

func newScheduleFrontendController(client transport.ScheduleClient) *ScheduleFrontendController {
	return &ScheduleFrontendController{
		client: client,
	}
}

func (c *ScheduleFrontendController) CreateProgram(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	parameters := new(transport.PutProgramRequest)
	if err := json.NewDecoder(r.Body).Decode(&parameters); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(errorResponse(err))
		return
	}

	program, err := c.client.PutProgram(r.Context(), parameters)
	switch {
	case err != nil:
		w.WriteHeader(grpcCodeToHTTPStatus(err))
		json.NewEncoder(w).Encode(errorResponse(err))
	default:
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(program)
	}
}

func (c *ScheduleFrontendController) CreateProgramEpisode(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	// save the uploaded file
	// r.ParseMultipartForm(10 << 20)
	// file, handler, err := r.FormFile("myFile")
	// if err != nil {
	//     fmt.Println("Error Retrieving the File")
	//     fmt.Println(err)
	//     return
	// }
	// defer file.Close()
	//
	// tempFile, err := ioutil.TempFile("temp-images", "upload-*.png")
	// if err != nil {
	//     fmt.Println(err)
	// }
	// defer tempFile.Close()
	//
	// fileBytes, err := ioutil.ReadAll(file)
	// if err != nil {
	//     fmt.Println(err)
	// }
	// tempFile.Write(fileBytes)

	parameters := new(transport.PutProgramEpisodeRequest)
	if err := json.NewDecoder(r.Body).Decode(&parameters); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(errorResponse(err))
		return
	}

	parameters.ProgramId = chi.URLParam(r, "program_id")
	parameters.AirDate = time.Now().Format(transport.TimeFormat)
	parameters.RecordingPath = "get the path"

	program, err := c.client.PutProgramEpisode(r.Context(), parameters)
	switch {
	case err != nil:
		w.WriteHeader(grpcCodeToHTTPStatus(err))
		json.NewEncoder(w).Encode(errorResponse(err))
	default:
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(program)
	}
}

func (c *ScheduleFrontendController) CreateProgramScheduleEntry(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	parameters := new(transport.PutProgramScheduleEntryRequest)
	if err := json.NewDecoder(r.Body).Decode(&parameters); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(errorResponse(err))
		return
	}

	parameters.ProgramId = chi.URLParam(r, "program_id")

	program, err := c.client.PutProgramScheduleEntry(r.Context(), parameters)
	switch {
	case err != nil:
		w.WriteHeader(grpcCodeToHTTPStatus(err))
		json.NewEncoder(w).Encode(errorResponse(err))
	default:
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(program)
	}
}

type errorResponseBody struct {
	Error string `json:"error"`
}

func errorResponse(err error) errorResponseBody {
	return errorResponseBody{
		Error: err.Error(),
	}
}

func grpcCodeToHTTPStatus(err error) int {
	switch grpc.Code(err) {
	case codes.InvalidArgument:
		return http.StatusBadRequest
	case codes.OK:
		return http.StatusOK
	case codes.NotFound:
		return http.StatusNotFound
	case codes.AlreadyExists:
		return http.StatusConflict
	case codes.PermissionDenied:
		return http.StatusUnauthorized
	case codes.FailedPrecondition:
		return http.StatusPreconditionRequired
	case codes.OutOfRange:
		return http.StatusRequestedRangeNotSatisfiable
	case codes.Unimplemented:
		return http.StatusNotFound
	case codes.Internal:
		return http.StatusInternalServerError
	case codes.Unavailable:
		return http.StatusServiceUnavailable
	case codes.Unauthenticated:
		return http.StatusUnauthorized
	}

	return http.StatusInternalServerError
}
