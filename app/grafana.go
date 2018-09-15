package app

import (
	"bytes"
	"context"
	"io"
	"os"
	"os/exec"
	"time"

	"github.com/revel/revel"
)

func init() {
	ctx, cancel := context.WithCancel(context.Background())
	var cmd *exec.Cmd
	AddExitEventHandler(func() {
		cancel()
		go func() {
			select {
			case <-time.After(30 * time.Second):
				cmd.Process.Kill()
			}
		}()
		revel.AppLog.Infof("shutting down grafana: %v", cmd.Wait())
	})
	revel.OnAppStart(func() {
		go func() {
			for {
				stdout := new(Writer)
				stdout.Str = "[Grafana]"
				cmd = exec.CommandContext(ctx, "grafana-server", "--config", "/usr/local/etc/grafana/grafana.ini", "--homepath", "/usr/local/opt/grafana/share/grafana", "cfg:default.paths.logs=/usr/local/var/log/grafana", "cfg:default.paths.data=/usr/local/var/lib/grafana", "cfg:default.paths.plugins=/usr/local/var/lib/grafana/plugins", "cfg:default.server.http_port=3001")
				cmd.Stdout = stdout
				cmd.Stderr = stdout
				cmd.Start()
				err := cmd.Wait()
				if err != nil {
					revel.AppLog.Errorf("Error running grafana %v", err)
				}
				<-time.After(30 * time.Second)
			}
		}()
	})
}

type Writer struct {
	Str string
}

func (w *Writer) Write(p []byte) (n int, err error) {
	var buf bytes.Buffer
	buf.WriteString(w.Str)
	n, err = buf.Write(p)
	io.Copy(os.Stdout, &buf)

	return
}
