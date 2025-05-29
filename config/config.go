package config

import (
	"fmt"

	C "github.com/Golevka2001/bassit/constant"

	"github.com/spf13/viper"
)

var (
	IsModified = false
	Config     = config{}
)

type config struct {
	Tuning [C.StringCnt]string

	PositionMarkerChar string
}

func Unmarshal() error {
	err := viper.Unmarshal(&Config)
	if err != nil {
		return fmt.Errorf("failed to unmarshal config: %w", err)
	}
	return nil
}
