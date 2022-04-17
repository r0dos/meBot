package main

import (
	"context"
	"fmt"
	"io/ioutil"
	"meBot/pkg/config"
	"os"
	"os/signal"
	"path/filepath"
	"strings"
	"syscall"
	"time"

	"go.uber.org/zap"
	tele "gopkg.in/telebot.v3"
	"gopkg.in/yaml.v2"
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
	bytes, err := ioutil.ReadFile(path)
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

	// Command: /ping
	b.Handle("/ping", func(c tele.Context) error {
		return c.Send("pong")
	})

	// Command: /me
	b.Handle("/me", func(c tele.Context) error {
		b.Delete(c.Message())

		userName := ""
		if c.Message().Sender.FirstName != "" {
			userName += c.Message().Sender.FirstName
		}
		if c.Message().Sender.LastName != "" {
			userName += " " + c.Message().Sender.LastName
		}

		args := c.Args()
		if len(args) == 0 {
			return c.Send(fmt.Sprintf("*%s* думает...", userName), tele.ModeMarkdown)
		}

		res := fmt.Sprintf("*%s* %s", userName, strings.Join(args, " "))

		return c.Send(res, tele.ModeMarkdown)
	})

	b.Start()
	defer b.Close()
	<-ctx.Done()
	return nil
}
