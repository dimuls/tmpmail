package main

import (
	"context"
	"crypto/tls"
	"errors"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/spf13/cobra"
	"go.uber.org/zap"
	"golang.org/x/crypto/acme/autocert"

	"tmpmail"
	"tmpmail/redis"
)

const (
	domain     = "tmp-mail.ru"
	mailDomain = "smtp.tmp-mail.ru"
	certsCache = "/var/lib/tmpmail/certs"
	emailTTL   = 10 * time.Minute
)

var (
	smtpAddr, httpAddr string
	redisAddr          string
)

func server(_ *cobra.Command, _ []string) {

	logger, err := zap.NewDevelopment()
	if err != nil {
		zap.L().Fatal("init logger", zap.Error(err))
	}

	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer cancel()

	rs := redis.NewStorage(redisAddr, emailTTL)
	defer func() {
		err = rs.Close()
		if err != nil {
			logger.Error("redis storage close", zap.Error(err))
			return
		}
		logger.Info("redis storage closed", zap.Error(err))
	}()

	cm := autocert.Manager{
		Prompt:     autocert.AcceptTOS,
		HostPolicy: autocert.HostWhitelist(domain, mailDomain),
		Cache:      autocert.DirCache(certsCache),
	}

	cmHTTPSrv := &http.Server{
		Addr:    "0.0.0.0:80",
		Handler: cm.HTTPHandler(nil),
	}

	go func() {
		err := cmHTTPSrv.ListenAndServe()
		if err != nil {
			if errors.Is(err, http.ErrServerClosed) {
				return
			}
			logger.Error("cert manager http server start", zap.Error(err))
		}
	}()

	defer func() {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		err = cmHTTPSrv.Shutdown(ctx)
		if err != nil {
			logger.Error("cert manager http server stop", zap.Error(err))
			return
		}
		logger.Info("cert manager http server shutdown")
	}()

	tlsCfg := cm.TLSConfig()
	origGetCertificate := tlsCfg.GetCertificate
	tlsCfg.GetCertificate = func(info *tls.ClientHelloInfo) (*tls.Certificate, error) {
		info.ServerName = mailDomain
		return origGetCertificate(info)
	}
	tlsCfg.ServerName = mailDomain

	smtpSrv := tmpmail.NewSMTPServer(logger, rs, tlsCfg, smtpAddr, domain, mailDomain)

	go func() {
		defer cancel()
		err := smtpSrv.ListenAndServe()
		if err != nil {
			if errors.Is(err, http.ErrServerClosed) {
				return
			}
			logger.Error("smtp server listen and serve", zap.Error(err))
		}
	}()

	defer func() {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		err = smtpSrv.Shutdown(ctx)
		if err != nil {
			logger.Error("smtp server shutdown", zap.Error(err))
			return
		}
		logger.Info("smtp server shutdown")
	}()

	tlsCfg = cm.TLSConfig()
	tlsCfg.ServerName = domain

	httpSrv := tmpmail.NewHTTPServer(logger, rs, httpAddr, tlsCfg)

	go func() {
		defer cancel()
		err := httpSrv.ListenAndServe()
		if err != nil {
			if errors.Is(err, http.ErrServerClosed) {
				return
			}
			logger.Error("http server start", zap.Error(err))
		}
	}()

	defer func() {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		err = httpSrv.Shutdown(ctx)
		if err != nil {
			logger.Error("http server stop", zap.Error(err))
			return
		}
		logger.Info("http server shutdown")
	}()

	<-ctx.Done()
}

func main() {
	rootCmd := &cobra.Command{
		Use: "tmpmail",
	}

	serverCmd := &cobra.Command{
		Use: "server",
		Run: server,
	}

	serverCmd.Flags().StringVar(&smtpAddr, "smtp-addr", "0.0.0.0:25", "")
	serverCmd.Flags().StringVar(&httpAddr, "http-addr", "0.0.0.0:443", "")
	serverCmd.Flags().StringVar(&redisAddr, "redis-addr", "127.0.0.1:6379", "")

	rootCmd.AddCommand(serverCmd)

	err := rootCmd.Execute()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
