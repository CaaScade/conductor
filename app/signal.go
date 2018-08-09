package app

import (
	"context"
	"os"
	"os/signal"

	"github.com/revel/revel"
)

func init() {
	sig := make(chan os.Signal, 1)
	signal.Notify(sig, os.Interrupt)
	go func() {
		<-sig
		revel.CurrentEngine.(*revel.GoHttpServer).Server.Shutdown(context.Background())
	}()
}
