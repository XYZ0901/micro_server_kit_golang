package initialize

import (
	"github.com/spf13/viper"
)

func cfgInit() {
	//TODO if you changed fileName please rename it
	configFileName := "config_dev"
	if !getEnvInfo("SERVER_DEBUG") {
		configFileName = "config_pro"
	}
	viper.SetConfigName(configFileName)
	viper.SetConfigType("yaml")
	viper.AddConfigPath("./config")
	err := viper.ReadInConfig()
	if err != nil {
		Logger.Sugar().Fatal("[init.cfgInit] viper.ReadInConfig error:", err)
	}
	if err = viper.Unmarshal(&Cfg); err != nil {
		Logger.Sugar().Fatal("[init.cfgInit] viper.Unmarshal error:", err)
	}
	go func() {
		for {
			viper.WatchConfig()
			//TODO if config file has changed what you want to do
		}
	}()
}

func getEnvInfo(env string) bool {
	viper.AutomaticEnv()
	return viper.GetBool(env)
}
