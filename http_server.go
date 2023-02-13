package tmpmail

import (
	"context"
	"crypto/rand"
	"crypto/tls"
	"embed"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/fs"
	"mime"
	"net/http"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"github.com/julienschmidt/httprouter"
	"github.com/jxskiss/base62"
	"github.com/rs/cors"
	"go.uber.org/zap"
	"golang.org/x/time/rate"

	"tmpmail/entity"
)

//go:embed all:ui/dist/*
var uiFS embed.FS

type HTTPServerStorage interface {
	CreateAccount(token, username string, ttl time.Duration) error
	ProlongAccount(token string, ttl time.Duration) error
	Account(token string) (entity.Account, error)
	RemoveAccount(token string) error
}

type HTTPServer struct {
	server            *http.Server
	storage           HTTPServerStorage
	authToken         string
	authRateLimiter   *ipRateLimiter
	defaultAccountTTL time.Duration
	ui                fs.FS
	api               http.Handler
	logger            *zap.Logger
}

func NewHTTPServer(l *zap.Logger, s HTTPServerStorage, addr string, tc *tls.Config,
	authToken string, defaultAccountTTL time.Duration) *HTTPServer {

	ui, _ := fs.Sub(uiFS, "ui/dist")

	srv := &HTTPServer{
		storage:           s,
		authToken:         authToken,
		authRateLimiter:   newIPRateLimiter(rate.Every(time.Hour), 5),
		defaultAccountTTL: defaultAccountTTL,
		ui:                ui,
		logger:            l,
	}

	api := httprouter.New()
	api.GET("/api/account", srv.getAPIAccount)
	api.POST("/api/account", srv.postAPIAccount)
	api.PUT("/api/account", srv.putAPIAccount)
	api.PATCH("/api/account", srv.patchAPIAccount)
	api.DELETE("/api/account", srv.deleteAPIAccount)

	corsHandler := cors.New(cors.Options{
		AllowedOrigins: []string{"https://tmp-mail.ru", "http://localhost:3000"},
		AllowedMethods: []string{
			http.MethodHead,
			http.MethodGet,
			http.MethodPost,
			http.MethodPut,
			http.MethodPatch,
			http.MethodDelete,
		},
		AllowedHeaders:   []string{"*"},
		AllowCredentials: true,
	})

	srv.api = corsHandler.Handler(api)

	r := httprouter.New()
	r.OPTIONS("/*path", srv.handler)
	r.GET("/*path", srv.handler)
	r.POST("/*path", srv.handler)
	r.PUT("/*path", srv.handler)
	r.PATCH("/*path", srv.handler)
	r.DELETE("/*path", srv.handler)

	srv.server = &http.Server{
		Addr:      addr,
		Handler:   r,
		TLSConfig: tc,
	}
	return srv
}

func (s *HTTPServer) ListenAndServe() error {
	err := s.server.ListenAndServeTLS("", "")
	if err != nil {
		return fmt.Errorf("listen and server: %w", err)
	}
	return nil
}

func (s *HTTPServer) Shutdown(ctx context.Context) error {
	return s.server.Shutdown(ctx)
}

const (
	authHeader   = "Authorization"
	acceptHeader = "Accept"
	tokenHeader  = "X-TOKEN"
	tokenLength  = 128
	emailsParam  = "emails"
)

func generateRandomString(length int) string {
	b := make([]byte, length)
	if _, err := rand.Read(b); err != nil {
		return ""
	}
	return base62.EncodeToString(b)
}

func generateEmail() string {
	const abc = "abcdefghijklmnopqrstuwxyz0123456789"
	var (
		l = len(abc)
		b strings.Builder
	)
	for i := int(time.Now().UnixNano()); i > 0; i = i / l {
		b.WriteByte(abc[i%l])
	}
	b.WriteString(strings.ToLower(generateRandomString(2)))
	return b.String()
}

func (s *HTTPServer) getAPIAccount(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	token := r.Header.Get(tokenHeader)

	a, err := s.storage.Account(token)
	if err != nil {
		if errors.Is(err, entity.ErrAccountDoesntExists) {
			w.WriteHeader(http.StatusNotFound)
			return
		}
		s.logger.Error("get account from storage", zap.Error(err))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(a)
}

func (s *HTTPServer) postAPIAccount(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	email := generateEmail()
	token := generateRandomString(tokenLength)
	err := s.storage.CreateAccount(token, email, s.defaultAccountTTL)
	if err != nil {
		s.logger.Warn("create email in storage", zap.Error(err))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(token)
}

func (s *HTTPServer) putAPIAccount(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	limiter := s.authRateLimiter.GetLimiter(r.RemoteAddr)

	reserve := limiter.Reserve()
	if !reserve.OK() {
		w.WriteHeader(http.StatusTooManyRequests)
		return
	}

	authToken := r.Header.Get(authHeader)
	if s.authToken != authToken {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	reserve.Cancel()

	emails := strings.Split(r.FormValue(emailsParam), ",")
	if len(emails) == 0 {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	for _, email := range emails {
		if len(email) == 0 {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
	}

	var (
		ttlStr = r.FormValue("ttl")
		ttl    = s.defaultAccountTTL
		err    error
	)
	if ttlStr != "" {
		ttl, err = time.ParseDuration(ttlStr)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
	}

	var tokens []string

	for _, email := range emails {
		token := generateRandomString(tokenLength)
		err = s.storage.CreateAccount(token, email, ttl)
		if err != nil {
			s.logger.Warn("create email in storage", zap.Error(err))
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		tokens = append(tokens, token)
	}

	accept := r.Header.Get(acceptHeader)

	switch accept {
	case "", "*/*", "plain/text":
		var b strings.Builder
		for i, email := range emails {
			token := tokens[i]
			b.WriteString(fmt.Sprintf("%s: %s\n", email, token))
		}
		w.Write([]byte(b.String()))
	case "application/json":
		json.NewEncoder(w).Encode(tokens)
	}
}

func (s *HTTPServer) patchAPIAccount(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	token := r.Header.Get(tokenHeader)
	err := s.storage.ProlongAccount(token, s.defaultAccountTTL)
	if err != nil {
		s.logger.Error("prolong account in storage", zap.Error(err))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

func (s *HTTPServer) deleteAPIAccount(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	s.storage.RemoveAccount(r.Header.Get(tokenHeader))
}

const indexHTMLFile = "index.html"

func (s *HTTPServer) serveFile(path string, w http.ResponseWriter) {
	file, err := s.ui.Open(path)
	if err != nil {
		file, _ = s.ui.Open(indexHTMLFile)
	}
	defer file.Close()

	w.Header().Set("Content-Type", mime.TypeByExtension(filepath.Ext(path)))
	io.Copy(w, file)
}

func (s *HTTPServer) handler(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	if strings.HasPrefix(r.URL.Path, "/api") {
		s.api.ServeHTTP(w, r)
		return
	}
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		w.Write([]byte("Method Not Allowed"))
		return
	}
	if r.URL.Path == "/" {
		s.serveFile(indexHTMLFile, w)
		return
	}
	s.serveFile(strings.TrimPrefix(r.URL.Path, "/"), w)
}

type ipRateLimiter struct {
	ips map[string]*rate.Limiter
	mu  *sync.RWMutex
	r   rate.Limit
	b   int
}

func newIPRateLimiter(r rate.Limit, b int) *ipRateLimiter {
	i := &ipRateLimiter{
		ips: make(map[string]*rate.Limiter),
		mu:  &sync.RWMutex{},
		r:   r,
		b:   b,
	}

	return i
}

func (i *ipRateLimiter) AddIP(ip string) *rate.Limiter {
	i.mu.Lock()
	defer i.mu.Unlock()

	limiter := rate.NewLimiter(i.r, i.b)

	i.ips[ip] = limiter

	return limiter
}

func (i *ipRateLimiter) GetLimiter(ip string) *rate.Limiter {
	i.mu.Lock()
	limiter, exists := i.ips[ip]

	if !exists {
		i.mu.Unlock()
		return i.AddIP(ip)
	}

	i.mu.Unlock()

	return limiter
}
