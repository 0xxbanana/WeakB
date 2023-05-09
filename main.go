package main

import (
	"BucketFileScan/Connect"
	"bufio"
	"fmt"
	"github.com/logrusorgru/aurora/v3"
	"os"
)

func main() {

	banner := `
██╗    ██╗███████╗ █████╗ ██╗  ██╗██████╗ 
██║    ██║██╔════╝██╔══██╗██║ ██╔╝██╔══██╗
██║ █╗ ██║█████╗  ███████║█████╔╝ ██████╔╝
██║███╗██║██╔══╝  ██╔══██║██╔═██╗ ██╔══██╗
╚███╔███╔╝███████╗██║  ██║██║  ██╗██████╔╝
 ╚══╝╚══╝ ╚══════╝╚═╝  ╚═╝╚═╝  ╚═╝╚═════╝ 

`
	aurora := aurora.NewAurora(true) // 强制启用彩色输出

	print(banner)
	ChooseCloud := 0
	fmt.Println(aurora.Red("欢迎使用由0xxbanana设计的存储桶扫描器，第一次写项目有很多不足的地方请在issue提交\n"))
	fmt.Println(aurora.Blue("请输入数字所对应要选择的云服务商"))
	fmt.Println(aurora.Green(
		"1 Alibaba Cloud" + "\n" +
			"2 Huawei Cloud" + "\n" +
			"3 Tencent Cloud"))
	fmt.Printf("%s", aurora.BrightWhite("==>"))
	fmt.Scanf("%d", &ChooseCloud)
	if ChooseCloud >= 1 && ChooseCloud <= 3 {
		switch ChooseCloud {
		case 1:
			Connect.Ali_input()
			Ali()
			break

		case 2:
			Connect.H_input()
			Huawei()
			break

		case 3:
			Connect.Tecent_input()
			Tecent()
			break

		}
	} else {
		fmt.Println(aurora.Red("输入有误"))
	}
	fmt.Println(aurora.Red("\n已完成，感谢使用"))
}

func Ali() {
	var id int

	for i, bucket := range Connect.Ali_Connect_Aliyun() {
		fmt.Printf("%d %s\n", aurora.Cyan(i), aurora.Green(bucket.Bucketname))
	}

	reader := bufio.NewReader(os.Stdin)
	_, _ = reader.ReadString('\n')
	fmt.Printf("%s", aurora.Blue("请输入存储桶序号"))
	fmt.Printf("\n%s", aurora.BrightWhite("==>"))
	fmt.Scanf("%d", &id)
	if id < 0 || id >= len(Connect.Ali_Connect_Aliyun()) {
		fmt.Println(aurora.Red("输入存储桶序号有误"))
	}
	Connect.Ali_SendScan(id)
}

func Huawei() {
	var id int

	for i, bucket := range Connect.H_Connect_Huaweiyun() {
		fmt.Printf("%d %s\n", aurora.Cyan(i), aurora.Green(bucket.Buketname))
	}

	reader := bufio.NewReader(os.Stdin)
	_, _ = reader.ReadString('\n')
	fmt.Printf("%s", aurora.Blue("请输入存储桶序号"))
	fmt.Printf("\n%s", aurora.BrightWhite("==>"))
	fmt.Scanf("%d", &id)
	if id < 0 || id >= len(Connect.H_Connect_Huaweiyun()) {
		fmt.Println(aurora.Red("输入存储桶序号有误"))
	}
	Connect.H_SendScan(id)
}

func Tecent() {
	var id int

	for i, bucket := range Connect.Tecent_Connect() {
		fmt.Printf("%d %s\n", aurora.Cyan(i), aurora.Green(bucket.Name))
	}

	reader := bufio.NewReader(os.Stdin)
	_, _ = reader.ReadString('\n')
	fmt.Printf("%s", aurora.Blue("请输入存储桶序号"))
	fmt.Printf("\n%s", aurora.BrightWhite("==>"))
	fmt.Scanf("%d", &id)

	if id < 0 || id >= len(Connect.Tecent_Connect()) {
		fmt.Println(aurora.Red("输入存储桶序号有误"))
	}

	Connect.Tecent_SendScan(id)
}
