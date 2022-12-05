package bot

import (
	"context"
	"meBot/pkg/log"
	"time"

	"go.uber.org/zap"
	"gopkg.in/telebot.v3"
)

type Storage interface {
	Inc(chatID, userID int64) error
	Get(chatID, userID int64) (int64, error)
	Reset(chatID, userID int64) error
}

type Registry interface {
	Add(key, code string) (context.Context, error)
	Update(key, code string) error
	Removal(key string) error
	Send(key, v string) error
}

type MeBot struct {
	bot      *telebot.Bot
	storage  Storage
	registry Registry
}

func NewMeBot(b *telebot.Bot, s Storage, r Registry) *MeBot {
	me := &MeBot{
		bot:      b,
		storage:  s,
		registry: r,
	}

	me.registerMiddlewares()
	me.registerHandlers()

	return me
}

func (m *MeBot) Start() {
	m.bot.Start()
}

func (m *MeBot) Close() {
	for i := 0; i < 5; i++ {
		c, err := m.bot.Close()
		if err != nil {
			log.Error("bot close", zap.Error(err))
		}

		if c {
			return
		}

		time.Sleep(time.Second)
	}
}
