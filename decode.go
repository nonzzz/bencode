package bencode

import (
	"fmt"
	"reflect"
	"strconv"
	"strings"
)

const (
	NumericStart   = 0x69 // i
	StringDelim    = 0x3A // :
	DirectroyStart = 0x64 // d
	SliceStart     = 0x6c // l
	EndOfType      = 0x65 // e
	PlusSign       = 0x2B // +
	MinusSign      = 0x2D // -
)

const bencodeSymbol = "bencode" // tag

type decode struct {
	buf []byte
	pos int
	end int
}

type bencodeErorr struct{ error }

func isNumeric(code byte) bool {
	return code >= 48 && code < 58
}

func Decode(buf []byte) (s interface{}, err error) {
	decode := &decode{
		buf: buf,
		end: len(buf),
	}

	defer func() {
		if r := recover(); r != nil {
			if e, ok := r.(bencodeErorr); ok {
				err = e.error
			}
		}
	}()

	return decode.next(), nil
}

func (decode *decode) step() {
	if decode.pos < decode.end {
		decode.pos++
	}
}

func (decode *decode) at(pos int) byte {
	if pos < decode.end {
		return decode.buf[pos]
	}
	return 0
}

// No need to eat operate kind
func (decode *decode) scanBinaryLen(end int) int {
	sum := 0
	// ASCII
	for {
		if decode.pos == end {
			break
		}
		if isNumeric(decode.buf[decode.pos]) {
			sum = sum*10 + (int(decode.buf[decode.pos]) - 48)
		} else {
			if decode.buf[decode.pos] == 46 {
				break
			}
			panic(bencodeErorr{error: fmt.Errorf("invalid binary len: wrong char '%s'", string(decode.buf[decode.pos]))})
		}
		decode.pos++
	}
	return sum
}

func (decode *decode) expect(kind byte) int {
	step := decode.pos
	for step < decode.end {
		if decode.buf[step] == kind {
			return step
		}
		step++
	}
	panic(bencodeErorr{error: fmt.Errorf("Invalid data: Missing delimiter '%s'", string(kind))})
}

func (decode *decode) next() interface{} {
	switch decode.buf[decode.pos] {
	case NumericStart:
		return decode.convertNumeric()
	case DirectroyStart:
		decode.step()
		return decode.convertDirectory()
	case SliceStart:
		decode.step()
		return decode.convertSlice()
	default:
		return decode.convertBinary()
	}
}

func (decode *decode) convertDirectory() (directory map[string]interface{}) {
	directory = make(map[string]interface{})
	for decode.buf[decode.pos] != EndOfType {
		binary := decode.convertBinary()
		directory[string(binary)] = decode.next()
	}
	decode.step()
	return directory
}
func (decode *decode) convertSlice() (list []interface{}) {
	for decode.buf[decode.pos] != EndOfType {
		list = append(list, decode.next())
	}
	decode.step()
	return list
}

func (decode *decode) convertNumeric() int {
	step := decode.expect(EndOfType)
	negative := 1
	//  consume all of operate symbol
	if decode.at(decode.pos+1) == MinusSign {
		negative = -1
		decode.step()
	}
	if decode.at(decode.pos+1) == PlusSign {
		step -= 1
		decode.step()
	}
	decode.step()
	num := decode.scanBinaryLen(step)
	for i := decode.pos; i < step; i++ {
		decode.step()
	}
	decode.step()
	return num * negative
}

func (decode *decode) convertBinary() []byte {
	step := decode.expect(StringDelim)
	l := decode.scanBinaryLen(step)
	for i := decode.pos; i < (step + l); i++ {
		decode.step()
	}
	decode.step()
	return decode.buf[step+1 : decode.pos]
}

type filedInfo struct {
	Alias   string
	TagName string
}

func scanTagField(field reflect.StructField) filedInfo {
	alias := field.Name
	if tag, ok := field.Tag.Lookup(bencodeSymbol); ok {
		tpl := strings.Split(tag, ",")
		// TODO
		if len(tpl) > 0 {
			alias = strings.TrimSpace(tpl[0])
		}
	}

	return filedInfo{
		Alias:   alias,
		TagName: field.Name,
	}
}

func bindTag(decodedMap map[string]interface{}, stu interface{}) error {
	stuType := reflect.TypeOf(stu)
	stuKind := stuType.Kind()

	if stuKind == reflect.Invalid {
		return fmt.Errorf("can't process empty value")
	}

	if stuKind != reflect.Ptr {
		return fmt.Errorf("invalid stu type: should be pointer")
	}

	stuType = stuType.Elem()
	stuKind = stuType.Kind()

	if stuKind != reflect.Struct {
		return fmt.Errorf("invalid stu type: should be struct")
	}

	stuValue := reflect.ValueOf(stu).Elem()

	// TODO
	convert := func(k reflect.Kind, data interface{}) reflect.Value {
		value := reflect.ValueOf(data)
		kind := value.Kind()
		switch k {
		case reflect.String:
			if kind == reflect.Slice || kind == reflect.Array {
				if value.Type().String() == "[]uint8" {
					return reflect.ValueOf(string(value.Bytes()))
				}
			}
			if kind == reflect.Int {
				return reflect.ValueOf(strconv.FormatInt(value.Int(), 10))
			}
		}
		return reflect.ValueOf(data)
	}

	for i := 0; i < stuType.NumField(); i++ {
		field := stuType.Field(i)
		info := scanTagField(field)
		each := stuValue.FieldByName(info.TagName)
		if value, ok := decodedMap[info.Alias]; ok {
			t := reflect.TypeOf(value)
			if !t.AssignableTo(each.Type()) {

				if each.Kind() == reflect.Struct && t.Kind() == reflect.Map {
					if each.CanInterface() {
						bindTag(value.(map[string]interface{}), each.Addr().Interface())
					}
				} else {
					val := convert(each.Kind(), value)
					each.Set(val)
				}

			} else {
				each.Set(reflect.ValueOf(value))
			}

		}
	}
	return nil
}

func UnMarshal(data interface{}, stu interface{}) error {
	value := reflect.ValueOf(data)

	kind := value.Kind()

	if kind == reflect.Invalid {
		return fmt.Errorf("can't process empty value")
	}

	if kind != reflect.Map {
		return fmt.Errorf("can't process type: %s ", kind)
	}

	decodedMap := make(map[string]interface{}, value.Len())
	iter := value.MapRange()
	for iter.Next() {
		decodedMap[iter.Key().String()] = iter.Value().Interface()
	}

	return bindTag(decodedMap, stu)
}
