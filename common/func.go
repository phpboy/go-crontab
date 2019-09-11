package common

import "encoding/json"

type CronObj struct {
	Name string `json:"name"`
	Cmd string `json:"cmd"`
	Exp string `json:"exp"`
}

type EventObj struct {
	CronObj CronObj
	Type int32
}

type Response struct {
	Code string `json:"code"`
	Msg string	`json:"msg"`
	Data interface{}	`json:"data"`
}

func ReturnJson(code,msg string,data interface{}) (jsonByte []byte,err error){

	bytes:=Response{
		Code:code,
		Msg:msg,
		Data:data,
	}

	jsonByte,err=json.Marshal(bytes)

	return

}


func BuildEvent(Type int32,val []byte) (E *EventObj){

	jsonObj := CronObj{}

	json.Unmarshal(val,&jsonObj)

	E = &EventObj{
		CronObj:jsonObj,
		Type:Type,
	}

	return

}
