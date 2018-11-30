package main

//多线程拷贝文件
//默认10线程并行
//去除表中空行数据
//author: paavo

import (
	"github.com/360EntSecGroup-Skylar/excelize"
	"fmt"
	"os"
	"strings"
	"log"
	"io"
)

const count = 10
const pwdpath = "/app/content/http/data"
func main() {
	nowPath, _ := os.Getwd()
	if len(os.Args) != 3 || nowPath != pwdpath {
		log.Println("Usage：filename  distpath  \nPlease into ",pwdpath, " directory")
		return
	}
	//if len(os.Args) != 3 {
	//	log.Println("Usage：filename  distpath  \nPlease into ", " directory")
	//	return
	//}
	filename := os.Args[1]
	dispath := os.Args[2]
	_, e := os.Stat(filename)
	if e != nil {
		fmt.Println("文件不存在.....")
		os.Exit(1)
	}

	xlsx, err := excelize.OpenFile(filename)
	if err != nil {
		fmt.Println(err)
		return
	}
	sheetName := xlsx.GetSheetName(1)
	fmt.Println(sheetName)

	//获取表格中的数据
	rows := xlsx.GetRows(sheetName)
	rowsCount := DeleteStrip(rows, 0)
	fmt.Println("非空行:",rowsCount)

	fileChan := make(chan []string, rowsCount)
	endChan := make(chan bool, rowsCount)
	for _, row := range rows {
		if row[0] != "" {
			str := strings.Split(row[0], "/")
			str = str[1:]
			fileChan <- str
		}

	}

	//开启十个协程进行操作
	for i := 0; i < count; i++ {
		go CopyFile(dispath, fileChan, endChan)
	}

	for i := 0; i < rowsCount; i++ {
		<-endChan
	}

}

func CopyFile(distpath string, filechan chan []string, endchan chan bool) {
	for {

		fileSlice := <-filechan
		FilePath := strings.Join(fileSlice[0:len(fileSlice)-1], "/")
		srcName := fileSlice[len(fileSlice)-1]
		dstName := srcName
		allPath := distpath + "/" + FilePath

		//create directory
		_, err := os.Stat(allPath)
		if err != nil {
			os.MkdirAll(allPath, 0755)
		}
		//fmt.Println("src:", FilePath+"/"+srcName, "\tdist:", allPath+"/"+dstName)

		src, err := os.Open(FilePath + "/" + srcName)
		if err != nil {
			continue
		}
		defer src.Close()
		dst, err := os.OpenFile(allPath+"/"+dstName, os.O_WRONLY|os.O_CREATE, 0744)
		if err != nil {
			continue
		}
		defer dst.Close()
		_, r := io.Copy(dst, src)
		if r != nil {
			continue
		}
		endchan <- true
	}

}

func DeleteStrip(rows [][]string, index int) int {
	Counts := 0
	for _, row := range rows {
		if row[index] == "" {
			continue
		} else {
			Counts++
		}
	}
	return Counts
}
