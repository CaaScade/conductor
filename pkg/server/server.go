package server

import (
	"fmt"
	"net/http"
	"time"

	"conductor/pkg/config"
	"conductor/pkg/router"

	log "github.com/sirupsen/logrus"
)

func Run() error {
	cfg := config.GlobalConfig()

	s := &http.Server{
		Handler:      &router.Router{},
		Addr:         fmt.Sprintf("%s:%d", cfg.Addr, cfg.Port),
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
	}
	log.Infof("listening on %s:%d", cfg.Addr, cfg.Port)
	return s.ListenAndServe()
}
