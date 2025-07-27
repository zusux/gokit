package encoder

import (
	"fmt"
	"go.uber.org/zap/buffer"
	"go.uber.org/zap/zapcore"
)

type AppPrefixEncoder struct {
	App             string
	Format          string
	zapcore.Encoder // ✅ 直接嵌入
}

func (a *AppPrefixEncoder) Clone() zapcore.Encoder {
	return &AppPrefixEncoder{
		App:     a.App,
		Encoder: a.Encoder.Clone(),
	}
}

func (a *AppPrefixEncoder) EncodeEntry(entry zapcore.Entry, fields []zapcore.Field) (*buffer.Buffer, error) {
	buf, err := a.Encoder.EncodeEntry(entry, fields)
	if err != nil {
		return nil, err
	}

	// 如果是 JSON 编码器，直接返回
	if a.Format == "json" {
		return buf, nil
	}

	// 如果是 console，才添加 app 前缀
	if a.App != "" {
		newBuf := buffer.NewPool().Get()
		newBuf.AppendString(fmt.Sprintf("[%s] ", a.App))
		newBuf.AppendString(buf.String())
		buf.Free()
		return newBuf, nil
	}

	return buf, nil
}
