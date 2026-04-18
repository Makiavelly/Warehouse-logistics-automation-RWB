package main

import (
	"context"
	"flag"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	JWTAdapter "logistics-service/logistics-service/adapters/JWT"
	dbAdapter "logistics-service/logistics-service/adapters/db"
	httpServer "logistics-service/logistics-service/adapters/http_server"
	"logistics-service/logistics-service/adapters/logger"
	middlewareAdapter "logistics-service/logistics-service/adapters/middleware"
	mlClientHTTP "logistics-service/logistics-service/adapters/ml_client"
	mlClientGRPC "logistics-service/logistics-service/adapters/ml_client_grpc"
	serveMux "logistics-service/logistics-service/adapters/mux"
	"logistics-service/logistics-service/adapters/rest"
	validatorAdapter "logistics-service/logistics-service/adapters/validator"
	"logistics-service/logistics-service/config"
	"logistics-service/logistics-service/core/models"
	"logistics-service/logistics-service/core/ports"
	coreService "logistics-service/logistics-service/core/service"
)

const defaultConfigPath = "logistics-service/config.yml"

func main() {
	var configPath string
	flag.StringVar(&configPath, "config", defaultConfigPath, "path to config file")
	flag.Parse()

	cfg := config.MustLoadConfig(configPath)
	log := logger.New(cfg.App.LogLevel)

	log.Info("logistics-service starting...")

	if err := run(log, cfg); err != nil {
		log.Error("application error", "error", err)
		os.Exit(1)
	}
}

func run(log *logger.Logger, cfg *config.Config) error {
	db, err := dbAdapter.New(log, cfg.DB.Address)
	if err != nil {
		return err
	}
	if err = db.Migrate(); err != nil {
		return err
	}

	jwt := JWTAdapter.NewJWT(cfg.Auth)

	var ml ports.MLClient
	if cfg.ML.UseGRPC {
		grpcML, err := mlClientGRPC.New(cfg.ML.Address)
		if err != nil {
			return err
		}
		ml = grpcML
		log.Info("ML client: gRPC", "address", cfg.ML.Address)
	} else {
		ml = mlClientHTTP.New(cfg.ML.Address)
		log.Info("ML client: HTTP", "address", cfg.ML.Address)
	}

	validator := validatorAdapter.NewValidator()
	middleware := middlewareAdapter.NewMiddleware(log, cfg.Auth.APIKey)

	svc := coreService.NewService(log, jwt, db, ml)

	mux := serveMux.NewMux()
	registerRoutes(log, mux, middleware, jwt, svc, validator)

	loggedMux := middleware.Logging(middleware.CORS(mux))

	serv := httpServer.New(cfg.Server.Address, loggedMux, cfg.Server.ReadTimeout, cfg.Server.WriteTimeout)
	log.Info("server listening", "address", cfg.Server.Address)

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	go func() {
		<-ctx.Done()
		log.Info("shutdown signal received")
		shutCtx, cancel := context.WithTimeout(context.Background(), cfg.Server.MaxShutdownTimeout)
		defer cancel()
		if err := serv.Stop(shutCtx); err != nil {
			log.Error("shutdown error", "error", err)
		}
	}()

	return serv.Start()
}

func registerRoutes(
	log ports.Logger,
	mux ports.ServeMux,
	md *middlewareAdapter.Middleware,
	jwt ports.TokenGenerator,
	svc ports.Service,
	v ports.Validator,
) {
	adminAuth := md.Auth(jwt, models.RoleAdmin)
	driverAuth := md.Auth(jwt, models.RoleDriver)
	anyAuth := md.Auth(jwt, "")

	// Health check
	mux.HandleFunc("GET /health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"status":"ok"}`))
	})

	// Auth (public)
	mux.HandleFunc("POST /api/auth/login", rest.NewLoginHandler(log, svc, v))
	mux.HandleFunc("POST /api/auth/register", adminAuth(rest.NewRegisterHandler(log, svc, v)))

	// Admin — Warehouses
	mux.HandleFunc("POST /api/warehouses", adminAuth(rest.NewCreateWarehouseHandler(log, svc, v)))
	mux.HandleFunc("GET /api/warehouses", anyAuth(rest.NewGetWarehousesHandler(log, svc)))
	mux.HandleFunc("DELETE /api/warehouses/{warehouse_id}", adminAuth(rest.NewDeleteWarehouseHandler(log, svc)))

	// Admin — Routes
	mux.HandleFunc("POST /api/warehouses/{warehouse_id}/routes", adminAuth(rest.NewCreateRouteHandler(log, svc, v)))
	mux.HandleFunc("GET /api/warehouses/{warehouse_id}/routes", anyAuth(rest.NewGetRoutesHandler(log, svc)))
	mux.HandleFunc("DELETE /api/warehouses/{warehouse_id}/routes/{route_id}", adminAuth(rest.NewDeleteRouteHandler(log, svc)))

	// Admin — Thresholds
	mux.HandleFunc("PUT /api/thresholds", adminAuth(rest.NewSetThresholdHandler(log, svc, v)))
	mux.HandleFunc("GET /api/thresholds", anyAuth(rest.NewGetThresholdsHandler(log, svc)))

	// Admin — Forecasts
	mux.HandleFunc("POST /api/forecasts/predict", adminAuth(rest.NewRequestForecastHandler(log, svc, v)))
	mux.HandleFunc("GET /api/forecasts", anyAuth(rest.NewGetForecastsHandler(log, svc)))

	// Admin — Truck calls & analytics
	mux.HandleFunc("GET /api/truck-calls", anyAuth(rest.NewGetTruckCallsHandler(log, svc)))
	mux.HandleFunc("GET /api/truck-calls/accuracy", anyAuth(rest.NewGetTruckCallAccuracyHandler(log, svc)))

	// Admin — Drivers management
	mux.HandleFunc("GET /api/drivers", adminAuth(rest.NewGetDriversHandler(log, svc)))
	mux.HandleFunc("PUT /api/drivers/assign", adminAuth(rest.NewAssignDriverHandler(log, svc, v)))

	// Driver endpoints
	mux.HandleFunc("GET /api/driver/signal", driverAuth(rest.NewGetDriverSignalHandler(log, svc)))
	mux.HandleFunc("POST /api/driver/truck-calls/{truck_call_id}/timeliness", driverAuth(rest.NewReportTimelinessHandler(log, svc, v)))
	mux.HandleFunc("GET /api/driver/stats", driverAuth(rest.NewGetDriverStatsHandler(log, svc)))

	// External integration (API key auth)
	mux.HandleFunc("POST /api/data/ingest", md.APIKeyAuth(rest.NewIngestDataHandler(log, svc, v)))
	mux.HandleFunc("POST /api/model/retrain", adminAuth(rest.NewRetrainModelHandler(log, svc, v)))
}