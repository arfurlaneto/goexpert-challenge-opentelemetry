package main

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/arfurlaneto/goexpert-challenge-opentelemetry/service-a/internal/handlers"
	"github.com/arfurlaneto/goexpert-challenge-opentelemetry/service-a/internal/services"
	"github.com/arfurlaneto/goexpert-challenge-opentelemetry/service-a/internal/usecases"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/chi/v5"
	"github.com/joho/godotenv"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.17.0"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	godotenv.Load(".env")

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
	defer stop()

	otelShutdown, err := setupOTelSDK(ctx)
	if err != nil {
		panic(err)
	}

	defer func() {
		err = errors.Join(err, otelShutdown(context.Background()))
	}()

	r := chi.NewRouter()
	r.Use(middleware.Recoverer)

	getTemperatureByCepHandler := handlers.NewGetTemperatureByCepHandler(
		usecases.NewGetTemperatureByCepUseCase(
			services.NewServiceBService(),
		),
	)

	r.Get("/healthz", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	})

	r.Post("/", getTemperatureByCepHandler.Handle)

	server := &http.Server{Addr: ":8080", Handler: r}
	srvErr := make(chan error, 1)
	go func() {
		srvErr <- server.ListenAndServe()
	}()

	select {
	case err = <-srvErr:
		panic(err)
	case <-ctx.Done():
		stop()
	}

	err = server.Shutdown(context.Background())
}

func setupOTelSDK(ctx context.Context) (shutdown func(context.Context) error, err error) {
	res, err := resource.New(ctx, resource.WithAttributes(semconv.ServiceName(os.Getenv("OTEL_SERVICE_NAME"))))
	if err != nil {
		return nil, fmt.Errorf("failed to create resource: %w", err)
	}

	ctx, cancel := context.WithTimeout(ctx, time.Second)
	defer cancel()

	conn, err := grpc.DialContext(
		ctx,
		os.Getenv("OTEL_COLLECTOR_URL"),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithBlock(),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create gRPC connection to collector: %w", err)
	}

	tracerExporter, err := otlptracegrpc.New(ctx, otlptracegrpc.WithGRPCConn(conn))
	if err != nil {
		return nil, fmt.Errorf("failed to create trace exporter: %w", err)
	}

	bsp := sdktrace.NewBatchSpanProcessor(tracerExporter)
	traceProvider := sdktrace.NewTracerProvider(
		sdktrace.WithSampler(sdktrace.AlwaysSample()),
		sdktrace.WithResource(res),
		sdktrace.WithSpanProcessor(bsp),
	)
	otel.SetTracerProvider(traceProvider)

	otel.SetTextMapPropagator(propagation.TraceContext{})

	return traceProvider.Shutdown, nil
}
