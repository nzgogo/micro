package transport

import (
	"bytes"
	"encoding/base64"

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

		msgChunk := codec.NewMessage(constant.REQUEST)
		msgChunk.ContextID = msg.ContextID
		msgChunk.Node = constant.FILE_SERVICE_UPLOAD_NODE

		msgChunk.Set("size", total)
		msgChunk.Set("index", counter)
		msgChunk.Set("fileChunk", chunk)

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
