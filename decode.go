package bencode

import "errors"

const (
	NumericStart   = 0x69 // i
	StringDelim    = 0x3A // :
	DirectroyStart = 0x64 // d
	SliceStart     = 0x6c // l
	EndOfType      = 0x65 // e
)

type decode struct {
	buf []byte
	pos int
	end int
}

type bencodeErorr struct{ error }

// TODO support negative number
// unsigned integer
func IsNumeric(code byte) bool {
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

func (decode *decode) scanneBinaryLen(end int) int {
	start := decode.pos
	sum := 0
	negative := 1
	// ASCII
	for i := start; i < end; i++ {
		if IsNumeric(decode.buf[i]) {
			sum = sum*10 + (int(decode.buf[i]) - 48)
			continue
		}
		if i == start && decode.buf[start] == 45 {
			negative = -1
			continue
		}
		if decode.buf[start] == 46 {
			break
		}
	}
	return sum * negative
}

func (decode *decode) expect(kind byte) int {
	step := decode.pos
	for step < decode.end {
		if decode.buf[step] == kind {
			return step
		}
		step++
	}
	panic(bencodeErorr{error: errors.New("Invalid data: Missing delimiter ")})
}

func (decode *decode) next() interface{} {
	switch decode.buf[decode.pos] {
	case NumericStart:
		decode.step()
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
	num := decode.scanneBinaryLen(step)
	for i := decode.pos; i < step; i++ {
		decode.step()
	}
	decode.step()
	return num
}

func (decode *decode) convertBinary() []byte {
	step := decode.expect(StringDelim)
	l := decode.scanneBinaryLen(step)
	for i := decode.pos; i < (step + l); i++ {
		decode.step()
	}
	decode.step()
	return decode.buf[step+1 : decode.pos]
}
