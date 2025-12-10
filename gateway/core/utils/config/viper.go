package viper

import (
	"RedRockMidAssessment/core/models"

	"github.com/spf13/viper"
)

func InitConfig(path string) (*models.Config, error) {
	Config := new(models.Config)

	v := viper.New()
	v.SetConfigType("yaml")
	v.SetConfigFile(path)

	if err := v.ReadInConfig(); err != nil {
		return nil, err
	}

	if err := v.Unmarshal(Config); err != nil {
		return nil, err
	}

	return Config, nil
}
