package viper

import (
	"RedRockMidAssessment-Consumer/core"
	"RedRockMidAssessment-Consumer/core/models"

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

	if err := v.Unmarshal(core.Config); err != nil {
		return nil, err
	}

	return Config, nil
}
