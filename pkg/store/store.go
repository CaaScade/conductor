package store

import (
	"time"

	log "github.com/sirupsen/logrus"

	"conductor/internal/etcd/embed"
)

func init() {
	go Init()
}

func Init() {
	cfg := embed.NewConfig()
	cfg.Dir = "default.etcd"
	e, err := embed.StartEtcd(cfg)
	if err != nil {
		log.Fatal(err)
	}
	defer e.Close()
	select {
	case <-e.Server.ReadyNotify():
		log.Infof("Store is ready!")
	case <-time.After(60 * time.Second):
		e.Server.Stop() // trigger a shutdown
		log.Printf("Store took too long to setup!")
	}
	log.Fatal(<-e.Err())
}