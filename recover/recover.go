package recover

import (
	"log"

	"github.com/multiplay/go-slack/chat"
	"github.com/multiplay/go-slack/webhook"
	"github.com/nzgogo/micro/codec"
	"github.com/nzgogo/micro/constant"
)

var slackChannel = webhook.New(constant.SLACKCHANNELADDR)

func Recover(srvName, funcName string, rMsg, nMsg interface{}) {
	log.Println("Recovered in server: " + srvName + " func: " + funcName)
	log.Println(rMsg)
	sendSlackMessage(srvName, funcName, rMsg, nMsg)
}

func sendSlackMessage(srvName, funcName string, rMsg, nMsg interface{}) {
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
	slack_msg := "*Message from* \n> 👉" + srvName + " " + funcName
	m := &chat.Message{Text: slack_msg, Attachments: attachments}
	m.Send(slackChannel)
}
