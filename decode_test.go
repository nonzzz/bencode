package bencode

import (
	"reflect"
	"testing"
)

func TestDecodeNumeric(t *testing.T) {
	input := "i-42e"
	expected := -42
	decoded, _ := Decode([]byte(input))

	value, _ := decoded.(int)
	if value != expected {
		t.Fatalf("%d != %d", value, expected)
	}
}

func TestDecodeStrin(t *testing.T) {
	input := "4:spam"
	expected := "spam"
	decoded, _ := Decode([]byte(input))
	value, _ := decoded.([]byte)
	if string(value) != expected {
		t.Fatalf("%s != %s", string(value), expected)
	}
}

func TestDecodeList(t *testing.T) {
	input := "l4:spami42ee"
	expected1 := "spam"
	var expected2 int64 = 42
	decoded, _ := Decode([]byte(input))
	value, _ := decoded.([]interface{})
	n1 := string(reflect.ValueOf(value[0]).Bytes())
	n2 := reflect.ValueOf(value[1]).Int()
	if n1 != expected1 {
		t.Fatalf("%s != %s", n1, expected1)
	}
	if n2 != expected2 {
		t.Fatalf("%d != %d", n2, expected2)
	}
}

func TestDecodeDirectory(t *testing.T) {
	input := "d3:bar4:spam3:fooi42ee"
	decoded, _ := Decode([]byte(input))
	value, _ := decoded.(map[string]interface{})
	expected1 := "spam"
	var expected2 int64 = 42
	n1 := string(reflect.ValueOf(value["bar"]).Bytes())
	n2 := reflect.ValueOf(value["foo"]).Int()
	if n1 != expected1 {
		t.Fatalf("%s != %s", n1, expected1)
	}
	if n2 != expected2 {
		t.Fatalf("%d != %d", n2, expected2)
	}
}
