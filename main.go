package main

import (
	"bytes"
	"fmt"
	"github.com/CatchZeng/dingtalk"
	"github.com/buger/jsonparser"
	"github.com/go-ini/ini"
	"github.com/gobuffalo/packr/v2"
	"github.com/shirou/gopsutil/v3/cpu"
	"github.com/shirou/gopsutil/v3/host"
	"github.com/shirou/gopsutil/v3/load"
	"github.com/shirou/gopsutil/v3/mem"
	"html/template"
	"io/ioutil"
	"log"
	"os/exec"
	"runtime"
	"time"
)

//go:generate packr2
func main() {

	data := make(map[string]string)

	loadAvg, _ := load.Avg()
	percent, _ := cpu.Percent(3*time.Second, false)
	memory, _ := mem.VirtualMemory()
	uptime, _ := host.Uptime()

	data["uptime"] = resolveTime(uptime)
	data["memoryUsedPercent"] = fmt.Sprintf("%.2f", memory.UsedPercent)
	data["loadAvg"] = fmt.Sprintf("%.2f, %.2f, %.2f", loadAvg.Load1, loadAvg.Load5, loadAvg.Load15)
	data["cpuUsedPercent"] = fmt.Sprintf("%.2f", percent[0])
	data["cpuTemp"] = GetTemp()
	data["wanIP"] = GetWanIP()

	msg := dingtalk.NewMarkdownMessage()

	box := packr.New("templates", "./templates")
	temp, err := box.FindString("notify.md")
	if err != nil {
		log.Panic(err)
	}

	buffer := bytes.NewBuffer([]byte{})
	parse, err := template.New("test").Parse(temp)
	if err != nil {
		log.Panic(err)
	}

	err = parse.Execute(buffer, data)
	if err != nil {
		log.Panic(err)
	}

	msg.SetAt([]string{}, false)
	msg.SetMarkdown("test", buffer.String())

	SendMessage(msg)

}

/*
发送钉钉信息
*/
func SendMessage(message dingtalk.Message) {

	box := packr.New("conf", "./conf")
	_conf, err := box.Find("conf.ini")
	if err != nil {
		log.Panic(err)
	}

	conf, err := ini.Load(_conf)
	if err != nil {
		log.Panic(err)
	}

	accessToken := conf.Section("config").Key("token").String()
	secret := conf.Section("config").Key("secret").String()
	client := dingtalk.NewClient(accessToken, secret)

	status, err := client.Send(message)
	if err != nil {
		log.Panic(err)
	}
	log.Println(status)
}

/*
获取温度
*/
func GetTemp() (temp string) {
	if runtime.GOOS != "linux" {
		temp = "0"
		return
	}

	t, err := ioutil.ReadFile("/sys/class/thermal/thermal_zone0/temp")
	if err != nil {
		log.Fatalln(err)
	}

	temp = string(t[0:2])

	return
}

const (

	//定义每分钟的秒数
	Minute = 60
	//定义每小时的秒数
	Hour = Minute * 60
	//定义每天的秒数
	Day = Hour * 24
)

/*
时间转换函数
*/
func resolveTime(seconds uint64) (timeString string) {

	if seconds < Minute {
		timeString = fmt.Sprintf("%d秒", seconds)
	} else if seconds < Hour {
		timeMin := seconds / Minute
		timeSec := timeMin * Minute
		timeSec = seconds - timeSec
		timeString = fmt.Sprintf("%d分%d秒", timeMin, timeSec)
	} else if seconds < Day {
		timeHour := seconds / Hour
		timeMin := timeHour * Hour
		timeMin = (seconds - timeMin) / Minute

		timeString = fmt.Sprintf("%d小时%d分", timeHour, timeMin)
	} else {
		timeDay := seconds / Day
		timeHour := timeDay * Day
		timeHour = (seconds - timeHour) / Hour
		timeString = fmt.Sprintf("%d天%d小时", timeDay, timeHour)

	}

	return
}

func GetWanIP() (ip string) {

	cmd := exec.Command("bash", "-c", `ubus call network.interface.wan status`)
	output, err := cmd.CombinedOutput()
	if err != nil {
		log.Panic(err)
	}

	ip, err = jsonparser.GetString(output, "ipv4-address", "[0]", "address")
	if err != nil {
		log.Panic(err)
	}

	return
}
