package Connect

import (
	"BucketFileScan/Scan"
	"bufio"
	"fmt"
	"github.com/huaweicloud/huaweicloud-sdk-go-obs/obs"
	"github.com/logrusorgru/aurora/v3"
	"os"
	"strconv"
)

var Access_Key_Id, Secret_Access_Key string

var endPoint = "https://obs.cn-north-4.myhuaweicloud.com"

type HuaweiBucket struct {
	Buketname string
	Location  string
}

type HuaweiFile struct {
	filename string
	size     int64
}

func H_input() {
	reader1 := bufio.NewReader(os.Stdin)
	_, _ = reader1.ReadString('\n')
	fmt.Println(aurora.Blue("输入Access_Key_Id"))
	fmt.Scanf("%s", &Access_Key_Id)

	reader2 := bufio.NewReader(os.Stdin)
	_, _ = reader2.ReadString('\n')
	fmt.Println(aurora.Blue("输入Secret_Access_Key"))
	fmt.Scanf("%s", &Secret_Access_Key)

	str := "Access_Key_Id: " + Access_Key_Id + "\n" + "Secret_Access_Key: " + Secret_Access_Key
	fmt.Println(aurora.BrightYellow(str))
}

// 创建ObsClient结构体
//var obsClient, _ = obs.New(Access_Key_Id, Secret_Access_Key, endPoint)

func H_Connect_Huaweiyun() (huaweibuckets []HuaweiBucket) {
	var obsClient, _ = obs.New(Access_Key_Id, Secret_Access_Key, endPoint)

	var vallocation string
	output, err := obsClient.ListBuckets(nil)
	if err == nil {
		for _, val := range output.Buckets {
			vallocation = H_GetLocation(val.Name)
			huaweibucket := HuaweiBucket{Buketname: val.Name, Location: vallocation}
			huaweibuckets = append(huaweibuckets, huaweibucket)
		}
	} else {
		if obsError, ok := err.(obs.ObsError); ok {
			fmt.Println(obsError.Code)
			fmt.Println(obsError.Message)
		} else {
			fmt.Println(err)
		}
	}
	return huaweibuckets
}

func H_GetLocation(bucket string) (location string) {
	var obsClient, _ = obs.New(Access_Key_Id, Secret_Access_Key, endPoint)

	var relocation string
	output, err := obsClient.GetBucketLocation(bucket)
	if err == nil {
		relocation = output.Location
	} else {
		if obsError, ok := err.(obs.ObsError); ok {
			fmt.Println(obsError.Code)
			fmt.Println(obsError.Message)
		} else {
			fmt.Println(err)
		}
	}
	return relocation
}

func H_Choose_Bucket(id int) (huaweifiles []HuaweiFile) {
	var obsClient, _ = obs.New(Access_Key_Id, Secret_Access_Key, endPoint)

	huaweibuckets := H_Connect_Huaweiyun()
	//fmt.Println(huaweibuckets[0].buketname)
	input := &obs.ListObjectsInput{}
	input.Bucket = huaweibuckets[id].Buketname
	endPoint = "https://obs." + huaweibuckets[0].Location + ".myhuaweicloud.com"
	obsClient, _ = obs.New(Access_Key_Id, Secret_Access_Key, endPoint)
	output, err := obsClient.ListObjects(input)
	if err == nil {
		for _, val := range output.Contents {
			//fmt.Printf("Key:%s  Size:%d\n", val.Key, val.Size)
			huaweifile := HuaweiFile{filename: val.Key, size: val.Size}
			huaweifiles = append(huaweifiles, huaweifile)

		}
	} else {
		if obsError, ok := err.(obs.ObsError); ok {
			fmt.Println(obsError.Code)
			fmt.Println(obsError.Message)
		} else {
			fmt.Println(err)
		}
	}
	return huaweifiles
}

func H_SendScan(id int) {
	t := 0
	filenames := H_Choose_Bucket(id)
	for i := range filenames {
		if Scan.FindWeakFile(filenames[i].filename) == true && filenames[i].size != 0 {
			fmt.Println(aurora.Magenta("[INFO] " + strconv.Itoa(t) + " 存在疑似敏感文件" + filenames[i].filename))
			t++
		}
	}

	if t == 0 {
		fmt.Println(aurora.Magenta("不存在敏感文件"))
	}
}
