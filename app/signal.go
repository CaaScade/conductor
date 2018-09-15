package app

import (
	"context"
	"github.com/revel/revel"
	"os"
	"os/signal"
	"sync"
)

type ExitHandler func()

var handlers []ExitHandler
var lock sync.Mutex

func AddExitEventHandler(handler ExitHandler) {
	lock.Lock()
	defer lock.Unlock()
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
