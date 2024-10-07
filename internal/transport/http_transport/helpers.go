package http_transport

import "gitlab.rusklimat.ru/ecom/go-lib/errs"

func (s *Server) collectLog(e *errs.Error) {
	// в канал с логами отправляем только ошибки, которые нужно отправлять в Sentry
	if e.GetLevel() >= errs.WarnLevel {
		s.logs <- e
	}
}
