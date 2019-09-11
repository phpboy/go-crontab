package core

import (
	"fmt"
	"github.com/gorhill/cronexpr"
	"go-crontab/common"
	"go.etcd.io/etcd/clientv3"
	"golang.org/x/net/context"
	"time"
)

var EventChan chan *common.EventObj
var Crontabs map[string]*Crontab

type Crontab struct {
	Name string
	Cmd string
	Expr *cronexpr.Expression
	NextTime time.Time
}

func InitScheduler(){
	EventChan = make(chan *common.EventObj,10)
	Crontabs = make(map[string]*Crontab)
	WatchEtcd()
	Loop()
}

func WatchEtcd(){

	var (
		resp *clientv3.GetResponse
		err error
		revStart int64
		watchChan clientv3.WatchChan
		watchResp clientv3.WatchResponse
		watchEvent *clientv3.Event
	)

	if resp,err = etcd.kv.Get(context.TODO(),"/job/",clientv3.WithPrefix());err!=nil{
		fmt.Println(err)
	}

	revStart = resp.Header.Revision+1

	for _,v:=range resp.Kvs{
		fmt.Println(string(v.Key),string(v.Value))
		PushToChan(common.BuildEvent(0,[]byte(v.Value)))
	}

	fmt.Println(revStart)

	go func() {

		watchChan = etcd.watcher.Watch(context.TODO(),"/job/",clientv3.WithRev(revStart),clientv3.WithPrefix())

		for watchResp = range watchChan{

			for _,watchEvent = range watchResp.Events{
				switch watchEvent.Type {
				case 0:
					PushToChan(common.BuildEvent(0,watchEvent.Kv.Value))
				case 1:
					PushToChan(common.BuildEvent(1,watchEvent.Kv.Value))
				}

			}
		}

	}()

	return
}

func Loop()  {

	var(
		Event *common.EventObj
		timer *time.Timer
		duration time.Duration
	)

	duration = ExecAll()

	timer = time.NewTimer(duration)

	for {
		select {
		case Event=<-EventChan:
			Execute(Event)
		case <-timer.C:
		}

		duration = ExecAll()

		timer.Reset(duration)
	}

}