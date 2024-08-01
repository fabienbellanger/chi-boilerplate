package chi_router

import (
	"chi_boilerplate/pkg/infrastructure/chi_router/handlers"
	"chi_boilerplate/pkg/infrastructure/logger"
	"chi_boilerplate/utils"
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/lestrrat-go/jwx/v2/jwt"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"github.com/go-chi/jwtauth/v5"
	"github.com/spf13/viper"
)

var tokenAuth *jwtauth.JWTAuth

func (s *ChiServer) initJWTToken() error {
	algo := viper.GetString("JWT_ALGO")
	key, err := utils.GetKeyFromAlgo(algo, viper.GetString("JWT_SECRET"), viper.GetString("JWT_PUBLIC_KEY_PATH"))
	if err != nil {
		return err
	}

	tokenAuth = jwtauth.New(algo, key, nil)

	return nil
}

func (s *ChiServer) initMiddlewares(r *chi.Mux) {
	r.Use(s.requestID) // Must be before the access logger
	if viper.GetBool("ENABLE_ACCESS_LOG") {
		r.Use(s.initAccessLogger())
	}
	r.Use(middleware.Recoverer)

	// Profiler
	if viper.GetBool("SERVER_PPROF") {
		r.Group(func(r chi.Router) {
			r.Use(s.initBasicAuth())

			r.Mount("/debug", middleware.Profiler())
		})
	}
}

func (s *ChiServer) initAccessLogger() func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		fn := func(w http.ResponseWriter, r *http.Request) {
			start := time.Now().UTC()
			ctxID := r.Context().Value(handlers.RequestIDKey("request_id"))
			var requestId string
			if ctxID != nil {
				requestId = fmt.Sprintf("%s", r.Context().Value(handlers.RequestIDKey("request_id")))
			}
			ww := middleware.NewWrapResponseWriter(w, r.ProtoMajor)

			next.ServeHTTP(ww, r)

			stop := time.Since(start)
			url := r.Host + r.RequestURI // TODO: Do better, missing https:// or http://
			fields := logger.Fields{
				logger.NewField("code", "int", ww.Status()),
				logger.NewField("method", "string", r.Method),
				logger.NewField("path", "string", r.URL.Path),
				logger.NewField("url", "string", url),
				logger.NewField("ip", "string", r.RemoteAddr), // TODO: Remove port
				logger.NewField("userAgent", "string", r.UserAgent()),
				logger.NewField("latency", "string", stop.String()),
				logger.NewField("request_id", "string", requestId),
			}

			s.Logger.Info("", fields)
		}
		return http.HandlerFunc(fn)
	}
}

func (s *ChiServer) initCORS() func(next http.Handler) http.Handler {
	return cors.Handler(cors.Options{
		AllowedOrigins:   viper.GetStringSlice("CORS_ALLOWED_ORIGINS"),
		AllowedMethods:   viper.GetStringSlice("CORS_ALLOWED_METHODS"),
		AllowedHeaders:   viper.GetStringSlice("CORS_ALLOWED_HEADERS"),
		ExposedHeaders:   viper.GetStringSlice("CORS_EXPOSED_HEADERS"),
		AllowCredentials: viper.GetBool("CORS_ALLOW_CREDENTIALS"),
		MaxAge:           viper.GetInt("CORS_MAX_AGE"),
	})
}

func (s *ChiServer) initBasicAuth() func(next http.Handler) http.Handler {
	creds := make(map[string]string, 1)
	creds[viper.GetString("SERVER_BASICAUTH_USERNAME")] = viper.GetString("SERVER_BASICAUTH_PASSWORD")

	return middleware.BasicAuth("Restricted", creds)
}

func (s *ChiServer) initJWT(r chi.Router) {
	r.Use(jwtauth.Verifier(tokenAuth))
	r.Use(s.jwtAuthenticator(tokenAuth))
}

func (s *ChiServer) jwtAuthenticator(ja *jwtauth.JWTAuth) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		hfn := func(w http.ResponseWriter, r *http.Request) {
			token, _, err := jwtauth.FromContext(r.Context())

			if err != nil {
				utils.Err401(w, err, "Unauthorized", nil) // TODO: Error not managed
				return
			}

			if token == nil || jwt.Validate(token, ja.ValidateOptions()...) != nil {
				utils.Err401(w, nil, "Unauthorized", nil) // TODO: Error not managed
				return
			}

			// Token is authenticated, pass it through
			next.ServeHTTP(w, r)
		}
		return http.HandlerFunc(hfn)
	}
}

func (s *ChiServer) requestID(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		id := uuid.New().String()
		ctx := context.WithValue(r.Context(), handlers.RequestIDKey("request_id"), id)

		w.Header().Add("X-Request-Id", id)

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
