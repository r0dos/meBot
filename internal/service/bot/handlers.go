package bot

import (
	"fmt"
	"strings"

	"gopkg.in/telebot.v3"
)

const (
	karmaMax = 3 // TODO in config
)

func (m *MeBot) registerHandlers() {
	// Command: /ping
	m.bot.Handle("/ping", m.handlerPong)

	groupOnly := m.bot.Group()
	groupOnly.Use(middlewareFromGroup)

	// Command: /me
	groupOnly.Handle("/me", m.handlerMe)

	// Command: /voteban
	groupOnly.Handle("/voteban", m.handlerVoteban)

	// Command: /karma
	groupOnly.Handle("/karma", m.handlerKarma)

	// Command: /ban
	groupOnly.Handle("/ban", m.handlerBan, middlewareCheckAdmins, middlewareCheckAdminsReplyTo)

	// Command: /unban
	groupOnly.Handle("/unban", m.handlerUnBan, middlewareCheckAdmins)

	// Command: /clearkarma
	groupOnly.Handle("/clearkarma", m.handlerClearKarma, middlewareCheckAdmins)

	//Command: /help
	groupOnly.Handle("/help", m.handlerHelp)

	// Command: /captcha
	groupOnly.Handle("/captcha", m.handlerToCaptcha, middlewareCheckAdmins)

	groupOnly.Handle(telebot.OnUserJoined, m.handlerWelcomeCaptcha)

	groupOnly.Handle("\fcode", m.handlerCode)
	groupOnly.Handle("\frefresh", m.handlerRefresh)
}

func (m *MeBot) handlerHelp(c telebot.Context) error {
	_, err := c.Bot().Reply(c.Message(),
		"Список команд:\n"+
			"/karma - покажет текущую карму\n"+
			"/ban - увеличит карму на 1 или забанит пользователя при достижении максимума\n"+
			"/unban - попробует разбанить пользователя, но карму не очистит"+
			"/clearkarma - очистит карму пользователю",
		telebot.ModeMarkdown)

	return err
}

func (m *MeBot) handlerPong(c telebot.Context) error {
	return c.Send("pong")
}

func (m *MeBot) handlerMe(c telebot.Context) error {
	err := c.Bot().Delete(c.Message())
	if err != nil {
		return fmt.Errorf("delete message: %v", err)
	}

	args := c.Args()
	if len(args) == 0 {
		return c.Send(fmt.Sprintf("*%s* думает...", getName(c.Message())), telebot.ModeMarkdown)
	}

	res := fmt.Sprintf("*%s* %s", getName(c.Message()), strings.Join(args, " "))

	if c.Message().ReplyTo != nil {
		_, err = c.Bot().Reply(c.Message().ReplyTo, res, telebot.ModeMarkdown)

		return err
	}

	return c.Send(res, telebot.ModeMarkdown)
}

func (m *MeBot) handlerVoteban(c telebot.Context) error {
	if c.Message().ReplyTo != nil {
		return c.Send(
			fmt.Sprintf("Давно уже пора забанить %s", getUsername(c.Message().ReplyTo)),
			telebot.ModeMarkdown)
	}

	return c.Send("Я тоже считаю, что это заслуживает бана!")
}

func (m *MeBot) handlerKarma(c telebot.Context) error {
	if c.Message().ReplyTo == nil {
		_, err := c.Bot().Reply(c.Message(), "У меня хорошая карма, спасибо")

		return err
	}

	karmaCurrent, err := m.storage.Get(c.Message().Chat.ID, c.Message().ReplyTo.Sender.ID)
	if err != nil {
		return fmt.Errorf("handler karma: %v", err)
	}

	_, err = c.Bot().Reply(
		c.Message(),
		fmt.Sprintf("Карма %s %d/%d", getUsername(c.Message().ReplyTo), karmaCurrent, karmaMax),
		telebot.ModeMarkdown,
	)

	return err
}

func (m *MeBot) handlerBan(c telebot.Context) error {
	if c.Message().ReplyTo == nil {
		_, err := c.Bot().Reply(c.Message(), "Пальцем показывать не хорошо, но иногда надо")

		return err
	}

	if err := m.storage.Inc(c.Message().Chat.ID, c.Message().ReplyTo.Sender.ID); err != nil {
		return fmt.Errorf("inc in ban: %v", err)
	}

	karmaCurrent, err := m.storage.Get(c.Message().Chat.ID, c.Message().ReplyTo.Sender.ID)
	if err != nil {
		return fmt.Errorf("handler karma: %v", err)
	}

	if karmaCurrent >= karmaMax {
		return c.Bot().Ban(c.Message().Chat, &telebot.ChatMember{
			User:            c.Message().ReplyTo.Sender,
			RestrictedUntil: telebot.Forever(),
		})
	}

	_, err = c.Bot().Reply(
		c.Message(),
		fmt.Sprintf("Карма %s %d/%d", getUsername(c.Message().ReplyTo), karmaCurrent, karmaMax),
		telebot.ModeMarkdown,
	)

	return err
}

func (m *MeBot) handlerUnBan(c telebot.Context) error {
	if c.Message().ReplyTo == nil {
		_, err := c.Bot().Reply(c.Message(), "И что?...")

		return err
	}

	return c.Bot().Unban(c.Message().Chat, c.Message().ReplyTo.Sender)
}

func (m *MeBot) handlerClearKarma(c telebot.Context) error {
	if c.Message().ReplyTo == nil {
		_, err := c.Bot().Reply(c.Message(), "И тебе лайк")

		return err
	}

	if err := m.storage.Reset(c.Message().Chat.ID, c.Message().ReplyTo.Sender.ID); err != nil {
		return fmt.Errorf("reset karma: %v", err)
	}

	return m.handlerKarma(c)
}
