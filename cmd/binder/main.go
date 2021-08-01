package main

import (
	"context"
	"github.com/alipniczkij/binder/internal/commands"
	"github.com/alipniczkij/binder/internal/commands/label"
	"github.com/alipniczkij/binder/internal/commands/subscribe"
	"github.com/alipniczkij/binder/internal/commands/unlabel"
	"github.com/alipniczkij/binder/internal/commands/unsubscribe"
	"github.com/alipniczkij/binder/internal/handler"
	"github.com/alipniczkij/binder/internal/sender/slack"
	"github.com/alipniczkij/binder/internal/storage/bbolt"
	"github.com/alipniczkij/binder/pkg/config"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

var configPath = "../../binder.json"

func main() {
	cfg := config.LoadConfig(configPath)

	if cfg.LogPath != "" {
		f, err := os.OpenFile(cfg.LogPath, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
		if err != nil {
			log.Panicf("Error opening log file: %v", err)
		}
		defer f.Close()
		log.SetOutput(f)
	}

	mapper := bbolt.New(cfg.MappingPath)
	defer mapper.Close()

	cmds := map[string]commands.Handler{
		commands.Subscribe:   subscribe.New(mapper),
		commands.Unsubscribe: unsubscribe.New(mapper),
		commands.Label:       label.New(mapper),
		commands.Unlabel:     unlabel.New(mapper),
	}

	sender := slack.New(cfg.Slack, mapper)
	handler := handler.New(sender, cmds)

	address := net.JoinHostPort(cfg.Server.Host, cfg.Server.Port)
	serv := http.Server{
		Addr:    address,
		Handler: handler.RegisterPublicHTTP(),
	}
	go serv.ListenAndServe()

	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt, syscall.SIGTERM, syscall.SIGINT)

	<-interrupt

	timeout, cancelFunc := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancelFunc()
	err := serv.Shutdown(timeout)
	if err != nil {
		log.Printf("Error when shutdown app: %v", err)
	}
}
