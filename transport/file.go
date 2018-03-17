package transport

import (
	"bytes"

	"github.com/nzgogo/micro/codec"
	"github.com/nzgogo/micro/constant"
)

func chunkCount(a int, b int) int {
	if a%b > 0 {
		return int(a/b) + 1
	}
	return int(a / b)
}

func (n *transport) SendFile(msg *codec.Message, sub string, file []byte) (err error) {
	fileSize := len(file)
	total := chunkCount(fileSize, constant.MAX_FILE_CHUNK_SIZE)
	fileReader := bytes.NewReader(file)

	for counter := 0; counter < total; counter++ {
		chunk := make([]byte, 0)
		fileReader.ReadAt(chunk, constant.MAX_FILE_CHUNK_SIZE)

		msgBody := make(map[string]interface{})
		msgBody["size"] = total
		msgBody["index"] = counter
		msgBody["fileChunk"] = chunk

		msgBodyBytes, err := codec.Marshal(msgBody)
		if err != nil {
			return err
		}

		msgChunk := *msg
		msg.Body = msgBodyBytes

		msgBytes, err := codec.Marshal(msgChunk)
		if err != nil {
			return err
		}

		err = n.Publish(sub, msgBytes)
		if err != nil {
			return err
		}
	}

	return
}
