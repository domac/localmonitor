package util
import (
	"path/filepath"
	"os"
	"fmt"
	"io/ioutil"
	"crypto/md5"
	"io"
	"errors"
	"study/log"
	"encoding/csv"
	"time"
	"container/list"
)

//检查错误
func CheckErr(err error) {
	if err != nil {
		panic(err)
	}
}

//获取绝对路径
func GetFullPath(path string) string {
	abPath, _ := filepath.Abs(path)
	return abPath
}

//文件复制
func CopyFile(src, dst string) (written int64, err error) {
	srcFile, err := os.Open(src)
	if err!=nil {
		fmt.Println(err.Error())
		return
	}
	defer srcFile.Close()

	srcFileName := GetFileName(src)
	srcFileName = srcFileName+"_"+time.Now().String()
	dstName := filepath.Join(dst, srcFileName)
	fmt.Println("dstName==>",dstName)

	dstFile, err := os.Create(dstName)

	if err!=nil {
		fmt.Println(err.Error())
		return
	}

	defer dstFile.Close()

	return io.Copy(dstFile, srcFile)

}


func GetFileName(fullpath string) string {
	fi, err := os.Stat(fullpath)
	if err!= nil {
		return ""
	}
	return fi.Name()
}

//检查文件是否存在
func CheckFileIsExist() bool {
	_, err := os.Stat(OutputFileName)
	return err ==nil || os.IsExist(err)
}

func GenerateMd5(fileName string) string {
	file, err := os.OpenFile(fileName, os.O_RDONLY, 0644)
	defer file.Close()

	Buf, buferror := ioutil.ReadFile(fileName)
	if buferror == nil {
		fmt.Sprintf("%x", md5.Sum(Buf))
	}

	if err == nil {
		md5hash := md5.New()
		io.Copy(md5hash, file)
		return fmt.Sprintf("%x", md5hash.Sum([]byte("")))
	}else {
		fmt.Println("generate md5 fail", fileName, err)
	}
	return ""
}

//进入文件获取信息
func JumpFile(path string) error {
	abPath := GetFullPath(path)
	fmt.Println("文件顶层路径信息:", abPath)
	return filepath.Walk(abPath, func(path string, info os.FileInfo, err error) error {
		if info == nil {
			return err
		}
		if info.IsDir() {//路径为目录
			//加入监控集合中
			LetItWatcher(path)
		}else {//路径为文件
			name := info.Name()
			if OutputFileName == name {
				errors.New("不能监控输出文件...")
			}
		}
		if md5Str := GenerateMd5(path); md5Str != "" {
			Md5Map[path] = md5Str
		}
		return nil
	})
}

//输出到指定文件中
func OutPutToFile() {
	f, err := os.Create(OutputFileName)
	CheckErr(err)
	defer f.Close()

	f.WriteString("\xEF\xBB\xBF") // 写入utf8-bom
	writer := csv.NewWriter(f)

	for k, v := range Md5Map {
		isDir := false
		fileInfo, err := os.Stat(k)
		if err == nil && fileInfo.IsDir() {
			isDir = true
		}
		writer.Write([]string{FileType[isDir], k, v})
	}
	writer.Write([]string{""})
	writer.Write([]string{"文件变更历史："})
	writer.Flush()
}

//追加变更历史
func AppendChangedToOutputFile(typeStr string, fileName string, isDir bool) {
	f, err := os.OpenFile(OutputFileName, os.O_APPEND | os.O_WRONLY, os.ModeAppend)
	CheckErr(err)
	defer f.Close()

	writer := csv.NewWriter(f)
	record := []string{typeStr, FileType[isDir], fileName, time.Now().String()}
	writer.Write(record)
	writer.Flush()
}

//加入监控集合中
func LetItWatcher(path string) {
	if _, ok := WatcherMap[path]; !ok {
		WatcherMap[path] = true
		err := Watcher.Add(path)
		if err != nil {
			log.Fatal(err)
		}
	}
}

//从监控集合中删除
func DeleteItWatcher(path string) {
	if _, ok := WatcherMap[path]; ok {
		Watcher.Remove(path)
		delete(WatcherMap, path)
	}
}

func LetItChanged(typeId int, fileName string) {
	Locker.Lock()
	if _, ok := ChangedMap[typeId]; !ok {
		ChangedMap[typeId] = list.New()
	}
	_, ok := WatcherMap[fileName]
	ChangedMap[typeId].PushBack(fileName)
	AppendChangedToOutputFile(PrefixMap[typeId], fileName, ok)
	Locker.Unlock()
}