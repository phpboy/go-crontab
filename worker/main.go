package main

import (
	"go-crontab/worker/core"
	"time"
)

func main()  {
	core.InitMongo()
	core.InitConfig()
	core.InitEtcd()
	core.InitScheduler()
	core.LoopLog()
	for{
		time.Sleep(1)
	}
}
