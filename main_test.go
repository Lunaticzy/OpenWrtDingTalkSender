package main

import (
	"bytes"
	"fmt"
	"github.com/CatchZeng/dingtalk"
	"github.com/gobuffalo/packr/v2"
	"html/template"
	"log"
	"testing"
)

func TestCheckNetworkConnect(t *testing.T) {
	res := CheckNetworkConnect()
	fmt.Println(res)
}

func TestResolveTime(t *testing.T) {
	timeString := ResolveTime(234234)
	fmt.Println(timeString)
}

func TestSendMessage(t *testing.T) {
	data := make(map[string]interface{})
	networkState := map[string]string{"Baidu": "成功", "Github": "成功", "Google": "成功"}

	data["uptime"] = "test"
	data["memoryUsedPercent"] = "test"
	data["loadAvg"] = "test"
	data["cpuUsedPercent"] = "test"
	data["cpuTemp"] = "test"
	data["wanIP"] = "test"
	data["networkState"] = networkState

	msg := dingtalk.NewMarkdownMessage()

	box := packr.New("templates", "./templates")
	temp, err := box.FindString("notify.md.temp")
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
