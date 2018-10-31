package server

import (
	"conductor/pkg/router"
	"fmt"
	"time"

	"conductor/pkg/config"
	log "github.com/sirupsen/logrus"
	"net/http"
	"../router/router.go"
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
