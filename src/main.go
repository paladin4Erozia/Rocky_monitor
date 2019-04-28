package main

import (
	"bytes"
	"fmt"
	"github.com/chenjiandongx/go-echarts/charts"
	_ "github.com/go-sql-driver/mysql"
	"io"
	"io/ioutil"
	"log"
	"mime/multipart"
	"net/http"
	"orm"
	"os"
	"strconv"
	"strings"
	"time"
)

func main(){
    //go uploadfile("d:/photo/cut.png", "http://localhost:9090/upload")
	//webser.Start()
	file_proc("./files/Sysinfo.txt")
	//visual()
}

func file_proc(filePth string){
	if fileObj, err := os.Open(filePth); err == nil {
		defer fileObj.Close()
		if contents, err := ioutil.ReadAll(fileObj); err == nil {
			result := strings.Replace(string(contents), "\n", "", 1)
			info := map[string]float64{}

			start := strings.Index(result, "CPU logical num: ")+strings.Count("CPU logical num: ","")-1
			end := strings.Index(result, "CPU physical cores num:")
			info["cpuLogicalNum"], err = strconv.ParseFloat(result[start:end-1], 32/64)

			start = strings.Index(result, "CPU physical cores num: ")+strings.Count("CPU physical cores num: ","")-1
			end = strings.Index(result, "scputimes")
			info["cpuPhysicalNum"], err = strconv.ParseFloat(result[start:end-2], 32/64)

			start = strings.Index(result, "scputimes(user=")+strings.Count("scputimes(user=","")-1
			end = strings.Index(result, ", system=")
			info["cpuUserTimeUsing"], err = strconv.ParseFloat(result[start:end], 32/64)

			start = strings.Index(result, ", system=")+strings.Count(", system=","")-1
			end = strings.Index(result, ", idle=")
			info["cpuSystemTimeUsing"], err = strconv.ParseFloat(result[start:end], 32/64)

			start = strings.Index(result, ", idle=")+strings.Count(", idle=","")-1
			end = strings.Index(result, ", interrupt=")
			info["cpuIdleTime"], err = strconv.ParseFloat(result[start:end], 32/64)

			start = strings.Index(result, "CPU using in near 10 sec:")+strings.Count("CPU using in near 10 sec:","")-1
			end = strings.Index(result, "svmem(total=")
			cpuUsingInfo := result[start+2:end-2]

			var cpuUsingPer [64] float64
			var tem, now uint8
			now = 0
			for i := 0; i < len(cpuUsingInfo); i++{
				if uint8(cpuUsingInfo[i]) >= uint8('0') && uint8(cpuUsingInfo[i] )<= uint8('9'){
					tem = uint8(cpuUsingInfo[i]) - uint8('0')
				}else if cpuUsingInfo[i] == '.'{
					tem = tem * 10 + uint8(cpuUsingInfo[i+1]) - uint8('0')
					cpuUsingPer[now] = cpuUsingPer[now] + float64(tem)
					now = now + 1
					if now >= uint8(info["cpuLogicalNum"]){
						now = 0
					}
				}
 			}

			for i := 0; i < int(info["cpuLogicalNum"]); i++{
				cpuUsingPer[i] = cpuUsingPer[i] / 10.0
				info["cpuPhysicalNum" + string(i + int('0'))] = cpuUsingPer[i]
			}

			start = strings.Index(result, "svmem(total=")+strings.Count("svmem(total=","")-1
			end = strings.Index(result, ", available=")
			info["memTotal"], err = strconv.ParseFloat(result[start:end], 32/64)

			start = strings.Index(result, " available=")+strings.Count(" available=","")-1
			end = strings.Index(result, ", percent=")
			info["memAvailable"], err = strconv.ParseFloat(result[start:end], 32/64)

			start = strings.Index(result, "sdiskusage(total=")+strings.Count("sdiskusage(total=","")-1
			for i := start; i<len(result); i++{
				if uint8(result[i]) < uint8('0') || uint8(result[i]) > uint8('9'){
					end = i
					break
				}
			}
			info["diskTotal"], err = strconv.ParseFloat(result[start:end], 32/64)

			start = end+strings.Count(", used=","")-1
			for i := start; i<len(result); i++{
				if uint8(result[i]) < uint8('0') || uint8(result[i]) > uint8('9'){
					end = i
					break
				}
			}
			info["diskUsed"], err = strconv.ParseFloat(result[start:end], 32/64)

			start = strings.Index(result, "snetio(bytes_sent=")+strings.Count("snetio(bytes_sent=","")-1
			end = strings.Index(result, ", bytes_recv=")
			info["snetioBytesSent"], err = strconv.ParseFloat(result[start:end], 32/64)

			start = strings.Index(result, ", bytes_recv=")+strings.Count(", bytes_recv=","")-1
			end = strings.Index(result, ", packets_sent=")
			info["snetioBytesRecv"], err = strconv.ParseFloat(result[start:end], 32/64)

			start = strings.Index(result, ", packets_sent=")+strings.Count(", packets_sent=","")-1
			end = strings.Index(result, ", packets_recv=")
			info["snetioPacketsSent"], err = strconv.ParseFloat(result[start:end], 32/64)

			start = strings.Index(result, ", packets_recv=")+strings.Count(", packets_recv=","")-1
			end = strings.Index(result, ", errin=")
			info["snetioPacketsRecv"], err = strconv.ParseFloat(result[start:end], 32/64)

			visual(info)

		}
	} else {
		print(err)
	}
}

func cpuNumInfoVisaul(info map[string]float64) *charts.Bar {
	var nameItems = []string{"Cpu Logical Num", "CPU Physical Num"}
	cpuNum := charts.NewBar()
	cpuNum.SetGlobalOptions(charts.TitleOpts{Title: "CPU NUM"}, charts.ToolboxOpts{Show: true})
	cpuNum.AddXAxis(nameItems).
		AddYAxis("Rocky's PC", [2]float64{info["cpuLogicalNum"], info["cpuPhysicalNum"]})
	return cpuNum
}

func cpuTimeInfoVisaul(info map[string]float64) *charts.Pie {
	var nameItems = []string{"User Time Using", "System Time Using", "Idle Time"}
	pie := charts.NewPie()
	pie.SetGlobalOptions(charts.TitleOpts{Title: "CPU Time Using"})
	cpuTimeUsingInfo := make(map[string]interface{})
	cpuTimeUsingInfo[nameItems[0]] = info["cpuUserTimeUsing"]
	cpuTimeUsingInfo[nameItems[1]] = info["cpuSystemTimeUsing"]
	cpuTimeUsingInfo[nameItems[2]] = info["cpuIdleTime"]
	pie.Add("area", cpuTimeUsingInfo,
		charts.PieOpts{Radius: []string{"30%", "75%"}, RoseType: "area", Center: []string{"25%", "50%"}},
	)
	pie.Add("radius", cpuTimeUsingInfo,
		charts.LabelTextOpts{Show: true, Formatter: "{b}: {c}"},
		charts.PieOpts{Radius: []string{"30%", "75%"}, RoseType: "radius", Center: []string{"75%", "50%"}},
	)
	return pie
}

func visual(info map[string]float64) {
	page := charts.NewPage()
	page.Add(
		cpuNumInfoVisaul(info),
		cpuTimeInfoVisaul(info),
	)
	f, err := os.Create("visaul.html")
	if err != nil {
		log.Println(err)
	}
	page.Render(f)
}

func uploadfile(filename string, url string){
    time.Sleep(30 * time.Second)
	bodyBuf := &bytes.Buffer{}
	bodyWriter :=multipart.NewWriter(bodyBuf)
	//模拟创建form表单字段
	strs := strings.Split(filename, "/")
	destname := strs[len(strs) - 1]
	filewriter, err := bodyWriter.CreateFormFile("uploadfile", destname)
	if err != nil{
		fmt.Println("error writing to buffer")
		return
	}
    //打开文件句柄操作
	fh, err := os.Open(filename)
	if err != nil{
		fmt.Println("error open file")
		return
	}
	defer fh.Close()

	//拷贝文件
	_, err = io.Copy(filewriter, fh)
	if err != nil{
		fmt.Println("error copy file")
		return
	}
	contentType := bodyWriter.FormDataContentType()
	bodyWriter.Close()

	resp, err := http.Post(url, contentType, bodyBuf)
	if err != nil{
		fmt.Println("error post buffer")
		return
	}
	defer resp.Body.Close()
	resp_body, err := ioutil.ReadAll(resp.Body)
	if err != nil{
		fmt.Println("error read all")
		return
	}
	fmt.Println(resp.Status)
	fmt.Println(string(resp_body))
}


type userInfo struct{
	TableName orm.TableName "userinfo"
	userNmae string `name:"username"`
	Uid int `name:"uid"PK:"true"auto:"true"`
	DepartName string `name:"departname"`
	Created string `name:"created"`
}

func ormTest(){
	ui := userInfo{userNmae:"CHAIN", DepartName:"TEST", Created:time.Now().String()}
	orm.Register(new(userInfo))
	db, err := orm.NewDb("mysql", "root:password@tcp(xxx.xx.xxx.xxxx:3306)/demo?charset=utf8")
	if err != nil {
		fmt.Println("打开SQL时出错:", err.Error())
		return
	}
	defer db.Close()
	
	//插入测试
	err = db.Insert(&ui)
	if err != nil {
		fmt.Println("插入时错误:", err.Error())
	}
	fmt.Println("插入成功")
    //修改测试
	ui.userNmae = "BBBB"
	err = db.Update(ui)
	if err != nil {
		fmt.Println("修改时错误:", err.Error())
	}
	fmt.Println("修改成功")
    //删除测试
	err = db.Delete(ui)
	if err != nil {
		fmt.Println("删除时错误:", err.Error())
	}
	fmt.Println("删除成功")
	//查询测试
	res, err := db.From("userinfo").
	Select("username", "departname", "uid").
	Where("uid__gt", 20).
	Where("username", "chain").Get()
	if err != nil{
		fmt.Println("err: ", err.Error())
	}
	fmt.Println(res)
}