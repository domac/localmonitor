package handler
import (
	"time"
	"study/servermonitor/util"
	"container/list"
	"github.com/go-fsnotify/fsnotify"
	"os"
	"fmt"
	"log"
)

//异步监听
func WatcherAsyncListen() {
	for {
		select {
		case event := <-util.Watcher.Events:
			log.Println("event:", event)
			if event.Op&fsnotify.Write == fsnotify.Write { //修改文件

				fmt.Println(">>",event.Name,"[edit]")

				oldMd5, ok := util.Md5Map[event.Name]
				if ok {
					newMd5Str := util.GenerateMd5(event.Name)
					if newMd5Str != oldMd5 {
						util.Md5Map[event.Name] = newMd5Str
						util.LetItChanged(1, event.Name)
					}

				}else {
					util.Md5Map[event.Name] = util.GenerateMd5(event.Name)
					util.LetItChanged(1, event.Name)
				}

			}else if event.Op&fsnotify.Create == fsnotify.Create {//创建文件

				fmt.Println(">>",event.Name,"[create]")

				fileInfo, err := os.Stat(event.Name)
				if err == nil && fileInfo.IsDir() {
					util.LetItWatcher(event.Name)
				}

				newMd5Str := util.GenerateMd5(event.Name)
				if len(newMd5Str) > 0 {
					util.Md5Map[event.Name] = newMd5Str
					util.LetItChanged(2, event.Name)
				}

			} else if event.Op&fsnotify.Remove == fsnotify.Remove {//移除文件

				fmt.Println(">>",event.Name,"[remove]")

				if _, ok := util.Md5Map[event.Name]; ok {
					delete(util.Md5Map, event.Name)
					util.LetItChanged(3, event.Name)
				}
			}else if event.Op&fsnotify.Rename == fsnotify.Rename {//重命名文件

				fmt.Println(">>",event.Name,"[rename]")

				if _, ok := util.Md5Map[event.Name]; ok {
					delete(util.Md5Map, event.Name)
					util.LetItChanged(4, event.Name)
				}

				util.DeleteItWatcher(event.Name)
			}
		case err := <-util.Watcher.Errors:
		//监听异常
			if err != nil {
				log.Println("error:", err)
			}
			return
		}
	}
}

//定时检测
func TimerCheck() {
	timer := time.NewTicker(10 * time.Second)
	for {
		select {
		case <-timer.C:
			util.Locker.Lock()
			content := ""
			for k, v := range util.ChangedMap {
				content += util.PrefixMap[k] + "列表: <br>"
				for el := v.Front(); el != nil; el = el.Next() {
					content+=el.Value.(string) + "<br>"
				}
				content += "<br>"
			}
			if len(content) >0 {
				util.ChangedMap = make(map[int]*list.List)
				go SendMail(content)//发送邮件
			}
			util.Locker.Unlock()
		}
	}
}


