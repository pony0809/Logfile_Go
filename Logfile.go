package Logfile

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"golangProject/Compress"
)
type Logfile struct{
	Logfile *os.File
}
func LogInit()(*Logfile){
	logFile, err := os.OpenFile("logfile.txt", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatal("Cannot create log file: ", err)
	}
	
	log.SetOutput(logFile)
	log.SetPrefix("[Log日誌]")
	log.SetFlags(log.Ldate | log.Ltime)
	return &Logfile{
		Logfile: logFile,
	}
}

func (logfile *Logfile)CheckLogfile() {
	fileStat,err :=logfile.Logfile.Stat()
	if err != nil {
		fmt.Print(err)
	}

	

	fmt.Printf("\n目前log日誌大小為: %d Byte\n", fileStat.Size())

	//logfile太大就壓縮
	if fileStat.Size() > 1024*4 {
		fmt.Println("logfile 超過4MB,進行壓縮...")
		dirPath := "."
		var fileNames []os.FileInfo
		filepath.Walk(dirPath, func(path string, info os.FileInfo, err error) error {
			if err != nil {
				fmt.Println("遍歷文件夾時發生錯誤:", err)
				log.Println("遍歷文件夾時發生錯誤:", err)
				return err
			}
			if !info.IsDir() && strings.HasPrefix(path, "logfile") && strings.HasSuffix(path, ".zip") {
				fileNames = append(fileNames, info)
				_, err := os.Stat(path)
				if err != nil {
					fmt.Println("獲取文件資訊時發生錯誤:", err)
					log.Println("獲取文件資訊時發生錯誤:", err)
				}
			}
			return nil
		})

		//將檔案根據上次修改時間進行排序
		sort.Slice(fileNames, func(i, j int) bool {
			return fileNames[i].ModTime().Before(fileNames[j].ModTime())
		})

		//壓縮並清空logfile
		if len(fileNames) < 3 {
			zipName := fmt.Sprintf("logfile%d.zip", len(fileNames)+1)
			fmt.Println("將logfile壓縮成:", zipName)
			Compress.CompressZip(zipName, []string{"logfile.txt"})
			err := os.WriteFile("logfile.txt", []byte(""), 0644)
			if err != nil {
				fmt.Println("无法清空文件:", err)
				return
			}
		} else {
			Compress.CompressZip(fileNames[0].Name(), []string{"logfile.txt"})
			fmt.Println("將logfile壓縮成:", fileNames[0].Name())
			err := os.WriteFile("logfile.txt", []byte(""), 0644)
			if err != nil {
				fmt.Println("无法清空文件:", err)
				return
			}
		}
	}
}


