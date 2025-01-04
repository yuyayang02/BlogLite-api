package config

import (
	"errors"
	"github.com/spf13/viper"
	"log"
)

var Cfg Config

type Config struct {
	PORT              int    `toml:"port" mapstructure:"port"`
	Mode              string `toml:"mode" mapstructure:"mode"`
	DatabaseDNS       string `toml:"database_dns" mapstructure:"database_dns"`
	AuthSecretKey     string `toml:"auth_secret_key" mapstructure:"auth_secret_key"`
	AuthAdminPassword string `toml:"auth_admin_password" mapstructure:"auth_admin_password"`
}

func (c Config) Validate() error {
	if c.Mode != "debug" && c.Mode != "release" {
		return errors.New("未设置或设置了错误的运行模式")
	}

	if c.PORT == 0 {
		return errors.New("未设置运行端口")
	}

	if c.DatabaseDNS == "" {
		return errors.New("未配置数据库")
	}

	if c.AuthSecretKey == "" {
		return errors.New("未设置身份认证密钥")
	}

	if c.AuthAdminPassword == "" {
		return errors.New("未设置身份认证密码")
	}
	return nil
}

func Init(config string) {

	viper.SetEnvPrefix("BL")
	viper.AutomaticEnv()

	viper.SetConfigFile(config)
	err := viper.ReadInConfig()
	if err != nil {
		log.Fatal("读取配置失败", "err", err)
	}

	err = viper.Unmarshal(&Cfg)
	if err != nil {
		log.Fatal("写入配置失败", "err", err)
	}

	err = Cfg.Validate()
	if err != nil {
		log.Fatal("初始化参数时发生错误", "err", err)
	}
}
