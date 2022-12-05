package bot

import (
	"fmt"
	"math/rand"
	"os"
	"strings"
	"time"

	"github.com/dchest/captcha"
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
func byteToInt(b []byte) []int {
	res := make([]int, 0, len(b))

	for _, item := range b {
		res = append(res, int(item))
	}

	return res
}

func intToString(i []int) string {
	b := strings.Builder{}
	for _, item := range i {
		b.WriteString(fmt.Sprint(item))
	}

	return b.String()
}

func saveFile(name string, img *captcha.Image) (string, error) {
	filePath := fmt.Sprintf("tmp/%s.png", name)

	file, err := os.OpenFile(filePath, os.O_WRONLY|os.O_CREATE, 0600)
	if err != nil {
		return "", fmt.Errorf("open file: %v", err)
	}
	defer file.Close()

	_, err = img.WriteTo(file)
	if err != nil {
		return "", fmt.Errorf("write file: %v", err)
	}

	return filePath, nil
}

var random = rand.New(rand.NewSource(time.Now().UnixNano()))

func getKeyboard(code []int) *telebot.ReplyMarkup {
	random.Shuffle(len(code), func(i, j int) {
		code[i], code[j] = code[j], code[i]
	})

	btnMap := make(map[int]struct{}, countBtn)
	for i := 0; len(btnMap) < 5; i++ {
		if i < len(code) {
			btnMap[code[i]] = struct{}{}

			continue
		}

		btnMap[random.Intn(10)] = struct{}{}
	}

	keyboard := &telebot.ReplyMarkup{}
	row := telebot.Row{}

	for key, _ := range btnMap {
		row = append(row, keyboard.Data(fmt.Sprint(key), "code", fmt.Sprint(key)))
	}

	row = append(row, keyboard.Data("\xF0\x9F\x94\x84", "refresh"))

	keyboard.Inline(row)

	return keyboard
}

func getKey(msg *telebot.Message) string {
	return fmt.Sprintf("%d_%d", msg.Chat.ID, msg.ReplyTo.Sender.ID)
}
