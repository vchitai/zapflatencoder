package zapflatencoder

import (
	"fmt"
	"sync"
	"time"

	"go.uber.org/zap/buffer"
	"go.uber.org/zap/zapcore"
)

var _ zapcore.PrimitiveArrayEncoder = &sliceEncoder{}

var _sliceEncoderPool = sync.Pool{
	New: func() interface{} {
		return &sliceEncoder{elems: make([]interface{}, 0, 2)}
	},
}

func getSliceEncoder() *sliceEncoder {
	return _sliceEncoderPool.Get().(*sliceEncoder)
}

func putSliceEncoder(e *sliceEncoder) {
	e.elems = e.elems[:0]
	_sliceEncoderPool.Put(e)
}

type sliceEncoder struct {
	elems []interface{}
}

func (s *sliceEncoder) flush(buf *buffer.Buffer) *buffer.Buffer {
	for i := range s.elems {
		if i > 0 {
			buf.AppendByte(tokenTab)
		}
		_, _ = fmt.Fprint(buf, s.elems[i])
	}
	return buf
}

func (s *sliceEncoder) AppendArray(v zapcore.ArrayMarshaler) error {
	enc := &sliceEncoder{}
	err := v.MarshalLogArray(enc)
	s.elems = append(s.elems, enc.elems)
	return err
}

func (s *sliceEncoder) AppendObject(v zapcore.ObjectMarshaler) error {
	m := zapcore.NewMapObjectEncoder()
	err := v.MarshalLogObject(m)
	s.elems = append(s.elems, m.Fields)
	return err
}

func (s *sliceEncoder) AppendReflected(v interface{}) error {
	s.elems = append(s.elems, v)
	return nil
}

func (s *sliceEncoder) AppendBool(v bool)              { s.elems = append(s.elems, v) }
func (s *sliceEncoder) AppendByteString(v []byte)      { s.elems = append(s.elems, v) }
func (s *sliceEncoder) AppendComplex128(v complex128)  { s.elems = append(s.elems, v) }
func (s *sliceEncoder) AppendComplex64(v complex64)    { s.elems = append(s.elems, v) }
func (s *sliceEncoder) AppendDuration(v time.Duration) { s.elems = append(s.elems, v) }
func (s *sliceEncoder) AppendFloat64(v float64)        { s.elems = append(s.elems, v) }
func (s *sliceEncoder) AppendFloat32(v float32)        { s.elems = append(s.elems, v) }
func (s *sliceEncoder) AppendInt(v int)                { s.elems = append(s.elems, v) }
func (s *sliceEncoder) AppendInt64(v int64)            { s.elems = append(s.elems, v) }
func (s *sliceEncoder) AppendInt32(v int32)            { s.elems = append(s.elems, v) }
func (s *sliceEncoder) AppendInt16(v int16)            { s.elems = append(s.elems, v) }
func (s *sliceEncoder) AppendInt8(v int8)              { s.elems = append(s.elems, v) }
func (s *sliceEncoder) AppendString(v string)          { s.elems = append(s.elems, v) }
func (s *sliceEncoder) AppendTime(v time.Time)         { s.elems = append(s.elems, v) }
func (s *sliceEncoder) AppendUint(v uint)              { s.elems = append(s.elems, v) }
func (s *sliceEncoder) AppendUint64(v uint64)          { s.elems = append(s.elems, v) }
func (s *sliceEncoder) AppendUint32(v uint32)          { s.elems = append(s.elems, v) }
func (s *sliceEncoder) AppendUint16(v uint16)          { s.elems = append(s.elems, v) }
func (s *sliceEncoder) AppendUint8(v uint8)            { s.elems = append(s.elems, v) }
func (s *sliceEncoder) AppendUintptr(v uintptr)        { s.elems = append(s.elems, v) }
