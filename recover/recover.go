package catch

import (
	"log"

	"github.com/multiplay/go-slack/chat"
	"github.com/multiplay/go-slack/webhook"
	"github.com/nzgogo/micro/codec"
)

func Recover(slackurl, srvName, funcName string, nMsg interface{}) {
	rMsg := recover()
	if rMsg != nil {
		log.Println("Recovered in server: " + srvName + " func: " + funcName)
		log.Println(rMsg)
		if nMsg != nil {
			log.Println(nMsg)
		}
		sendSlackMessage(slackurl, srvName, funcName, rMsg, nMsg)
	}
}

func PostProc(slackurl, srvName, funcName string, rMsg, nMsg interface{}) {
	if rMsg != nil {
		log.Println("Recovered in server: " + srvName + " func: " + funcName)
		log.Println(rMsg)
		if nMsg != nil {
			log.Println(nMsg)
		}
		sendSlackMessage(slackurl, srvName, funcName, rMsg, nMsg)
	}
}

func sendSlackMessage(slackurl, srvName, funcName string, rMsg, nMsg interface{}) {
	if slackurl == "" {
		return
	}
	var slackChannel = webhook.New(slackurl)
	var sendMsg = ""
	if rMsg != nil {
		rMsgMal, err := codec.Marshal(rMsg)
		if err == nil {
			sendMsg += string(rMsgMal)
		}
	}
	if nMsg != nil {
		nMsgMal, err := codec.Marshal(nMsg)
		if err == nil {
			sendMsg += string(nMsgMal)
		}
	}
	attachments := make([]*chat.Attachment, 1)
	msg := &chat.Attachment{
		Title: "PANICKING",
		Color: "#FF2D00",
		Text:  sendMsg,
	}
	attachments = append(attachments, msg)
	slack_msg := "*Message from 👉* \n> " + srvName + " " + funcName
	m := &chat.Message{Text: slack_msg, Attachments: attachments}
	m.Send(slackChannel)
}
