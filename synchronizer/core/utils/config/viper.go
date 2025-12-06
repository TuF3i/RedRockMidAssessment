package viper

import (
	"RedRockMidAssessment-Synchronizer/core/models"

	"github.com/spf13/viper"
)

func InitConfig(path string) (*models.Config, error) {
	Config := new(models.Config)

	v := viper.New()
	v.SetConfigType("yaml")
	v.SetConfigFile(path)

	if err := v.ReadInConfig(); err != nil {
		return Config, err
	}

	if err := v.Unmarshal(Config); err != nil {
		return Config, err
	}

	return Config, nil
}
