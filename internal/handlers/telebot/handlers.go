package telebot

import (
	"fmt"
	"gopkg.in/telebot.v3"
	"strings"
)

func InitHandlers(b *telebot.Bot) *telebot.Bot {
	// Command: /ping
	b.Handle("/ping", func(c telebot.Context) error {
		return c.Send("pong")
	})

	// Command: /me
	b.Handle("/me", func(c telebot.Context) error {
		err := b.Delete(c.Message())
		if err != nil {
			return err
		}

		var userName string

		if c.Message().Sender.FirstName != "" {
			userName += c.Message().Sender.FirstName
		}

		if c.Message().Sender.LastName != "" {
			userName += " " + c.Message().Sender.LastName
		}

		args := c.Args()
		if len(args) == 0 {
			return c.Send(fmt.Sprintf("*%s* думает...", userName), telebot.ModeMarkdown)
		}

		res := fmt.Sprintf("*%s* %s", userName, strings.Join(args, " "))

		if c.Message().ReplyTo != nil {
			msg, _ := b.Reply(c.Message().ReplyTo, res, telebot.ModeMarkdown)
			return c.Send(msg)
		}

		return c.Send(res, telebot.ModeMarkdown)
	})

	b.Handle("/voteban", func(c telebot.Context) error {
		if c.Message().ReplyTo != nil {
			if c.Message().ReplyTo.Sender.Username != "" {
				return c.Send(
					fmt.Sprintf("Давно уже пора забанить %s", c.Message().ReplyTo.Sender.Username),
					telebot.ModeMarkdown)
			}

			var userName string

			if c.Message().Sender.FirstName != "" {
				userName += c.Message().Sender.FirstName
			}

			if c.Message().Sender.LastName != "" {
				userName += " " + c.Message().Sender.LastName
			}
			return c.Send(
				fmt.Sprintf("Давно уже пора забанить %s", userName),
				telebot.ModeMarkdown)
		}

		return c.Send("Я тоже считаю, что это заслуживает бана!")
	})

	return b
}
