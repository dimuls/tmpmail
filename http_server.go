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

	"github.com/julienschmidt/httprouter"
	"github.com/jxskiss/base62"
	"github.com/rs/cors"
	"go.uber.org/zap"

	"tmpmail/entity"
)

//go:embed all:ui/dist/*
var uiFS embed.FS

func init() {

}

type HTTPServerStorage interface {
	CreateAccount(token, username string) error
	ProlongAccount(token string) error
	Account(token string) (entity.Account, error)
	RemoveAccount(token string) error
}

type HTTPServer struct {
	server  *http.Server
	storage HTTPServerStorage
	ui      fs.FS
	api     http.Handler
	logger  *zap.Logger
}

func NewHTTPServer(l *zap.Logger, s HTTPServerStorage, addr string, tc *tls.Config) *HTTPServer {
	ui, _ := fs.Sub(uiFS, "ui/dist")

	srv := &HTTPServer{
		storage: s,
		ui:      ui,
		logger:  l,
	}

	api := httprouter.New()
	api.POST("/api/account", srv.postAPIAccount)
	api.PATCH("/api/account", srv.patchAPIAccount)
	api.GET("/api/account", srv.getAPIAccount)
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
	tokenHeader = "X-TOKEN"
	tokenLength = 128
	emailLength = 8
)

func generateRandomString(length int) string {
	b := make([]byte, length)
	if _, err := rand.Read(b); err != nil {
		return ""
	}
	return base62.EncodeToString(b)
}

func (s *HTTPServer) postAPIAccount(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	token := generateRandomString(tokenLength)
	email := strings.ToLower(generateRandomString(emailLength))
	err := s.storage.CreateAccount(token, email)
	if err != nil {
		s.logger.Warn("create email in storage", zap.Error(err))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(token)
}

func (s *HTTPServer) patchAPIAccount(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	token := r.Header.Get(tokenHeader)
	err := s.storage.ProlongAccount(token)
	if err != nil {
		s.logger.Error("prolong account in storage", zap.Error(err))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
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
