package conf

import (
	"github.com/sirupsen/logrus"
	"gopkg.in/yaml.v3"
	"os"
)

type Conf struct {
	Lark *Lark `yaml:"lark"`
	Gpt  *Gpt  `yaml:"gpt"`
}

type Lark struct {
	AppID             string `yaml:"AppID"`
	AppSecret         string `yaml:"AppSecret"`
	EncryptKey        string `yaml:"EncryptKey"`
	VerificationToken string `yaml:"VerificationToken"`
}

type Gpt struct {
	ApiKey string `yaml:"ApiKey"`
}

var conf *Conf

func init() {
	bytes, err := os.ReadFile(os.Getenv("CONF_PATH"))
	if err != nil {
		logrus.Fatalf("os.Open err: %s", err.Error())
	}
	conf = &Conf{}
	err = yaml.Unmarshal(bytes, conf)
	if err != nil {
		logrus.Fatalf("yaml.Unmarshal err: %s", err.Error())
	}
}

func GetConf() *Conf {
	return conf
}
