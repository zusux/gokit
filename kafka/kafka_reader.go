// Copyright 2019 GitBitEx.com
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     https://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package kafka

import (
	"context"
	"fmt"
	"io"
	"log"
	"os"

	"github.com/segmentio/kafka-go"
)

type KafkaReader struct {
	reader *kafka.Reader
}

func NewKafkaReader(brokers []string, groupID, topic string, out io.Writer) *KafkaReader {
	s := &KafkaReader{}
	s.reader = kafka.NewReader(kafka.ReaderConfig{
		Brokers: brokers,
		Topic:   topic,
		//Partition: 0,
		GroupID:     groupID,
		MinBytes:    1,
		MaxBytes:    10e6,
		ErrorLogger: log.New(io.MultiWriter(out, os.Stderr), "[kafkaReaderERROR]: ", log.Ldate|log.Ltime|log.Lshortfile),
	})
	return s
}

func (s *KafkaReader) SetOffset(offset int64) error {
	return s.reader.SetOffset(offset)
}

func (s *KafkaReader) FetchMsg(ctx context.Context) (msg kafka.Message, err error) {
	defer func() {
		// 发生宕机时，获取panic传递的上下文并打印
		e := recover()
		if e != nil {
			err = fmt.Errorf("[KafkaReader]-runtime error:%v", e)
		}
		return
	}()
	msg, err = s.reader.FetchMessage(ctx)
	return
}

func (s *KafkaReader) CommitMsg(ctx context.Context, msg kafka.Message) error {
	return s.reader.CommitMessages(ctx, msg)
}

func (s *KafkaReader) Close() error {
	s.reader.Config().Logger.Printf("Close KafkaReader...")
	return s.reader.Close()
}
