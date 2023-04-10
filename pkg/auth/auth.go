package auth

import (
	"context"
	"fmt"
	"github.com/borealisdb/cli/pkg/config"
	"github.com/borealisdb/cli/pkg/templates"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"net/http"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"time"
)

type Auth struct {
	Log       *logrus.Entry
	CliConfig config.Wrapper
}

func (a *Auth) TokenListener() {
	gin.SetMode(gin.ReleaseMode)
	quit := make(chan int)
	r := gin.New()
	r.Use(gin.Recovery())
	_, err := templates.Asset(filepath.Join(config.TemplatesPath, "index.html"))
	if err != nil {
		a.Log.Fatal("could not load html template from pkg/templates/index.html")
	}

	r.LoadHTMLFiles(filepath.Join(config.TemplatesPath, "index.html"))

	r.GET("/token", func(c *gin.Context) {
		defer func() {
			quit <- 0
		}()
		idToken := c.Query("id_token")
		if err := a.CliConfig.SetToken(idToken); err != nil {
			a.Log.Errorf("could not save token: %v", err)
		}
		c.HTML(200, "index.html", gin.H{})
	})

	srv := &http.Server{
		Addr:    ":9999",
		Handler: r,
	}

	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			a.Log.Fatalf("listen: %s\n", err)
		}
	}()

	<-quit
	a.Log.Debugf("Shutdown Server ...")
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		a.Log.Fatal("Server Shutdown:", err)
	}
	select {
	case <-ctx.Done():
	}
	a.Log.Debugf("Server exiting")
}

func (a *Auth) Openbrowser(url string) error {
	var cmd *exec.Cmd
	switch runtime.GOOS {
	case "openbsd":
		fallthrough
	case "linux":
		cmd = exec.Command("xdg-open", url)
	case "darwin":
		cmd = exec.Command("open", url)
	case "windows":
		r := strings.NewReplacer("&", "^&")
		cmd = exec.Command("cmd", "/c", "start", r.Replace(url))
	}
	if cmd != nil {
		err := cmd.Start()
		if err != nil {
			return fmt.Errorf("failed to open browser: %v", err)
		}
		err = cmd.Wait()
		if err != nil {
			return fmt.Errorf("failed to wait for open browser command to finish: %v", err)
		}
		return nil
	} else {
		return fmt.Errorf("unsupported platform")
	}
}
