package main
import (
	"os/signal"
	"os"
	"fmt"
	"time"
	"github.com/go-fsnotify/fsnotify"
	"study/servermonitor/util"
	"log"
	"study/servermonitor/handler"
	"container/list"
)


func main() {

	util.Md5Map = make(map[string]string)
	util.WatcherMap = make(map[string]bool) // 监听的文件夹列表
	util.ChangedMap = make(map[int]*list.List)

	desPath := "/Users/lihaoquan/GoProjects/Reference/temp"

	// 创建监听者
	var err error
	util.Watcher, err = fsnotify.NewWatcher()
	if err != nil {
		log.Fatal("创建fsnotify watcher失败……", err)
		util.Watcher.Close()
		return
	}

	//文件变更监控(异步)
	go handler.WatcherAsyncListen()

	//遍历文件夹
	walkerr := util.JumpFile(desPath)
	if walkerr != nil {
		log.Println("遍历文件夹错误……", walkerr)
		return
	}

	//重新加载所有MD5,生成新的的csv文件中
	util.OutPutToFile()
	fmt.Println("load scv file done!")

	//定时任务
	go handler.TimerCheck()

	//终止信号
	done := make(chan os.Signal, 1)
	signal.Notify(done, os.Interrupt, os.Kill)
	<-done

	util.Watcher.Close()
	util.Watcher = nil
	fmt.Println("close wather!")
	time.Sleep(1*time.Second)
	fmt.Println("服务关闭")
}
