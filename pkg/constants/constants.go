package constants

import (
	"fmt"

	"github.com/spf13/viper"
)

var (
	OpenKommanderConfigFilename = ".openkommander_config"
)

func init() {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath("config")

	if err := viper.ReadInConfig(); err != nil {
		fmt.Printf("Error reading config file: %v", err)
		return
	}

	if sessionFilename := viper.GetString("openkommander.session_filename"); sessionFilename != "" {
		OpenKommanderConfigFilename = sessionFilename
	}
}
