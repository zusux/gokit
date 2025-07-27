package encoder

import (
	"go.uber.org/zap/buffer"
	"go.uber.org/zap/zapcore"
)

type AppPrefixEncoder struct {
	App             string
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
	if a.App == "" {
		return buf, nil
	}

	// 追加前缀 [App]
	prefixed := buffer.NewPool().Get()
	prefixed.AppendString("[" + a.App + "] ")
	prefixed.AppendString(buf.String())
	buf.Free()
	return prefixed, nil
}
