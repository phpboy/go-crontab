package core

import (
	"fmt"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"golang.org/x/net/context"
	"time"
)
var LogChan chan *LogObj
var MongoClient *Mongo

type Mongo struct{
	Client *mongo.Client
	Collection *mongo.Collection
}

type Log struct {
	Logs []interface{}
}

type LogObj struct {
	JobName string `json:"jobName" bson:"jobName"` // 任务名字
	Command string `json:"command" bson:"command"` // 脚本命令
	Err string `json:"err" bson:"err"` // 错误原因
	Output string `json:"output" bson:"output"`	// 脚本输出
	PlanTime int64 `json:"planTime" bson:"planTime"` // 计划开始时间
	ScheduleTime int64 `json:"scheduleTime" bson:"scheduleTime"` // 实际调度时间
	StartTime int64 `json:"startTime" bson:"startTime"` // 任务执行开始时间
	EndTime int64 `json:"endTime" bson:"endTime"` // 任务执行结束时间
}

func InitMongo()  {

	LogChan = make(chan *LogObj,100)

	var (
		client *mongo.Client
		collection *mongo.Collection
		err error
	)

	client,err = mongo.NewClient(options.Client().ApplyURI("mongodb://localhost:27017"))

	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)

	if err = client.Connect(ctx);err!=nil{
		panic(err)
	}

	collection = client.Database("cron").Collection("logs")

	MongoClient = &Mongo{
		Client:client,
		Collection:collection,
	}
}

func LoopLog()  {

	fmt.Println(111111)
	var(
		log *LogObj
		logs Log
	)

	go func() {
		for{
			fmt.Println(2222)
			select {
				case log = <-LogChan:
					logs.Logs = append(logs.Logs,log)
				}

				if len(logs.Logs)>=1{
					SaveLogs(logs)
				}
		}
	}()

}

func SaveLogs(logs Log)  {
	var(
		err error
	)
	if _,err = MongoClient.Collection.InsertMany(context.TODO(),logs.Logs);err!=nil{
		fmt.Println(err)
	}
}

func PushToLogChan(log *LogObj)  {
	LogChan <- log
}