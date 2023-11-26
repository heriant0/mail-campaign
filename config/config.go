package config

import "github.com/spf13/viper"

type Config struct {
	ConfigSmtpHost     string `mapstructure:"CONFIG_SMTP_HOST"`
	ConfigSmtpPort     int    `mapstructure:"CONFIG_SMTP_PORT"`
	ConfigSenderName   string `mapstructure:"CONFIG_SENDER_NAME"`
	ConfigAuthEmail    string `mapstructure:"CONFIG_AUTH_EMAIL"`
	ConfigAuthPassword string `mapstructure:"CONFIG_AUTH_PASSWORD"`
	AppPort            string `mapstructure:"APP_PORT"`
	MailPort           string `mapstructure:"MAIL_PORT"`
	BaseUrl            string `mapstructure:"BASE_URL"`
}

func LoadConfig() (Config, error) {
	var config Config

	viper.AddConfigPath(".")
	viper.SetConfigName("app")
	viper.SetConfigType("env")
	viper.AutomaticEnv()

	err := viper.ReadInConfig()
	if err != nil {
		return config, err
	}

	err = viper.Unmarshal(&config)
	if err != nil {
		return config, err
	}

	return config, nil
}
