package viper

import (
	"RedRockMidAssessment/core"
	"RedRockMidAssessment/core/models"

	"github.com/spf13/viper"
)

func InitConfig(path string) error {
	core.Config = new(models.Config)

	v := viper.New()
	v.SetConfigType("yaml")
	v.SetConfigFile(path)

	if err := v.ReadInConfig(); err != nil {
		return err
	}

	if err := v.Unmarshal(core.Config); err != nil {
		return err
	}

	return nil
}
