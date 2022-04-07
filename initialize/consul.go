package initialize

import (
	"fmt"
	consul "github.com/hashicorp/consul/api"
	"log"
	"math/rand"
)

var (
	serverId string
	err      error
)

func consulInit() {
	cfg := consul.DefaultConfig()
	cfg.Address = Cfg.ConsulConfig.Addr
	ConsulCli, err = consul.NewClient(consul.DefaultConfig())
	if err != nil {
		log.Fatalln(err)
	}
}

func ConsulRegister(addr string, port int) {
	serverId = addr + fmt.Sprintf(":%04d:%04d", port, rand.Intn(10000))

	check := &consul.AgentServiceCheck{
		GRPC:                           fmt.Sprintf("%s:%d", addr, port),
		Interval:                       "10s",
		Timeout:                        "10s",
		DeregisterCriticalServiceAfter: "10s",
	}
	register := &consul.AgentServiceRegistration{
		ID:      serverId,
		Name:    Cfg.ConsulConfig.Service.Name,
		Tags:    Cfg.ConsulConfig.Service.Tags,
		Port:    port,
		Address: addr,
		Check:   check,
	}
	if err = ConsulCli.Agent().ServiceRegister(register); err != nil {
		panic(err)
	}
}

func ConsulDeregister() {
	err = ConsulCli.Agent().ServiceDeregister(serverId)
	if err != nil {
		panic(err)
	}
}

func ConsulFilterServices(tags []string, name string) (map[string]*consul.AgentService, error) {
	filter := ""
	for _, v := range tags {
		filter += `"` + v + `" in Tags and `
	}
	filter += `Service matches "` + name + `"`
	return ConsulCli.Agent().ServicesWithFilter(filter)
}
