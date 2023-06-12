package helper_buffer_formatter

import (
	"bytes"
)

func EncodeBuffer(data []byte) *bytes.Buffer{
	return bytes.NewBuffer(data)
}