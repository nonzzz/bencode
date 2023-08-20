package bencode

import "testing"

func TestEncodeNumeric(t *testing.T) {
	input := -42
	expected := "i-42e"
	encoded, _ := Encode(input)

	if string(encoded) != expected {
		t.Fatalf("%s != %s", encoded, expected)
	}
}

func TestEncodeString(t *testing.T) {
	input := "spam"
	expected := "4:spam"
	encoded, _ := Encode(input)

	if string(encoded) != expected {
		t.Fatalf("%s != %s", encoded, expected)
	}
}

func TestEncodeList(t *testing.T) {
	input := []interface{}{"spam", 42}
	expected := "l4:spami42ee"
	encoded, _ := Encode(input)

	if string(encoded) != expected {
		t.Fatalf("%s != %s", encoded, expected)
	}
}

func TestEncodeDirectory(t *testing.T) {
	input := map[string]interface{}{
		"foo": 42,
		"bar": "spam",
	}
	expected := "d3:bar4:spam3:fooi42ee"
	encoded, _ := Encode(input)

	if string(encoded) != expected {
		t.Fatalf("%s != %s", encoded, expected)
	}
}
