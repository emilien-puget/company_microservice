package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"github.com/caarlos0/env/v6"
	"github.com/labstack/echo-contrib/prometheus"
	"github.com/labstack/echo/v4"
	"github.com/prometheus/client_golang/prometheus/promhttp"

	"github.com/emilien-puget/company_microservice/configuration"
)

var ErrStopSignalReceived = errors.New("stop signal received")

// Version set via ldflags
var Version = "local"

const service = "company_microservice"

func main() {
	eCfg := configuration.Api{}
	if err := env.Parse(&eCfg, (env.Options{RequiredIfNoDef: true})); err != nil {
		log.Printf("%+v\n", err)
		os.Exit(-1)
	}

	ctx, cl := Init()
	defer cl(nil)

	// validate := validator.New()

	e := echo.New()
	defer e.Shutdown(context.Background())
	p := prometheus.NewPrometheus(service, nil)
	p.Use(e)

	go func() {
		err := e.Start(fmt.Sprintf(":%s", eCfg.Port))
		if err != nil {
			cl(fmt.Errorf("exposed server:%w", err))
		}
	}()

	srv := initInternalSrv(eCfg.InternalPort)
	defer srv.Shutdown(context.Background())
	go func() {
		err := srv.ListenAndServe()
		if err != nil {
			cl(fmt.Errorf("internal server:%w", err))
		}
	}()
	log.Printf("starting %s", Version)
	<-ctx.Done()
}

func initInternalSrv(internalPort string) *http.Server {
	mux := http.NewServeMux()
	mux.HandleFunc("/ping", func(writer http.ResponseWriter, request *http.Request) {
		writer.WriteHeader(http.StatusOK)
	})
	mux.Handle("/metrics", promhttp.Handler())
	srv := &http.Server{
		Addr:    fmt.Sprintf(":%s", internalPort),
		Handler: mux,
	}
	return srv
}

func Init() (ctx context.Context, cl func(err error)) {
	ctx = context.Background()
	stopSignal := notifyStopSignal()
	ctx, cancelFunc := context.WithCancelCause(ctx)
	go func() {
		<-stopSignal
		cancelFunc(ErrStopSignalReceived)
	}()

	var once sync.Once
	return ctx, func(err error) {
		once.Do(func() {
			cancelFunc(err)
			loggingExit(ctx)
		})
	}
}

func notifyStopSignal() <-chan os.Signal {
	gracefulStop := make(chan os.Signal, 1)
	signal.Notify(gracefulStop, syscall.SIGTERM, syscall.SIGINT, syscall.SIGQUIT)
	return gracefulStop
}

func loggingExit(ctx context.Context) {
	err := context.Cause(ctx)
	if err != nil {
		if errors.Is(err, ErrStopSignalReceived) {
			log.Print("stop signal received")
			return
		}
		if errors.Is(err, context.Canceled) {
			log.Print("context cancel without cause")
			return
		}
		log.Print(err)
		return
	}
}
