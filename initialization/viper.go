/**
 @author: 15973
 @date: 2022/07/13
 @note:
**/
package initialization

import (
	"fmt"
	"github.com/spf13/viper"
)

func InitializeViper() {
	viper.SetConfigName("config")
	viper.AddConfigPath("./config")
	err := viper.ReadInConfig() // Find and read the config file
	if err != nil {             // Handle errors reading the config file
		panic(fmt.Errorf("fatal error config file: %w", err))
	}
}
