package util
import (
	"container/list"
	"github.com/go-fsnotify/fsnotify"
	"sync"
)

var (

	BAK_PATH string = "/Users/lihaoquan/GoProjects/Reference/temp_bak"

	Locker sync.Mutex

	OutputFileName string = "filesName.csv"

	Md5Map map[string]string

	ChangedMap map[int]*list.List

	WatcherMap map[string]bool

	Watcher *fsnotify.Watcher

	PrefixMap = map[int]string{
		1: "修改",
		2: "新建",
		3: "删除",
		4: "重命名",
	}

	FileType = map[bool]string{
		true:  "文件夹",
		false: "文件",
	}
)