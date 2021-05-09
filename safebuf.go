package zapflatencoder

import (
	"unicode/utf8"

	"go.uber.org/zap/buffer"
)

func getBuffer() *safeBuf {
	return &safeBuf{_bufPool.Get()}
}

type safeBuf struct {
	*buffer.Buffer
}

func (buf *safeBuf) safeAddString(s string) {
	for i := 0; i < len(s); {
		if buf.tryAddRuneSelf(s[i]) {
			i++
			continue
		}
		r, size := utf8.DecodeRuneInString(s[i:])
		if buf.tryAddRuneError(r, size) {
			i++
			continue
		}
		buf.AppendString(s[i : i+size])
		i += size
	}
}

func (buf *safeBuf) safeAddByteString(s []byte) {
	for i := 0; i < len(s); {
		if buf.tryAddRuneSelf(s[i]) {
			i++
			continue
		}
		r, size := utf8.DecodeRune(s[i:])
		if buf.tryAddRuneError(r, size) {
			i++
			continue
		}
		_, _ = buf.Write(s[i : i+size])
		i += size
	}
}

func (buf *safeBuf) tryAddRuneSelf(b byte) bool {
	if b >= utf8.RuneSelf {
		return false
	}
	buf.AppendByte(b)
	return true
}

func (buf *safeBuf) tryAddRuneError(r rune, size int) bool {
	if r == utf8.RuneError && size == 1 {
		buf.AppendString(TokenReplacement)
		return true
	}
	return false
}
