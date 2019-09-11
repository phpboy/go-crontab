package core

import (
	"fmt"
	"github.com/gorhill/cronexpr"
	"go-crontab/common"
	"golang.org/x/net/context"
	"os/exec"
	"time"
)

func PushToChan(obj *common.EventObj)  {
	EventChan <- obj
}

func ModifyCrontabs(eventInfo *common.EventObj) {
	var (
		exists bool
		expr *cronexpr.Expression
	)
	switch eventInfo.Type {
		case 0:
			expr,_ = cronexpr.Parse(eventInfo.CronObj.Exp)
			Crontabs[eventInfo.CronObj.Name] = &Crontab{
				Name:eventInfo.CronObj.Name,
				Cmd:eventInfo.CronObj.Cmd,
				Expr:expr,
				NextTime:time.Now(),
			}
		case 1:
			if _,exists=Crontabs[eventInfo.CronObj.Name];exists{
				delete(Crontabs,eventInfo.CronObj.Name)
			}
		case 2:
	}

	return
}

func RunCrontabs() (duration time.Duration) {

	var(
		cmd *exec.Cmd
		now time.Time
		output []byte
		err error
		nearTime *time.Time
		log *LogObj
	)

	if len(Crontabs)<1{
		duration = 1*time.Second
		return
	}

	now=time.Now()
	for _,cron:=range Crontabs{
		fmt.Println(cron.NextTime,now)
		if cron.NextTime.Before(now) || cron.NextTime.Equal(now){
			cmd = exec.CommandContext(context.TODO(),"C:\\cygwin64\\bin\\bash.exe","-c",cron.Cmd)
			if output,err = cmd.CombinedOutput();err !=nil{
				fmt.Println("err exec:",err)
				continue
			}else{
				cron.NextTime = cron.Expr.Next(time.Now())

				log = &LogObj{
					JobName:cron.Name,
					Command:cron.Cmd,
					Err:"",
					Output:string(output),
					PlanTime:cron.NextTime.UnixNano()/1000/1000,
				}
				fmt.Println("push to log ",log)
				PushToLogChan(log)
			}

		}
		if nearTime == nil || cron.NextTime.Before(*nearTime){
			nearTime = &cron.NextTime
		}
	}

	duration = (*nearTime).Sub(now)
	return
}
