package bencode

import (
	"os"
	"testing"
)

var buf, _ = os.ReadFile("./test.torrent")
var decoded, _ = Decode(buf)

func BenchmarkEncode(b *testing.B) {

	for i := 0; i < b.N; i++ {
		_, _ = Encode(decoded)
	}
}

func BenchmarkDecode(b *testing.B) {

	for i := 0; i < b.N; i++ {
		_, _ = Decode(buf)
	}

}
