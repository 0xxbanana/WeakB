package Connect

import (
	"BucketFileScan/Scan"
	"bufio"
	"context"
	"fmt"
	"github.com/logrusorgru/aurora/v3"
	"github.com/tencentyun/cos-go-sdk-v5"
	"net/http"
	"net/url"
	"os"
	"strconv"
)

var SecretID, SecretKey string

type cosfile struct {
	filename string
	size     int64
}

func Tecent_input() {
	reader1 := bufio.NewReader(os.Stdin)
	_, _ = reader1.ReadString('\n')
	fmt.Println(aurora.Blue("输入SecretID"))
	fmt.Scanf("%s", &SecretID)

	reader2 := bufio.NewReader(os.Stdin)
	_, _ = reader2.ReadString('\n')
	fmt.Println(aurora.Blue("输入SecretKey"))
	fmt.Scanf("%s", &SecretKey)

	str := "SecretID: " + SecretID + "\n" + "SecretKey: " + SecretKey
	fmt.Println(aurora.BrightYellow(str))
}

func Tecent_Connect() (buckets []cos.Bucket) {
	u, _ := url.Parse("cos_url")
	b := &cos.BaseURL{BucketURL: u}
	c := cos.NewClient(b, &http.Client{
		Transport: &cos.AuthorizationTransport{
			SecretID:  SecretID,
			SecretKey: SecretKey,
		},
	})

	s, _, err := c.Service.Get(context.Background())

	if err != nil {
		panic(err)
	}

	return s.Buckets
}

// 选择存储桶，返回存储桶地址
func Tecent_Choose_Bucket(id int) (cosurl string) {
	buckets := Tecent_Connect()
	//for _, bucket := range buckets {
	//	fmt.Println("存储桶名称:" + bucket.Name + "区域:" + bucket.Region)
	//}

	//需要重新选择存储桶地址
	url := "https://" + buckets[id].Name + ".cos." + buckets[id].Region + ".myqcloud.com"
	return url
}

// 列出所有文件
func Tecent_ShowFolder(id int) (filenames []cosfile) {
	u, _ := url.Parse(Tecent_Choose_Bucket(id))
	b := &cos.BaseURL{BucketURL: u}
	client := cos.NewClient(b, &http.Client{
		Transport: &cos.AuthorizationTransport{
			SecretID:  SecretID,
			SecretKey: SecretKey,
		},
	})
	keyMarker := ""
	versionIdMarker := ""
	isTruncated := true
	opt := &cos.BucketGetObjectVersionsOptions{}
	for isTruncated {
		opt.KeyMarker = keyMarker
		opt.VersionIdMarker = versionIdMarker
		v, _, err := client.Bucket.GetObjectVersions(context.Background(), opt)
		if err != nil {
			// ERROR
			break
		}
		for _, vc := range v.Version {
			cosFile := cosfile{filename: vc.Key, size: vc.Size}
			filenames = append(filenames, cosFile)
		}
		keyMarker = v.NextKeyMarker
		versionIdMarker = v.NextVersionIdMarker
		isTruncated = v.IsTruncated
	}
	return filenames
}

func Tecent_SendScan(id int) {
	t := 0
	filenames := Tecent_ShowFolder(id)

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
