package Scan

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

func FindWeakFile(filename string) (findresult bool) {
	flag := false
	ext := ""
	dotIndex := len(filename) - 1
	for i := dotIndex; i >= 0; i-- {
		if filename[i] == '.' {
			ext = filename[i:]
			break
		}
	}
	// 创建两个channel
	dictChan := make(chan string) // 用于将字典文件中的后缀名传递给比对任务
	resultChan := make(chan bool) // 用于将比对结果传递给主任务

	// 创建读取字典文件的goroutine
	go func() {
		// 加载字典文件
		dictFile, err := os.Open("Dict/WeakSuffixName.txt")
		if err != nil {
			fmt.Println("打开字典文件失败：", err)
			return
		}
		defer dictFile.Close()

		// 逐行读取字典文件
		scanner := bufio.NewScanner(dictFile)
		for scanner.Scan() {
			// 将字典文件中的后缀名发送到channel中
			dictChan <- strings.TrimSpace(scanner.Text())
		}

		if err := scanner.Err(); err != nil {
			fmt.Println("读取字典文件出错：", err)
			return
		}

		// 关闭channel
		close(dictChan)
	}()

	// 创建比对任务的goroutine
	go func() {
		// 循环从channel中读取字典文件中的后缀名，直到channel关闭
		for suffix := range dictChan {
			// 对比文件后缀名和字典文件中的后缀名
			if strings.ToLower(ext) == strings.ToLower(suffix) {
				// 将比对结果发送到channel中
				resultChan <- true
				return
			}
		}

		// 没有匹配项，将比对结果发送到channel中
		resultChan <- false
	}()

	flag = <-resultChan
	return flag
}
