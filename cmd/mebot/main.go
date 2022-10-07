package main

import (
	"context"
	"fmt"
	"meBot/internal/config"
	"meBot/internal/handlers/telebot"
	"os"
	"os/signal"
	"path/filepath"
	"syscall"
	"time"

	"go.uber.org/zap"
	tele "gopkg.in/telebot.v3"
	"gopkg.in/yaml.v3"
)

func main() {
	// Cancel context if got Ctrl+C signal.
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	c := make(chan os.Signal)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-c
		// Run Cleanup
		cancel()
		os.Exit(1)
	}()

	if err := run(ctx); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func run(ctx context.Context) error {
	logger, _ := zap.NewDevelopment()
	defer func() { _ = logger.Sync() }()

	// === config
	config := &config.Config{}
	pwd, err := os.Getwd()
	if err != nil {
		logger.Fatal("get project work dir", zap.Error(err))
	}
	path := filepath.Join(pwd, "config.yml")
	bytes, err := os.ReadFile(path)
	if err != nil {
		logger.Fatal("read config file path: "+path, zap.Error(err))
	}
	err = yaml.Unmarshal(bytes, &config)
	if err != nil {
		logger.Fatal("unmarshal config", zap.Error(err))
	}
	// === config.end

	// telebot
	pref := tele.Settings{
		Token:  config.BotToken,
		Poller: &tele.LongPoller{Timeout: 10 * time.Second},
	}

	b, err := tele.NewBot(pref)
	if err != nil {
		logger.Fatal("init telebot", zap.Error(err))
	}

	b = telebot.InitHandlers(b)

	b.Start()
	defer func() {
		_, _ = b.Close()
	}()

	<-ctx.Done()

	return nil
}
