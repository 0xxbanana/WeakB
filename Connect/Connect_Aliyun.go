package Connect

import (
	"BucketFileScan/Scan"
	"bufio"
	"fmt"
	"github.com/aliyun/aliyun-oss-go-sdk/oss"
	"github.com/logrusorgru/aurora/v3"
	"os"
	"strconv"
)

var AccessKeyID, AccessKeySecret string

type AliBuckets struct {
	Bucketname string
	Location   string
}

type AliFiles struct {
	filename string
	size     int64
}

func Ali_input() {
	reader1 := bufio.NewReader(os.Stdin)
	_, _ = reader1.ReadString('\n')
	fmt.Println(aurora.Blue("输入AccessKeyID"))
	fmt.Scanf("%s", &AccessKeyID)

	reader2 := bufio.NewReader(os.Stdin)
	_, _ = reader2.ReadString('\n')
	fmt.Println(aurora.Blue("输入AccessKeySecret"))
	fmt.Scanf("%s", &AccessKeySecret)

	str := "AccessKeyID: " + AccessKeyID + "\n" + "AccessKeySecret: " + AccessKeySecret
	fmt.Println(aurora.BrightYellow(str))
}

// 列出所有存储空间
func Ali_Connect_Aliyun() (bucketnames []AliBuckets) {
	var client, err = oss.New("oss-cn-hangzhou.aliyuncs.com", AccessKeyID, AccessKeySecret)
	if err != nil {
		fmt.Println("Error:", err)
		os.Exit(-1)
	}

	// 列举当前账号所有地域下的存储空间。
	marker := ""
	for {
		lsRes, err := client.ListBuckets(oss.Marker(marker))
		if err != nil {
			fmt.Println("Error:", err)
			os.Exit(-1)
		}

		// 默认情况下一次返回100条记录。
		for _, bucket := range lsRes.Buckets {
			Alifile := AliBuckets{Bucketname: bucket.Name, Location: bucket.Location}
			bucketnames = append(bucketnames, Alifile)
		}

		if lsRes.IsTruncated {
			marker = lsRes.NextMarker
		} else {
			break
		}
	}
	return bucketnames
}

// 列出所有文件
func HandleError(err error) {
	fmt.Println("Error:", err)
	os.Exit(-1)
}

func Ali_Choose_Ali_Bucket(id int) (filenames []AliFiles) {
	AliBucketnames := Ali_Connect_Aliyun()
	bucketName := AliBucketnames[id].Bucketname
	var client, err = oss.New(AliBucketnames[id].Location+".aliyuncs.com", AccessKeyID, AccessKeySecret)
	bucket, err := client.Bucket(bucketName)
	if err != nil {
		HandleError(err)
	}
	// 列举所有文件。
	marker := ""
	for {
		lsRes, err := bucket.ListObjects(oss.Marker(marker))
		if err != nil {
			HandleError(err)
		}
		// 打印列举结果。默认情况下，一次返回100条记录。
		for _, object := range lsRes.Objects {
			Alifilekey := AliFiles{filename: object.Key, size: object.Size}
			filenames = append(filenames, Alifilekey)
		}
		if lsRes.IsTruncated {
			marker = lsRes.NextMarker
		} else {
			break
		}
	}
	return filenames
}

func Ali_SendScan(id int) {
	t := 0
	filenames := Ali_Choose_Ali_Bucket(id)

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
