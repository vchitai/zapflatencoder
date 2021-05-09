package zapflatencoder

import (
	"encoding/base64"
	"encoding/json"
	"math"
	"sync"
	"time"

	"go.uber.org/zap/buffer"
	"go.uber.org/zap/zapcore"
)

var _ zapcore.ObjectEncoder = &objectEncoder{}

type objectEncoder struct {
	*zapcore.EncoderConfig
	namespaces int
	buf        *safeBuf
}

func (enc *objectEncoder) flush(buf *buffer.Buffer) *buffer.Buffer {
	enc.closeOpenNamespaces()
	if enc.buf.Len() == 0 {
		return buf
	}
	if buf.Len() > 0 {
		buf.AppendByte(tokenTab)
	}
	_, _ = buf.Write(enc.buf.Bytes())
	return buf
}

var _objectEncoderPool = sync.Pool{New: func() interface{} {
	return &objectEncoder{
		namespaces: 0,
	}
}}

func getObjectEncoder(cfg *zapcore.EncoderConfig) *objectEncoder {
	enc := _objectEncoderPool.Get().(*objectEncoder)
	enc.EncoderConfig = cfg
	enc.buf = getBuffer()
	return enc
}

func putObjectEncoder(enc *objectEncoder) {
	enc.EncoderConfig = nil
	enc.buf.Free()
	_objectEncoderPool.Put(enc)
}

func (enc *objectEncoder) AddArray(key string, arr zapcore.ArrayMarshaler) error {
	enc.addKey(key)
	return enc.AppendArray(arr)
}

func (enc *objectEncoder) AddObject(key string, obj zapcore.ObjectMarshaler) error {
	enc.addKey(key)
	return enc.AppendObject(obj)
}

func (enc *objectEncoder) AddBinary(key string, val []byte) {
	enc.AddString(key, base64.StdEncoding.EncodeToString(val))
}

func (enc *objectEncoder) AddByteString(key string, val []byte) {
	enc.addKey(key)
	enc.AppendByteString(val)
}

func (enc *objectEncoder) AddBool(key string, val bool) {
	enc.addKey(key)
	enc.AppendBool(val)
}

func (enc *objectEncoder) AddComplex128(key string, val complex128) {
	enc.addKey(key)
	enc.AppendComplex128(val)
}

func (enc *objectEncoder) AddDuration(key string, val time.Duration) {
	enc.addKey(key)
	enc.AppendDuration(val)
}

func (enc *objectEncoder) AddFloat64(key string, val float64) {
	enc.addKey(key)
	enc.AppendFloat64(val)
}

func (enc *objectEncoder) AddInt64(key string, val int64) {
	enc.addKey(key)
	enc.AppendInt64(val)
}

func (enc *objectEncoder) AddReflected(key string, obj interface{}) error {
	marshaled, err := json.Marshal(obj)
	if err != nil {
		return err
	}
	enc.addKey(key)
	_, err = enc.buf.Write(marshaled)
	return err
}

func (enc *objectEncoder) OpenNamespace(key string) {
	enc.addKey(key)
	enc.buf.AppendByte(tokenNamespaceOpen)
	enc.namespaces++
}

func (enc *objectEncoder) AddString(key, val string) {
	enc.addKey(key)
	enc.AppendString(val)
}

func (enc *objectEncoder) AddTime(key string, val time.Time) {
	enc.addKey(key)
	enc.AppendTime(val)
}

func (enc *objectEncoder) AddUint64(key string, val uint64) {
	enc.addKey(key)
	enc.AppendUint64(val)
}

func (enc *objectEncoder) AppendArray(arr zapcore.ArrayMarshaler) error {
	enc.buf.AppendByte(tokenArrayOpen)
	err := arr.MarshalLogArray(enc)
	enc.buf.AppendByte(tokenArrayClose)
	return err
}

func (enc *objectEncoder) AppendObject(obj zapcore.ObjectMarshaler) error {
	enc.buf.AppendByte(tokenNamespaceOpen)
	err := obj.MarshalLogObject(enc)
	enc.buf.AppendByte(tokenNamespaceClose)
	return err
}

func (enc *objectEncoder) AppendBool(val bool) {
	enc.buf.AppendBool(val)
}

func (enc *objectEncoder) AppendByteString(val []byte) {
	enc.buf.AppendByte(tokenStringEnclosed)
	enc.buf.safeAddByteString(val)
	enc.buf.AppendByte(tokenStringEnclosed)
}

func (enc *objectEncoder) AppendComplex128(val complex128) {
	r, i := real(val), imag(val)
	enc.buf.AppendByte(tokenStringEnclosed)
	enc.buf.AppendFloat(r, 64)
	enc.buf.AppendByte('+')
	enc.buf.AppendFloat(i, 64)
	enc.buf.AppendByte('i')
	enc.buf.AppendByte(tokenStringEnclosed)
}

func (enc *objectEncoder) AppendDuration(val time.Duration) {
	cur := enc.buf.Len()
	enc.EncodeDuration(val, enc)
	if cur == enc.buf.Len() {
		enc.AppendInt64(int64(val))
	}
}

func (enc *objectEncoder) AppendInt64(val int64) {
	enc.buf.AppendInt(val)
}

func (enc *objectEncoder) AppendReflected(val interface{}) error {
	marshaled, err := json.Marshal(val)
	if err != nil {
		return err
	}
	_, err = enc.buf.Write(marshaled)
	return err
}

func (enc *objectEncoder) AppendString(val string) {
	enc.buf.AppendByte(tokenStringEnclosed)
	enc.buf.safeAddString(val)
	enc.buf.AppendByte(tokenStringEnclosed)
}

func (enc *objectEncoder) AppendTime(val time.Time) {
	cur := enc.buf.Len()
	enc.EncodeTime(val, enc)
	if cur == enc.buf.Len() {
		enc.AppendInt64(val.UnixNano())
	}
}

func (enc *objectEncoder) AppendUint64(val uint64) {
	enc.buf.AppendUint(val)
}

func (enc *objectEncoder) AddComplex64(k string, v complex64) { enc.AddComplex128(k, complex128(v)) }
func (enc *objectEncoder) AddFloat32(k string, v float32)     { enc.AddFloat64(k, float64(v)) }
func (enc *objectEncoder) AddInt(k string, v int)             { enc.AddInt64(k, int64(v)) }
func (enc *objectEncoder) AddInt32(k string, v int32)         { enc.AddInt64(k, int64(v)) }
func (enc *objectEncoder) AddInt16(k string, v int16)         { enc.AddInt64(k, int64(v)) }
func (enc *objectEncoder) AddInt8(k string, v int8)           { enc.AddInt64(k, int64(v)) }
func (enc *objectEncoder) AddUint(k string, v uint)           { enc.AddUint64(k, uint64(v)) }
func (enc *objectEncoder) AddUint32(k string, v uint32)       { enc.AddUint64(k, uint64(v)) }
func (enc *objectEncoder) AddUint16(k string, v uint16)       { enc.AddUint64(k, uint64(v)) }
func (enc *objectEncoder) AddUint8(k string, v uint8)         { enc.AddUint64(k, uint64(v)) }
func (enc *objectEncoder) AddUintptr(k string, v uintptr)     { enc.AddUint64(k, uint64(v)) }
func (enc *objectEncoder) AppendComplex64(v complex64)        { enc.AppendComplex128(complex128(v)) }
func (enc *objectEncoder) AppendFloat64(v float64)            { enc.appendFloat(v, 64) }
func (enc *objectEncoder) AppendFloat32(v float32)            { enc.appendFloat(float64(v), 32) }
func (enc *objectEncoder) AppendInt(v int)                    { enc.AppendInt64(int64(v)) }
func (enc *objectEncoder) AppendInt32(v int32)                { enc.AppendInt64(int64(v)) }
func (enc *objectEncoder) AppendInt16(v int16)                { enc.AppendInt64(int64(v)) }
func (enc *objectEncoder) AppendInt8(v int8)                  { enc.AppendInt64(int64(v)) }
func (enc *objectEncoder) AppendUint(v uint)                  { enc.AppendUint64(uint64(v)) }
func (enc *objectEncoder) AppendUint32(v uint32)              { enc.AppendUint64(uint64(v)) }
func (enc *objectEncoder) AppendUint16(v uint16)              { enc.AppendUint64(uint64(v)) }
func (enc *objectEncoder) AppendUint8(v uint8)                { enc.AppendUint64(uint64(v)) }
func (enc *objectEncoder) AppendUintptr(v uintptr)            { enc.AppendUint64(uint64(v)) }

func (enc *objectEncoder) closeOpenNamespaces() {
	for i := 0; i < enc.namespaces; i++ {
		enc.buf.AppendByte(tokenNamespaceClose)
	}
}

func (enc *objectEncoder) addKey(key string) {
	enc.buf.AppendByte(tokenTab)
	enc.buf.safeAddString(key)
	enc.buf.AppendByte(tokenKeyValueSeparator)
}

func (enc *objectEncoder) appendFloat(val float64, bitSize int) {
	switch {
	case math.IsNaN(val):
		enc.buf.AppendString(`"NaN"`)
	case math.IsInf(val, 1):
		enc.buf.AppendString(`"+Inf"`)
	case math.IsInf(val, -1):
		enc.buf.AppendString(`"-Inf"`)
	default:
		enc.buf.AppendFloat(val, bitSize)
	}
}
