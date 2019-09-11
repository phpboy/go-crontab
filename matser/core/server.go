package core

import (
	"encoding/json"
	"fmt"
	"go-crontab/common"
	"go.etcd.io/etcd/clientv3"
	"golang.org/x/net/context"
	"net"
	"net/http"
	"strconv"
	"time"
)

func Save(h http.ResponseWriter, r *http.Request){
	var (
		key string
		value string
		resp *clientv3.PutResponse
		err error
	)

	if err=r.ParseForm();err!=nil{
		fmt.Println(err)
	}

	key=r.PostForm.Get("key")
	value=r.PostForm.Get("value")

	if resp,err=etcd.kv.Put(context.TODO(),key,value);err!=nil{
		fmt.Println(err)
	}

	fmt.Println("Revision:",resp.Header.Revision)

	if data,err:=common.ReturnJson("0","ok",string(resp.Header.Revision));err==nil{
		h.Write(data)
	}else{
		fmt.Println("error return json:",err)
	}

}

func Del(h http.ResponseWriter, r *http.Request){
	var (
		key string
		resp *clientv3.DeleteResponse
		err error
	)

	if err=r.ParseForm();err!=nil{
		fmt.Println(err)
	}

	key=r.PostForm.Get("key")

	if resp,err=etcd.kv.Delete(context.TODO(),key);err!=nil{
		fmt.Println(err)
	}

	fmt.Println("Revision:",resp.Header.Revision)

	if data,err:=common.ReturnJson("0","ok",string(resp.Header.Revision));err==nil{
		h.Write(data)
	}else{
		fmt.Println("error return json:",err)
	}

}

func List(h http.ResponseWriter, r *http.Request){
	var (
		key string
		resp *clientv3.GetResponse
		err error
		CronList []common.CronObj
		CronTab common.CronObj
	)

	if err=r.ParseForm();err!=nil{
		fmt.Println(err)
	}

	key=r.PostForm.Get("key")

	if resp,err=etcd.kv.Get(context.TODO(),key,clientv3.WithPrefix());err!=nil{
		fmt.Println(err)
	}

	for _,kp:=range resp.Kvs{

		json.Unmarshal(kp.Value,&CronTab)

		CronList = append(CronList,CronTab)
	}

	if data,err:=common.ReturnJson("0","ok",CronList);err==nil{
		h.Write(data)
	}else{
		fmt.Println("error return json:",err)
	}

}


func Kill(h http.ResponseWriter, r *http.Request){
	var (
		key string
		resp *clientv3.LeaseGrantResponse
		respPut *clientv3.PutResponse
		err error
	)

	if err=r.ParseForm();err!=nil{
		fmt.Println(err)
	}

	key=r.PostForm.Get("name")

	killerKey:="/job_kill/"+key

	resp,err = etcd.lease.Grant(context.TODO(),1)

	fmt.Println("resp.ID:",resp.ID)

	if respPut,err=etcd.kv.Put(context.TODO(),killerKey,"",clientv3.WithLease(resp.ID));err!=nil{
		fmt.Println(err)
	}

	if data,err:=common.ReturnJson("0","ok",string(respPut.Header.Revision));err==nil{
		h.Write(data)
	}else{
		fmt.Println("error return json:",err)
	}

}

func Server()  {

	mux:= http.NewServeMux()

	mux.HandleFunc("/save",Save)
	mux.HandleFunc("/del",Del)
	mux.HandleFunc("/list",List)
	mux.HandleFunc("/kill",Kill)

	server:=http.Server{
		ReadTimeout:time.Duration(Config.ReadTimeOut)*time.Second,
		WriteTimeout:time.Duration(Config.WriteTimeOut)*time.Second,
		Handler:mux,
	}

	l,err:=net.Listen("tcp",":"+strconv.Itoa(Config.Port))

	if err!=nil{
		panic(err)
	}

	err=server.Serve(l)

	if err!=nil{
		panic(err)
	}
}

