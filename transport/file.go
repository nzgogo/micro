package transport

import (
	"bytes"
	"encoding/base64"
	"log"

	"github.com/nzgogo/micro/codec"
	"github.com/nzgogo/micro/constant"
)

func chunkCount(a int, b int) int {
	if a%b > 0 {
		return int(a/b) + 1
	}
	return int(a / b)
}

func (n *transport) SendFile(msg *codec.Message, sub string, file string) (err error) {
	b, err := base64.StdEncoding.DecodeString(file)
	if err != nil {
		return err
	}
	fileSize := len(b)
	total := chunkCount(fileSize, constant.MAX_FILE_CHUNK_SIZE)
	fileReader := bytes.NewReader(b)

	for counter := 0; counter < total; counter++ {
		var chunk []byte
		if fileReader.Size() >= constant.MAX_FILE_CHUNK_SIZE {
			chunk = make([]byte, constant.MAX_FILE_CHUNK_SIZE)
		} else {
			chunk = make([]byte, fileReader.Size())
		}
		fileReader.Read(chunk)
		log.Printf("Chunk %d size: %d\n", counter, len(chunk))

		msgBody := make(map[string]interface{})
		msgBody["size"] = total
		msgBody["index"] = counter
		msgBody["fileChunk"] = chunk

		msgBodyBytes, err := codec.Marshal(msgBody)
		if err != nil {
			return err
		}

		msgChunk := &codec.Message{
			ContextID: msg.ContextID,
			Type:      constant.REQUEST,
			Node:      constant.FILE_SERVICE_UPLOAD_NODE,
			Body:      msgBodyBytes,
		}
		log.Println(msgChunk)

		msgBytes, err := codec.Marshal(msgChunk)
		if err != nil {
			return err
		}

		log.Printf("Message %d size: %d\n", counter, len(msgBytes))

		err = n.Publish(sub, msgBytes)
		if err != nil {
			return err
		}
	}

	return
}
