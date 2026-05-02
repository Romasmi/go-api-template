package http

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"path/filepath"

	api "github.com/Romasmi/s-shop-microservices/internal/api"
	"github.com/Romasmi/s-shop-microservices/internal/interface/http/middleware"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	httpSwagger "github.com/swaggo/http-swagger/v2"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func NewGatewayServer(grpcAddr string, httpPort uint) (*http.Server, error) {
	ctx := context.Background()
	mux := runtime.NewServeMux()
	opts := []grpc.DialOption{grpc.WithTransportCredentials(insecure.NewCredentials())}

	err := api.RegisterUserServiceHandlerFromEndpoint(ctx, mux, grpcAddr, opts)
	if err != nil {
		return nil, err
	}

	mainMux := http.NewServeMux()
	mainMux.Handle("/", mux)
	mainMux.Handle("/swagger/", httpSwagger.Handler(
		httpSwagger.URL("/swagger-static/api/user.swagger.json"),
	))
	mainMux.HandleFunc("/swagger-static/", serveSwaggerStatic)
	mainMux.HandleFunc("/proto/", serveProto)
	mainMux.Handle("/metrics", promhttp.Handler())

	return &http.Server{
		Addr:    fmt.Sprintf(":%d", httpPort),
		Handler: middleware.MetricsMiddleware(mainMux),
	}, nil
}

func serveSwaggerStatic(w http.ResponseWriter, r *http.Request) {
	path := filepath.Clean(r.URL.Path)
	prefix := "/swagger-static/"
	http.ServeFile(w, r, filepath.Join("api/swagger", path[len(prefix):]))
}

func serveProto(w http.ResponseWriter, r *http.Request) {
	path := filepath.Clean(r.URL.Path)
	if path == "/proto/" || path == "/proto" {
		http.ServeFile(w, r, "api/user.proto")
		return
	}

	trimmedPath := path[len("/proto/"):]

	// Try api/ folder
	apiPath := filepath.Join("api", trimmedPath)
	if _, err := os.Stat(apiPath); err == nil {
		http.ServeFile(w, r, apiPath)
		return
	}

	// Try google/ folder
	googlePath := filepath.Join("google", trimmedPath)
	if _, err := os.Stat(googlePath); err == nil {
		http.ServeFile(w, r, googlePath)
		return
	}

	http.NotFound(w, r)
}
