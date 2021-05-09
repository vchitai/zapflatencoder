package zapflatencoder

import (
	"sync"

	"go.uber.org/zap/buffer"
	"go.uber.org/zap/zapcore"
)

var _bufPool = buffer.NewPool()

var _consolePool = sync.Pool{New: func() interface{} {
	return &encoder{}
}}

func getEncoder() *encoder {
	return _consolePool.Get().(*encoder)
}

type encoder struct {
	*zapcore.EncoderConfig
	*objectEncoder
	lineEnding string
}

func newEncoder(cfg zapcore.EncoderConfig) *encoder {
	lineEnding := tokenLineEnding
	if len(cfg.LineEnding) > 0 {
		lineEnding = cfg.LineEnding
	}
	return &encoder{
		EncoderConfig: &cfg,
		objectEncoder: getObjectEncoder(&cfg),
		lineEnding:    lineEnding,
	}
}

func (enc *encoder) Clone() zapcore.Encoder {
	clone := getEncoder()
	clone.EncoderConfig = enc.EncoderConfig
	clone.lineEnding = enc.lineEnding
	clone.objectEncoder = getObjectEncoder(enc.EncoderConfig)
	return clone
}

func (enc *encoder) EncodeMessage(message string, aEnc zapcore.PrimitiveArrayEncoder) {
	aEnc.AppendString(message)
}

func (enc *encoder) EncodeEntry(ent zapcore.Entry, fields []zapcore.Field) (*buffer.Buffer, error) {
	newLine := _bufPool.Get()

	newLine = enc.encodeEntry(newLine, ent)
	// Add any structured context.
	newLine = enc.encodeContext(newLine, fields)

	newLine.AppendString(enc.lineEnding)

	newLine = enc.encodeStacktrace(newLine, ent)

	return newLine, nil
}
func (enc *encoder) encodeStacktrace(buf *buffer.Buffer, ent zapcore.Entry) *buffer.Buffer {
	// If there's no traceback key, honor that; this allows users to force single-line output.
	if len(ent.Stack) > 0 && len(enc.StacktraceKey) > 0 {
		buf.AppendString(ent.Stack)
		buf.AppendString(enc.lineEnding)
	}
	return buf
}
func (enc *encoder) encodeEntry(buf *buffer.Buffer, ent zapcore.Entry) *buffer.Buffer {
	sliceEnc := getNextSliceEncoder()
	defer putNextSliceEncoder(sliceEnc)

	if len(enc.TimeKey) > 0 && enc.EncodeTime != nil {
		enc.EncodeTime(ent.Time, sliceEnc)
	}
	if len(enc.LevelKey) > 0 && enc.EncodeLevel != nil {
		enc.EncodeLevel(ent.Level, sliceEnc)
	}
	if len(enc.MessageKey) > 0 {
		enc.EncodeMessage(ent.Message, sliceEnc)
	}
	if ent.Caller.Defined && len(enc.CallerKey) > 0 && enc.EncodeCaller != nil {
		enc.EncodeCaller(ent.Caller, sliceEnc)
	}

	return sliceEnc.flush(buf)
}

func (enc *encoder) encodeContext(buf *buffer.Buffer, extra []zapcore.Field) *buffer.Buffer {
	objEnc := enc.objectEncoder
	defer func() {
		putObjectEncoder(objEnc)
		enc.objectEncoder = getObjectEncoder(enc.EncoderConfig)
	}()

	for i := range extra {
		extra[i].AddTo(objEnc)
	}

	return objEnc.flush(buf)
}
