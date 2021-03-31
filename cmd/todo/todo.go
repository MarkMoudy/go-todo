package main

import (
	"context"
	"fmt"
	stdlog "log"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"github.com/MarkMoudy/go-todo/todo"
	"github.com/go-kit/kit/log"
	"github.com/go-stack/stack"
	"github.com/gorilla/mux"
)

func main() {
	baseLogger := log.NewLogfmtLogger(log.NewSyncWriter(os.Stdout))
	srvLogger := newLogContext(baseLogger, "server")
	apiLogger := newLogContext(baseLogger, "api")

	// ensure standard library logger writes to go-kit log, including caller
	// key.
	stdlog.SetFlags(stdlog.Llongfile)
	stdlog.SetOutput(log.NewStdlibAdapter(srvLogger, log.FileKey("caller")))

	// Services
	todoSvc := todo.NewInmemService()
	todoSvc = todo.NewLoggingService(todoSvc, apiLogger)

	// Router
	r := mux.NewRouter()
	apiRouter := r.PathPrefix("/api").Subrouter()
	todo.MakeRoutes(apiRouter, todoSvc, apiLogger)

	// Server
	httpServer := http.Server{
		// TODO: make this configurable via a CLI flag
		Addr:    ":6060",
		Handler: r,
	}

	errC := make(chan error)

	go func() {
		defer httpServer.Shutdown(context.Background())
		c := make(chan os.Signal, 1)
		signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)
		errC <- fmt.Errorf("%s", <-c)

	}()

	go func() {
		srvLogger.Log("event", "starting HTTP Server")
		errC <- httpServer.ListenAndServe()
	}()

	srvLogger.Log("event", "shutting down", "err", <-errC)
}

func newLogContext(logger log.Logger, app string) log.Logger {
	return log.With(logger,
		"time", log.DefaultTimestampUTC,
		"app", app,
		"caller", log.Valuer(func() interface{} {
			return pkgCaller{c: stack.Caller(3)}
		}))
}

type pkgCaller struct {
	c stack.Call
}

func (p pkgCaller) String() string {
	caller := fmt.Sprintf("%+v", p.c)
	caller = strings.TrimPrefix(caller, "github.com/MarkMoudy/go-todo/")
	return caller
}
