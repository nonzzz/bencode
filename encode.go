// This Package implements a simple bencode encode/decode
// bencode(https://en.wikipedia.org/wiki/Bencode)
// Lexicographic order(https://en.wikipedia.org/wiki/Lexicographic_order)

package bencode

import (
	"fmt"
	"reflect"
	"sort"
	"strconv"
	"strings"
)

type encode struct {
	output []byte
}

func Encode(input interface{}) (s []byte, err error) {
	encode := &encode{}

	defer func() {
		if r := recover(); r != nil {
			if e, ok := r.(bencodeErorr); ok {
				err = e.error
			}
		}
	}()

	encode.next(input)
	return encode.output, nil
}

func (encode *encode) next(input interface{}) {
	value := reflect.ValueOf(input)
	kind := value.Kind()
	switch kind {
	case reflect.Slice, reflect.Array:
		// string alias
		if value.Type().String() == "[]uint8" {
			encode.encodeString(string(value.Bytes()))
		} else {
			encode.encodeSlice(value)
		}
	case reflect.Map:
		encode.encodeDirectory(value)
	case reflect.String:
		encode.encodeString(value.String())
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32,
		reflect.Int64, reflect.Uint, reflect.Uint8, reflect.Uint16,
		reflect.Uint32, reflect.Uint64:
		encode.encodeNumeric(strconv.FormatInt(value.Int(), 10))
	default:
		panic(bencodeErorr{error: fmt.Errorf("Invalid type: can't support %s", kind)})
	}
}

// Someone need convert WTF string(non-standard ASCII) if necessary.

func (encode *encode) encodeString(s string) {
	var sb strings.Builder
	l := strconv.Itoa(len(s))
	sb.WriteString(l)
	sb.WriteByte(StringDelim)
	sb.WriteString(s)
	encode.output = append(encode.output, sb.String()...)
}

// Should sort directory
func (encode *encode) encodeDirectory(value reflect.Value) {
	encode.output = append(encode.output, DirectroyStart)
	keys := value.MapKeys()
	sort.Slice(keys, func(a, b int) bool {
		return keys[a].String() < keys[b].String()
	})
	for _, k := range keys {
		encode.encodeString(k.String())
		encode.next(value.MapIndex(k).Interface())
	}

	encode.output = append(encode.output, EndOfType)
}

func (encode *encode) encodeSlice(value reflect.Value) {
	encode.output = append(encode.output, SliceStart)
	for i := 0; i < value.Len(); i++ {
		encode.next(value.Index(i).Interface())
	}
	encode.output = append(encode.output, EndOfType)
}

func (encode *encode) encodeNumeric(s string) {
	content := make([]byte, 0, 2+len(s))
	content = append(content, NumericStart)
	content = append(content, s...)
	content = append(content, EndOfType)
	encode.output = append(encode.output, content...)
}
