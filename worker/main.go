package main

import "go-crontab/worker/core"

func main()  {
	core.InitConfig()
	core.InitEtcd()
	core.InitScheduler()
}
