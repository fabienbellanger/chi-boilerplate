package chi_router

import (
	"log"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"github.com/spf13/viper"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func (s *ChiServer) initMiddlewares(r *chi.Mux, log *zap.Logger) {
	r.Use(middleware.RequestID)
	if viper.GetBool("ENABLE_ACCESS_LOG") {
		r.Use(initLogger(log))
	}
	r.Use(middleware.Recoverer)

	// Profiler
	if viper.GetBool("SERVER_PPROF") {
		r.Group(func(r chi.Router) {
			r.Use(initBasicAuth())

			r.Mount("/debug", middleware.Profiler())
		})
	}
}

func initCORS() func(next http.Handler) http.Handler {
	return cors.Handler(cors.Options{
		AllowedOrigins:   viper.GetStringSlice("CORS_ALLOWED_ORIGINS"),
		AllowedMethods:   viper.GetStringSlice("CORS_ALLOWED_METHODS"),
		AllowedHeaders:   viper.GetStringSlice("CORS_ALLOWED_HEADERS"),
		ExposedHeaders:   viper.GetStringSlice("CORS_EXPOSED_HEADERS"),
		AllowCredentials: viper.GetBool("CORS_ALLOW_CREDENTIALS"),
		MaxAge:           viper.GetInt("CORS_MAX_AGE"),
	})
}

func initLogger(log *zap.Logger) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		fn := func(w http.ResponseWriter, r *http.Request) {
			start := time.Now().UTC()

			ww := middleware.NewWrapResponseWriter(w, r.ProtoMajor)

			next.ServeHTTP(ww, r)

			stop := time.Since(start)
			url := r.Host + r.RequestURI // TODO: Do better, missing https/http
			fields := []zapcore.Field{
				zap.Int("code", ww.Status()),
				zap.String("method", r.Method),
				zap.String("path", r.URL.Path),
				zap.String("url", url),
				zap.String("ip", r.RemoteAddr), // TODO: Remove port
				zap.String("userAgent", r.UserAgent()),
				zap.String("latency", stop.String()),
				// zap.String("requestId", fmt.Sprintf("%s", c.Locals("requestid"))),
			}

			log.Info("", fields...)
		}
		return http.HandlerFunc(fn)
	}
}

func initBasicAuth() func(next http.Handler) http.Handler {
	creds := make(map[string]string, 1)
	creds[viper.GetString("SERVER_BASICAUTH_USERNAME")] = viper.GetString("SERVER_BASICAUTH_PASSWORD")

	log.Printf("%v\n", creds)

	return middleware.BasicAuth("Restricted", creds)
}
