package bot

import (
	"fmt"

	"gopkg.in/telebot.v3"
)

func getUsername(mes *telebot.Message) (userName string) {
	if mes == nil || mes.Sender == nil {
		return
	}

	if mes.Sender.Username != "" {
		return fmt.Sprintf("@%s", mes.Sender.Username)
	}

	return getName(mes)
}

func getName(mes *telebot.Message) (userName string) {
	if mes == nil || mes.Sender == nil {
		return
	}

	if mes.Sender.FirstName != "" {
		userName += mes.Sender.FirstName
	}

	if mes.Sender.LastName != "" {
		userName += " " + mes.Sender.LastName
	}

	return
}
