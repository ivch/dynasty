package auth

import (
	"context"
	"encoding/json"
	"net/http"
	"strconv"
	"strings"

	"github.com/go-chi/chi"
	httptransport "github.com/go-kit/kit/transport/http"
	"github.com/jinzhu/gorm"
	"github.com/rs/zerolog"
	"gopkg.in/go-playground/validator.v9"
)

func New(log *zerolog.Logger, db *gorm.DB, usrv UserService, jwtSecret string) http.Handler {
	svc := newService(log, db, usrv, jwtSecret)
	options := []httptransport.ServerOption{
		httptransport.ServerErrorEncoder(encodeHTTPError),
	}

	r := chi.NewRouter()

	r.Method("POST", "/v1/login", httptransport.NewServer(
		makeLoginEndpoint(svc),
		decodeLoginRequest(log),
		encodeHTTPResponse,
		options...))

	r.Method("POST", "/v1/refresh", httptransport.NewServer(
		makeRefreshEndpoint(svc),
		decodeRefreshTokenRequest(log),
		encodeHTTPResponse,
		options...))

	r.Method("GET", "/v1/gwfa", authCheck(log, svc))

	return r
}

func decodeRefreshTokenRequest(log *zerolog.Logger) httptransport.DecodeRequestFunc {
	return func(_ context.Context, r *http.Request) (interface{}, error) {
		var req refreshTokenRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			log.Error().Err(err).Msg("failed to decode request")
			return nil, err
		}

		if err := validator.New().Struct(&req); err != nil {
			log.Error().Err(err).Msg("error validating request")
			return nil, err
		}

		return &req, nil
	}
}

func decodeLoginRequest(log *zerolog.Logger) httptransport.DecodeRequestFunc {
	return func(_ context.Context, r *http.Request) (interface{}, error) {
		var req loginRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			log.Error().Err(err).Msg("failed to decode request")
			return nil, err
		}

		// ip, _, err := net.SplitHostPort(r.RemoteAddr)
		// if err != nil {
		//	log.Error().Err(err).Msg("failed to parse host")
		//	return nil, fmt.Errorf("userip: %q is not IP:port", r.RemoteAddr)
		// }
		//
		// userIP := net.ParseIP(ip)
		// if userIP == nil {
		//	log.Error().Err(err).Msg("failed to parse ip")
		//	return nil, fmt.Errorf("userip: %q is not IP:port", r.RemoteAddr)
		// }
		//
		// req.IP = userIP
		// req.Ua = r.Header.Get("User-Agent")

		if err := validator.New().Struct(&req); err != nil {
			log.Error().Err(err).Msg("error validating request")
			return nil, err
		}

		return &req, nil
	}
}

func authCheck(log *zerolog.Logger, srv Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		token := r.Header.Get("Authorization")
		if token == "" {
			log.Warn().Msg("gateway forward auth: there is no Authorization header in request")
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		id, err := srv.Gwfa(strings.TrimPrefix(token, "Bearer "))
		if err != nil {
			log.Warn().Msgf("gateway forward auth: %v", err)
			w.WriteHeader(http.StatusForbidden)
			return
		}
		w.Header().Set("X-Auth-User", strconv.FormatUint(uint64(id), 10))
		w.WriteHeader(http.StatusOK)
	}
}

func encodeHTTPResponse(_ context.Context, w http.ResponseWriter, response interface{}) error {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusOK)

	return json.NewEncoder(w).Encode(response)
}

func encodeHTTPError(_ context.Context, err error, w http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusInternalServerError)
	json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
}