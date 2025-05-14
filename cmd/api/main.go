package main

import (
	"context"
	"fmt"
	"insider-mert/internal/handlers/messages"
	"insider-mert/internal/handlers/service"
	"insider-mert/internal/helpers/telemetry"
	"insider-mert/internal/repository"
	"insider-mert/internal/workers"
	"insider-mert/internal/workers/message"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/jackc/pgx/v5"
	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
)

const (
	BasePath = "/api/v1"
)

//	@title			Sample Messaging Worker
//	@version		1.0
//	@description	This is a sample messaging worker.

//	@license.name	Apache 2.0
//	@license.url	http://www.apache.org/licenses/LICENSE-2.0.html

// @host		localhost:8080
// @BasePath	/api/v1
func main() {
	mux := http.NewServeMux()

	ctx := context.Background()
	otelShutdown, err := telemetry.SetupOTelSDK(ctx)
	if err != nil {
		panic(err)
	}
	defer otelShutdown(ctx)

	conn, err := pgx.Connect(context.Background(), os.Getenv("DATABASE_URL"))
	if err != nil {
		panic(err)
	}
	defer conn.Close(context.Background())

	// check if we can actually connect to it.
	var data any
	err = conn.QueryRow(ctx, "select version();").Scan(&data)
	if err != nil {
		panic(err)
	}

	worker, err := workers.Create(ctx, time.Second*3, message.SendMessage, conn)
	if err != nil {
		panic(err)
	}
	worker.Start()

	httpHandler, err := RegisterRoutes(mux, worker, conn)
	if err != nil {
		panic(err)
	}

	go func() {
		err = http.ListenAndServe(":8080", httpHandler)
		if err != nil {
			panic(err)
		}
	}()

	exitSignal := make(chan os.Signal, 1)
	signal.Notify(exitSignal, syscall.SIGTERM)
	<-exitSignal
	ctx.Done()
}

// this code based on https://opentelemetry.io/docs/languages/go/getting-started/
func RegisterRoutes(mux *http.ServeMux, worker *workers.Worker, conn *pgx.Conn) (http.Handler, error) {
	handleFunc := func(method string, pattern string, handlerFunc func(http.ResponseWriter, *http.Request)) {
		fullPath := fmt.Sprintf("%s %s%s", method, BasePath, pattern)
		handler := otelhttp.WithRouteTag(fullPath, http.HandlerFunc(handlerFunc))
		mux.Handle(fullPath, handler)
	}

	messagesRepository := repository.NewMessagesRepository(conn)

	serviceHandler, err := service.NewServiceHandler(worker)
	if err != nil {
		return nil, err
	}
	handleFunc("POST", "/service/toggle", serviceHandler.Toggle)

	messagesHandler, err := messages.NewMessagesHandler(messagesRepository)
	if err != nil {
		return nil, err
	}
	handleFunc("GET", "/messages/list", messagesHandler.List)

	handler := otelhttp.NewHandler(mux, "/")
	return handler, nil
}
