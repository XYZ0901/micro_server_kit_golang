package initialize

import (
	sentinel "github.com/alibaba/sentinel-golang/api"
	"github.com/alibaba/sentinel-golang/core/circuitbreaker"
	"github.com/alibaba/sentinel-golang/core/config"
	"github.com/alibaba/sentinel-golang/core/flow"
	"github.com/alibaba/sentinel-golang/logging"
	"log"
)

type stateChangeListener struct {
}

func (s *stateChangeListener) OnTransformToClosed(prev circuitbreaker.State, rule circuitbreaker.Rule) {
	//TODO implement me
	//panic("implement me")
}

func (s *stateChangeListener) OnTransformToOpen(prev circuitbreaker.State, rule circuitbreaker.Rule, snapshot interface{}) {
	//TODO implement me
	//panic("implement me")
}

func (s *stateChangeListener) OnTransformToHalfOpen(prev circuitbreaker.State, rule circuitbreaker.Rule) {
	//TODO implement me
	//panic("implement me")
}

func sentinelInit() {
	conf := config.NewDefaultConfig()
	conf.Sentinel.Log.Logger = logging.NewConsoleLogger()
	err := sentinel.InitWithConfig(conf)
	if err != nil {
		log.Fatalln(err)
	}
	// https://sentinelguard.io/zh-cn/docs/golang/flow-control.html
	_, err = flow.LoadRules([]*flow.Rule{
		{
			Resource:               "qps_Reject",
			TokenCalculateStrategy: flow.Direct,
			ControlBehavior:        flow.Reject,
			Threshold:              100,
			StatIntervalInMs:       1000,
		},
		{
			Resource:               "qps_Throttling",
			TokenCalculateStrategy: flow.Direct,
			ControlBehavior:        flow.Throttling,
			MaxQueueingTimeMs:      5000,
			Threshold:              10,
			StatIntervalInMs:       1000,
		},
		{
			Resource:               "warm_Reject",
			TokenCalculateStrategy: flow.WarmUp,
			ControlBehavior:        flow.Reject,
			Threshold:              1000,
			WarmUpPeriodSec:        10,
			WarmUpColdFactor:       3,
			StatIntervalInMs:       1000,
		},
	})
	if err != nil {
		log.Fatalf("Unexpected error: %+v\n", err)
		return
	}

	circuitbreaker.RegisterStateChangeListeners(&stateChangeListener{})
	_, err = circuitbreaker.LoadRules([]*circuitbreaker.Rule{
		{
			Resource:                     "slow_rate",
			Strategy:                     circuitbreaker.SlowRequestRatio,
			RetryTimeoutMs:               3000,
			MinRequestAmount:             10,
			StatIntervalMs:               5000,
			StatSlidingWindowBucketCount: 10,
			MaxAllowedRtMs:               50,
			Threshold:                    0.5,
		},
		{
			Resource:                     "err_rate",
			Strategy:                     circuitbreaker.ErrorRatio,
			RetryTimeoutMs:               3000,
			MinRequestAmount:             10,
			StatIntervalMs:               5000,
			StatSlidingWindowBucketCount: 10,
			Threshold:                    0.4,
		},
		{
			Resource:                     "err_count",
			Strategy:                     circuitbreaker.ErrorCount,
			RetryTimeoutMs:               3000,
			MinRequestAmount:             10,
			StatIntervalMs:               5000,
			StatSlidingWindowBucketCount: 10,
			Threshold:                    50,
		},
	})
	if err != nil {
		log.Fatalln(err)
		return
	}
}
