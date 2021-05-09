package zapflatencoder

import (
	"fmt"
	"sync"
	"time"

	"go.uber.org/zap/buffer"
	"go.uber.org/zap/zapcore"
)

var _ zapcore.PrimitiveArrayEncoder = &nextSliceEncoder{}

var _nextSliceEncoderPool = sync.Pool{
	New: func() interface{} {
		return &nextSliceEncoder{}
	},
}

func getNextSliceEncoder() *nextSliceEncoder {
	enc := _nextSliceEncoderPool.Get().(*nextSliceEncoder)
	enc.buf = _bufPool.Get()
	return enc
}

func putNextSliceEncoder(e *nextSliceEncoder) {
	e.buf.Free()
	_nextSliceEncoderPool.Put(e)
}

type nextSliceEncoder struct {
	buf *buffer.Buffer
}

func (s *nextSliceEncoder) flush(buf *buffer.Buffer) *buffer.Buffer {
	_, _ = buf.Write(s.buf.Bytes())
	return buf
}

func (s *nextSliceEncoder) add(x interface{}) {
	if s.buf.Len() > 0 {
		s.buf.AppendByte(TokenTab)
	}
	_, _ = fmt.Fprint(s.buf, x)
}

func (s *nextSliceEncoder) AppendArray(v zapcore.ArrayMarshaler) error {
	enc := getNextSliceEncoder()
	defer putNextSliceEncoder(enc)
	err := v.MarshalLogArray(enc)
	s.buf = enc.flush(s.buf)
	return err
}

func (s *nextSliceEncoder) AppendObject(v zapcore.ObjectMarshaler) error {
	m := zapcore.NewMapObjectEncoder()
	err := v.MarshalLogObject(m)
	s.add(m.Fields)
	return err
}

func (s *nextSliceEncoder) AppendReflected(v interface{}) error {
	s.add(v)
	return nil
}

func (s *nextSliceEncoder) AppendBool(v bool)              { s.add(v) }
func (s *nextSliceEncoder) AppendByteString(v []byte)      { s.add(v) }
func (s *nextSliceEncoder) AppendComplex128(v complex128)  { s.add(v) }
func (s *nextSliceEncoder) AppendComplex64(v complex64)    { s.add(v) }
func (s *nextSliceEncoder) AppendDuration(v time.Duration) { s.add(v) }
func (s *nextSliceEncoder) AppendFloat64(v float64)        { s.add(v) }
func (s *nextSliceEncoder) AppendFloat32(v float32)        { s.add(v) }
func (s *nextSliceEncoder) AppendInt(v int)                { s.add(v) }
func (s *nextSliceEncoder) AppendInt64(v int64)            { s.add(v) }
func (s *nextSliceEncoder) AppendInt32(v int32)            { s.add(v) }
func (s *nextSliceEncoder) AppendInt16(v int16)            { s.add(v) }
func (s *nextSliceEncoder) AppendInt8(v int8)              { s.add(v) }
func (s *nextSliceEncoder) AppendString(v string)          { s.add(v) }
func (s *nextSliceEncoder) AppendTime(v time.Time)         { s.add(v) }
func (s *nextSliceEncoder) AppendUint(v uint)              { s.add(v) }
func (s *nextSliceEncoder) AppendUint64(v uint64)          { s.add(v) }
func (s *nextSliceEncoder) AppendUint32(v uint32)          { s.add(v) }
func (s *nextSliceEncoder) AppendUint16(v uint16)          { s.add(v) }
func (s *nextSliceEncoder) AppendUint8(v uint8)            { s.add(v) }
func (s *nextSliceEncoder) AppendUintptr(v uintptr)        { s.add(v) }
