package config

// Config is the options of r0bot.
type Config struct {
	BotToken   string `yaml:"bot_token"`
	SendUsers  string `yaml:"send_users"`
	BinanceApi string `yaml:"binance_api"`
	BinanceSec string `yaml:"binance_sec"`
}
