package main

import (
	"go-crontab/matser/core"
)

func main(){

	core.InitConfig()

	core.InitEtcd()

	core.Server()

}
