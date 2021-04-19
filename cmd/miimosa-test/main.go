package main

import (
	"github.com/cockroachdb/errors"
	"github.com/kelseyhightower/envconfig"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"log"
	"miimosa-test/internal/config"
	"miimosa-test/internal/server"
	"miimosa-test/pkg/sessions"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"
)

func main() {

	var cfg = config.App{}
	err := envconfig.Process("", &cfg)

	if err != nil {
		logrus.WithError(err).Fatalf("unable to load configuration file")
	}

	lis, err := net.Listen("tcp", ":"+cfg.Port)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer()
	sessions.RegisterSessionsServer(s, server.NewSessionServer(cfg))

	logrus.WithField("port", cfg.Port).Info("starting server")

	// this is only here because of https://github.com/gusaul/grpcox great tool for interactively test grpc
	// this must not be present in production
	reflection.Register(s)

	err = gracefulServe(
		func() error {
			return s.Serve(lis)
		},
		s.GracefulStop,
	)

	if err != nil {
		logrus.WithError(err).Fatal("server stopped abnormally")
	}

	logrus.Info("stopped")
}

func gracefulServe(start func() error, shutdown func()) error {
	var (
		stopChan = make(chan os.Signal)
		closed   = make(chan error)
	)

	go func() {
		signal.Notify(stopChan, os.Interrupt, syscall.SIGTERM)

		<-stopChan

		logrus.Info("stopping the server")

		shutdown()

		closed <- nil
	}()

	if err := start(); err != http.ErrServerClosed {
		return errors.Wrap(err, "unable to start server")
	}

	<-closed

	return nil
}
