package app

import (
	"context"
	"os"
	"os/signal"

	"github.com/revel/revel"
)

type ExitHandler func()

var handlers []ExitHandler

func AddExitEventHandler(handler ExitHandler) {
	handlers = append(handlers, handler)
}

func init() {
	handlers = []ExitHandler{}
	sig := make(chan os.Signal, 1)
	signal.Notify(sig, os.Interrupt)
	go func() {
		<-sig
		for _, h := range handlers {
			h()
		}
		revel.CurrentEngine.(*revel.GoHttpServer).Server.Shutdown(context.Background())
	}()
}
