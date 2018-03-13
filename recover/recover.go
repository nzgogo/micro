package recover

import (
	"log"

	"github.com/multiplay/go-slack/chat"
	"github.com/multiplay/go-slack/webhook"
	"github.com/nzgogo/micro/constant"
)

var slackChannel= webhook.New(constant.SLACKCHANNELADDR)

func Recover(srvName, funcName string) {
	if r := recover(); r != nil {
		log.Println("Recovered in server: "+srvName + " func: " + funcName)
		log.Println(r)
		sendSlackMessage(srvName, funcName, r)
	}
}

func sendSlackMessage(srvName , funcName string, recoverMsg interface{}) {
	attachments := make([]*chat.Attachment, 1)

		msg := &chat.Attachment{
		Title: "func " + funcName + " panic",
		Color: "#FF2D00",
		Text:  recoverMsg.(string),
	}
	attachments = append(attachments,msg)
	slack_msg := "*Message from* \n> 👉" + srvName
	m := &chat.Message{Text: slack_msg, Attachments: attachments}
	m.Send(slackChannel)
}