package http_transport

import (
	"errors"
	"gitlab.rusklimat.ru/ecom/go-lib/errs"
	"os"
	"os/signal"
	"syscall"
)

func (s *Server) Listen() error {
	var err error

	if (s.cfg.TLSPem != "" && s.cfg.TLSKey == "") || (s.cfg.TLSPem == "" && s.cfg.TLSKey != "") {
		return errors.New("(*Server).Listen() error: cfg.TLSPem or cfg.TLSKey doesn't have value")
	}

	switch s.cfg.TLSPem != "" && s.cfg.TLSKey != "" {
	case true:
		err = s.a.ListenTLS(s.cfg.Host, s.cfg.TLSPem, s.cfg.TLSKey)

	case false:
		err = s.a.Listen(s.cfg.Host)
	}

	return err
}

func (s *Server) GracefulShutdown(connectionsClosed chan struct{}) {
	sigint := make(chan os.Signal, 1)
	signal.Notify(sigint, syscall.SIGTERM, syscall.SIGINT, syscall.SIGQUIT)

	<-sigint

	if err := s.a.Shutdown(); err != nil {
		s.logs <- errs.NewError(errs.FatalLevel, err.Error()) // wrap not required
	}

	connectionsClosed <- struct{}{}
}
