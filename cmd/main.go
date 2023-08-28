package main

import (
	"fmt"
	"os"

	"github.com/nonzzz/bencode"
)

var buf, _ = os.ReadFile("./test.torrent")

func main() {
	s, _ := bencode.Decode(buf)
	fmt.Println(s)
}
