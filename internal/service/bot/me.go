package bot

import (
	"go.uber.org/zap"
	"gopkg.in/telebot.v3"
	"meBot/pkg/log"
)

type Storage interface {
	Inc(chatID, userID int64) error
	Get(chatID, userID int64) (int64, error)
	Reset(chatID, userID int64) error
}

type MeBot struct {
	bot     *telebot.Bot
	storage Storage
}

func NewMeBot(b *telebot.Bot, s Storage) *MeBot {
	me := &MeBot{
		bot:     b,
		storage: s,
	}

	me.registerMiddlewares()
	me.registerHandlers()

	return me
}

func (m *MeBot) Start() {
	m.bot.Start()
}

func (m *MeBot) Close() {
	_, err := m.bot.Close()
	if err != nil {
		log.Error("bot close", zap.Error(err))
	}
}
