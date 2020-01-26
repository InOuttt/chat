package adiraFinance

import (
	"bytes"
	"encoding/gob"
	"encoding/json"
	"log"
	"path/filepath"
	"runtime"
	"strings"
)

// IfcToString ...
func IfcToString(src interface{}) string {
	if src == nil || src.(string) == "" {
		return ""
	}

	var buf bytes.Buffer
	enc := gob.NewEncoder(&buf)
	err := enc.Encode(src)
	if err != nil {
		Log.Error("Failed gob encode")
		return ""
	}

	var toString string = string(buf.Bytes())
	toString = ClearText(toString)
	// Log.Ln("IfcToString(): ", len(toString), ":", len(src.(string)))

	return toString
}

// ClearText ...
func ClearText(text string) string {
	if text == "" {
		return ""
	}

	return strings.Trim(
		text,
		"\t \n \x03 \a \f \v \x0e \x10 \b \x15 \x00 \x12 \x04 \x1b \x06 \x05  "+
			"\x13 \x11 \x17 \x14 # \x16 þ",
	)
}

// ClearJSON ...
func ClearJSON(text string) string {
	// fix first char is "
	if len(text) > 2 && text[:2] == "\"\f" {
		text = text[2:]
	}

	strArr := strings.Split(text, "")
	var removeString string

	for _, s := range strArr {
		if s != "{" && s != "[" && s != "\"" {
			removeString += s
			continue
		}
		break
	}

	cleaned := strings.Trim(text, removeString)

	if cleaned[:1] == "{" && cleaned[len(cleaned)-1:] != "}" {
		cleaned += "}"
	} else if cleaned[:1] == "[" && cleaned[len(cleaned)-1:] != "]" {
		cleaned += "]"
	}

	return cleaned
}

// StrToJSON ...
func StrToJSON(j string) (interface{}, error) {
	log.Println(j[:10] + " ... " + j[len(j)-10:])
	var out interface{}
	err := json.Unmarshal([]byte(j), &out)
	if out != nil {
		return out, nil
	} else if err != nil {
		if len(j) <= 512 {
			Log.Ln(j)
		}
		Log.Error(err)
		_, pathfile, line, ok := runtime.Caller(1)
		if ok {
			log.Printf("called by %s#%d", filepath.Base(pathfile), line)
		}
	}

	return nil, err
}
