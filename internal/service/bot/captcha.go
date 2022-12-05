package bot

import (
	"context"
	"errors"
	"fmt"
	"meBot/pkg/log"
	"os"
	"time"

	"github.com/dchest/captcha"
	"go.uber.org/zap"
	"gopkg.in/telebot.v3"
)

const countBtn = 5

const waitTime = time.Minute * 2

func (m *MeBot) handlerToCaptcha(c telebot.Context) error {
	if c.Message().ReplyTo == nil {
		return nil
	}

	return m.handlerCaptcha(c.Message().ReplyTo)
}

func (m *MeBot) handlerWelcomeCaptcha(c telebot.Context) error {
	return m.handlerCaptcha(c.Message())
}

func (m *MeBot) handlerCaptcha(message *telebot.Message) error {
	dig := captcha.RandomDigits(countBtn)
	img := captcha.NewImage("", dig, 240, 80)

	codeInt := byteToInt(dig)
	codeStr := intToString(codeInt)

	filePath, err := saveFile(codeStr, img)
	if err != nil {
		return fmt.Errorf("save capcha: %v", err)
	}
	defer func() {
		if err := os.Remove(filePath); err != nil {
			log.Error("remove captcha", zap.Error(err))
		}
	}()

	keyboard := getKeyboard(codeInt)

	username := getUsername(message)
	text := fmt.Sprintf("Привет, %s! Пройди анти-спам проверку за 2 мин.", username)

	photo := &telebot.Photo{
		File: telebot.File{
			FileLocal: filePath,
		},
		Caption: text,
	}

	msg, err := m.bot.Reply(message, photo, keyboard)
	if err != nil {
		return fmt.Errorf("send reply captcha: %v", err)
	}

	ctx, err := m.registry.Add(getKey(msg), codeStr)
	if err != nil {
		_ = m.bot.Delete(msg)

		return fmt.Errorf("registry add: %v", err)
	}

	go m.wait(ctx, msg, getKey(msg))

	return nil
}

func (m *MeBot) handlerRefresh(c telebot.Context) error {
	dig := captcha.RandomDigits(countBtn)
	img := captcha.NewImage("", dig, 240, 80)

	codeInt := byteToInt(dig)
	codeStr := intToString(codeInt)

	filePath, err := saveFile(codeStr, img)
	if err != nil {
		return fmt.Errorf("save capcha: %v", err)
	}
	defer func() {
		if err := os.Remove(filePath); err != nil {
			log.Error("remove captcha", zap.Error(err))
		}
	}()

	keyboard := getKeyboard(codeInt)

	photo := &telebot.Photo{
		File: telebot.File{
			FileLocal: filePath,
		},
		Caption: c.Text(),
	}

	if err := m.registry.Update(getKey(c.Message()), codeStr); err != nil {
		return fmt.Errorf("update captcha: %v", err)
	}

	return c.Edit(photo, keyboard)
}

func (m *MeBot) handlerCode(c telebot.Context) error {
	data := c.Data()
	if len(data) > 1 {
		return errors.New("bad data")
	}

	return m.registry.Send(getKey(c.Message()), data)
}

func (m *MeBot) handlerOnChat(c telebot.Context) error {
	fmt.Println("chat id:", c.Chat().ID)
	fmt.Println("sender id:", c.Message().Sender.ID)
	fmt.Printf("%#v\n", c.Message())

	return nil
}

func (m *MeBot) wait(ctx context.Context, msg *telebot.Message, key string) {
	defer func() {
		if err := m.registry.Removal(key); err != nil {
			log.Error("remove from registry", zap.Error(err))
		}
	}()

	defer func() {
		if err := m.bot.Delete(msg); err != nil {
			log.Error("delete message", zap.Error(err))
		}
	}()

	if err := m.bot.Restrict(msg.Chat, &telebot.ChatMember{
		Rights:          telebot.NoRights(),
		User:            msg.ReplyTo.Sender,
		RestrictedUntil: time.Now().Add(waitTime).Unix(),
	}); err != nil {
		log.Error("mute user", zap.Error(err))

		return
	}

	timer := time.NewTimer(waitTime)
	defer timer.Stop()

	select {
	case <-ctx.Done():
		if err := m.bot.Restrict(msg.Chat, &telebot.ChatMember{
			Rights:          telebot.NoRestrictions(),
			User:            msg.ReplyTo.Sender,
			RestrictedUntil: time.Now().Add(waitTime).Unix(),
		}); err != nil {
			log.Error("unmute user", zap.Error(err))
		}

	case <-timer.C:
		if err := m.bot.Ban(msg.Chat, &telebot.ChatMember{
			User:            msg.ReplyTo.Sender,
			RestrictedUntil: time.Now().Add(time.Hour).Unix(),
		}); err != nil {
			log.Error("ban user", zap.Error(err))

			return
		}

		if _, err := m.bot.Send(msg.Chat, fmt.Sprintf("%s не прошёл проверку и был забанен", getUsername(msg.ReplyTo))); err != nil {
			log.Error("send message ban user", zap.Error(err))
		}
	}
}
