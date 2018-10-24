package app

import (
	"context"
	"os"
	"os/exec"
	"time"

	"github.com/revel/revel"
)

func init() {
	ctx := context.Background()
	var cmd *exec.Cmd
	AddExitEventHandler(func() {
		ctx.Done()
		revel.AppLog.Infof("shutting down prometheus: %v", cmd.Wait())
	})
	revel.OnAppStart(func() {
		go func() {
			for {
				cmd = exec.CommandContext(ctx, "prometheus", "--config.file", "prometheus.yaml")
				cmd.Stdout = os.Stdout
				cmd.Stderr = os.Stderr
				cmd.Start()
				err := cmd.Wait()
				if err != nil {
					revel.AppLog.Errorf("Error running prometheus %v", err)
				}
				<-time.After(30 * time.Second)
			}
		}()
	})
}
