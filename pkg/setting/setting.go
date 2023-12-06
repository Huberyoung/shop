package setting

import (
	"github.com/spf13/viper"
)

type Setting struct {
	vp *viper.Viper
}

type DatabaseSettingS struct {
	Username  string
	Password  string
	Protocol  string
	Host      string
	Port      string
	DBName    string
	ParseTime bool
	Loc       string
	Charset   string
}

func NewSetting() (*Setting, error) {
	vp := viper.New()
	vp.AddConfigPath("configs/")
	vp.SetConfigName("config")
	vp.SetConfigType("yaml")
	err := vp.ReadInConfig()
	if err != nil {
		return nil, err
	}
	return &Setting{vp}, nil
}

func (s *Setting) ReadSection(k string, v any) error {
	return s.vp.UnmarshalKey(k, v)
}
