package fixmssql

import (
	"bytes"
	"encoding/gob"
	"log"
	"strings"
)

func clearText(text interface{}) interface{} {
	if text == nil {
		return nil
	}

	var buf bytes.Buffer
	enc := gob.NewEncoder(&buf)
	err := enc.Encode(text)
	if err != nil {
		return err
	}
	var toString string = string(buf.Bytes())
	toString = strings.Trim(toString, "\t \n \x03 \a \f \v \x0e \x10 \b \x15 \x00 \x12 \x04 \x1b")
	log.Println("clearText: ", text, toString)
	var out interface{}
	out = toString

	return out
}
