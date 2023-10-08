package api

import (
	"context"
	"fmt"
	"net/http"

	"go.uber.org/zap"

	"github.com/asfourco/templates/backend/db"
	"github.com/gorilla/mux"
	"github.com/streamingfast/dhttp"
	"github.com/streamingfast/shutter"
)

type Server struct {
	ctx        context.Context
	listenPort uint16

	router  *mux.Router
	handler http.Handler

	postgresClient *db.PostgresClient

	DefaultPageSize uint16
}

func NewServer(ctx context.Context, listenPort uint16, postgresClient *db.PostgresClient) (*Server, error) {
	s := &Server{
		ctx:             ctx,
		listenPort:      listenPort,
		postgresClient:  postgresClient,
		DefaultPageSize: db.DEFAULT_PAGE_SIZE,
	}
	if err := s.configureHttpServer(); err != nil {
		return nil, err
	}
	return s, nil
}

func (s *Server) Start(app *shutter.Shutter) error {
	stdLogger, err := zap.NewStdLogAt(zlog, zap.InfoLevel)
	if err != nil {
		return fmt.Errorf("unable to create error logger: %w", err)
	}

	server := &http.Server{
		Addr:     fmt.Sprintf(":%d", s.listenPort),
		Handler:  s.handler,
		ErrorLog: stdLogger,
	}

	go func() {
		zlog.Info("serving HTTP", zap.Int("port", int(s.listenPort)))
		go app.Shutdown(server.ListenAndServe())
	}()

	return nil
}

func (s *Server) configureHttpServer() error {
	zlog.Info("configuring HTTP server")

	// monitoring
	monitoringRouter := s.router.PathPrefix("/").Subrouter()
	monitoringRouter.Path("/healthz").Handler(dhttp.JSONHandler(getHealth))

	// API v1 router
	apiV1Router := s.router.PathPrefix("/api/v1").Subrouter()
	apiV1Router.Use(dhttp.NewCORSMiddleware("*"))

	// API REST router
	apiRestRouter := apiV1Router.PathPrefix("/").Subrouter()

	// API Todo router
	apiTodoRouter := apiRestRouter.PathPrefix("/todo").Subrouter()
	apiTodoRouter.Methods("GET", "OPTIONS").Path("/{id}").Handler(dhttp.JSONHandler(s.GetTodo))
	apiTodoRouter.Methods("GET", "OPTIONS").Path("/").Handler(dhttp.JSONHandler(s.GetTodoList))
	apiTodoRouter.Methods("POST", "OPTIONS").Path("/").Handler(dhttp.JSONHandler(s.CreateTodo))
	apiTodoRouter.Methods("PUT", "OPTIONS").Path("/{id}").Handler(dhttp.JSONHandler(s.UpdateTodo))
	apiTodoRouter.Methods("DELETE", "OPTIONS").Path("/{id}").Handler(dhttp.JSONHandler(s.DeleteTodo))

	// walk configured routes
	err := s.router.Walk(func(r *mux.Route, _ *mux.Router, _ []*mux.Route) error {
		if (*route)(r).String() != "" {
			zlog.Info("routes", zap.Stringer("route", (*route)(r)))
		}
		return nil
	})

	if err != nil {
		return fmt.Errorf("unable to walk routes: %w", err)
	}

	zlog.Info("HTTP server configured")
	s.handler = s.router
	return nil
}
