package controller

import (
	"fmt"
	"github.com/go-chi/chi/v5"
	middlewares "github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/httplog"
	"github.com/sethvargo/go-limiter/httplimit"
	"github.com/sethvargo/go-limiter/memorystore"
	logging "github.com/sirupsen/logrus"
	"gorm.io/gorm"
	http "net/http"
	"recordgram/botapi"
	"recordgram/config"
)

func CORS(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
		next.ServeHTTP(w, r)
	})
}
func setupServer(mux *chi.Mux, config config.Config, db *gorm.DB, bot *botapi.Bot) {
	mux.Use(CORS)
	mux.Use(middlewares.Recoverer)
	mux.Use(middlewares.Logger)
	setupRateLimit(mux, config)
	if config.Debug {
		mux.Use(httplog.Handler(httplog.NewLogger("records_server", httplog.Options{
			JSON:    true,
			Concise: true,
		})))
	}
	mux.Route("/records", RecordsRouter(db, bot))
}

func setupRateLimit(mux *chi.Mux, config config.Config) {
	limitConfig := config.Server.RateLimit
	if !limitConfig.Enabled {
		logging.Warn("RateLimit: disabled")
		return
	}
	store, err := memorystore.New(&memorystore.Config{
		Tokens:   limitConfig.Tokens,
		Interval: limitConfig.Interval,
	})
	if err != nil {
		logging.WithError(err).Panic("RateLimit: error creating memorystore")
	}
	middleware, err := httplimit.NewMiddleware(store, httplimit.IPKeyFunc())
	if err != nil {
		logging.Panic(err)
	}
	mux.Use(middleware.Handle)
	logging.WithFields(logging.Fields{"tokens": limitConfig.Tokens, "interval": limitConfig.Interval}).Infof("RateLimit initialized")
}

type Server struct {
	start func()
}

func NewServer(config config.Config, db *gorm.DB, bot *botapi.Bot) *Server {
	router := chi.NewMux()
	setupServer(router, config, db, bot)
	protocol := "http"
	if config.Server.Https.Enabled {
		protocol = "https"
	}
	logging.WithFields(logging.Fields{"protocol": protocol, "port": config.Server.Port}).Info(fmt.Sprintf("Server: initilalizing"))
	return &Server{
		func() {
			if config.Server.Https.Enabled {
				err := http.ListenAndServeTLS(config.Server.Host+":"+config.Server.Port, config.Server.Https.Cert, config.Server.Https.Key, router)
				if err != nil {
					logging.WithError(err).Panic("Server: error")
				}
			} else {
				err := http.ListenAndServe(config.Server.Host+":"+config.Server.Port, router)
				if err != nil {
					logging.WithError(err).Panic("Server: er ror")
				}
			}
		},
	}
}

func (server Server) Start() {
	server.start()
}
