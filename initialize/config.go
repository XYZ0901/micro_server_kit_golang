package initialize

import (
	"fmt"
	"github.com/nacos-group/nacos-sdk-go/clients"
	"github.com/nacos-group/nacos-sdk-go/common/constant"
	"github.com/nacos-group/nacos-sdk-go/vo"
	"github.com/spf13/viper"
	"gopkg.in/yaml.v3"
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
	nc := nacosConfig{}
	if err = viper.Unmarshal(&nc); err != nil {
		Logger.Sugar().Fatal("[init.cfgInit] viper.Unmarshal error:", err)
	}
	go func() {
		for {
			viper.WatchConfig()
			//TODO if config file has changed what you want to do
		}
	}()
	getConfigFromNacos(nc)
}

func getConfigFromNacos(nc nacosConfig) {
	nServerCfg := nc.Nacos.Server
	nClientCfg := nc.Nacos.Client
	nDataCfg := nc.Nacos.Data

	sc := []constant.ServerConfig{
		*constant.NewServerConfig(
			nServerCfg.Host, uint64(nServerCfg.Port),
			constant.WithScheme(nServerCfg.Scheme),
			constant.WithContextPath(nServerCfg.ContextPath)),
	}

	cc := *constant.NewClientConfig(
		constant.WithNamespaceId(nClientCfg.NamespaceId),
		constant.WithTimeoutMs(uint64(nClientCfg.TimeoutMs)),
		constant.WithNotLoadCacheAtStart(nClientCfg.NotLoadCacheAtStart),
		constant.WithLogDir(nClientCfg.LogDir),
		constant.WithCacheDir(nClientCfg.CacheDir),
		constant.WithLogLevel(nClientCfg.LogLevel),
	)

	client, err := clients.NewConfigClient(vo.NacosClientParam{
		ClientConfig:  &cc,
		ServerConfigs: sc,
	})
	if err != nil {
		panic(err)
	}

	content, err := client.GetConfig(vo.ConfigParam{
		DataId: nDataCfg.DataId,
		Group:  nDataCfg.Group,
		Type:   vo.ConfigType(nDataCfg.Type),
	})
	if err != nil {
		panic(err)
	}
	fmt.Println(content)

	err = yaml.Unmarshal([]byte(content), &Cfg)
	if err != nil {
		panic(err)
	}

	err = client.ListenConfig(vo.ConfigParam{
		DataId: nDataCfg.DataId,
		Group:  nDataCfg.Group,
		Type:   vo.ConfigType(nDataCfg.Type),
		OnChange: func(namespace, group, dataId, data string) {
			fmt.Println("config changed group:" + group + ", dataId:" + dataId + ", content:" + data)
			err = yaml.Unmarshal([]byte(data), &Cfg)
			if err != nil {
				_ = yaml.Unmarshal([]byte(content), &Cfg)
			}
		},
	})
	if err != nil {
		panic(err)
	}

}

func getEnvInfo(env string) bool {
	viper.AutomaticEnv()
	return viper.GetBool(env)
}
