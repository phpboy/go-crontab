package core

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"os"
)

var Config *C


type C struct {
	ReadTimeOut,WriteTimeOut,Port int
	Etcd string
}

var conf = flag.String("config","D:/code/go/src/go-crontab/matser/config.json","config filename")
func InitConfig() {

	flag.Parse()

	if *conf == ""{
		fmt.Println("config filename required")
		return
	}

	var (
		file * os.File
		err error
		content []byte
	)

	if file,err=os.Open(*conf);err!=nil{
		fmt.Println(err)
		return
	}

	if content,err=ioutil.ReadAll(file);err!=nil{
		fmt.Println(err)
		return
	}

	err = json.Unmarshal(content,&Config)

}