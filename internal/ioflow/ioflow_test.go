package ioflow

import (
	"bytes"
	"errors"
	"io"
	"reflect"
	"strings"
	"testing"
)

type errReader struct{}

func (errReader) Read([]byte) (int, error) { return 0, errors.New("read error") }

type errWriter struct{}

func (errWriter) Write([]byte) (int, error) { return 0, errors.New("write error") }

func TestReadAllText(t *testing.T) {
	for _, tt := range []struct{ input, want string }{
		{"", ""},
		{"a", "a"},
		{"hello", "hello"},
		{"line\n", "line\n"},
		{"рус", "рус"},
		{"123", "123"},
		{" space ", " space "},
		{strings.Repeat("x", 10), strings.Repeat("x", 10)},
		{"a\nb", "a\nb"},
		{"!", "!"}} {
		got, err := ReadAllText(strings.NewReader(tt.input))
		if err != nil || got != tt.want {
			t.Fatalf("got %q/%v want %q", got, err, tt.want)
		}
	}
	if _, err := ReadAllText(errReader{}); err == nil {
		t.Fatal("expected read error")
	}
}

func TestCountBytes(t *testing.T) {
	for _, tt := range []struct {
		input string
		want  int
	}{
		{"", 0},
		{"a", 1},
		{"hello", 5},
		{"line\n", 5},
		{"рус", 6},
		{"123", 3},
		{" space ", 7},
		{strings.Repeat("x", 10), 10},
		{"a\nb", 3},
		{"!", 1}} {
		got, err := CountBytes(strings.NewReader(tt.input))
		if err != nil || got != tt.want {
			t.Fatalf("got %d/%v want %d", got, err, tt.want)
		}
	}
	if _, err := CountBytes(errReader{}); err == nil {
		t.Fatal("expected read error")
	}
}

func TestWriteString(t *testing.T) {
	for _, input := range []string{"", "a", "hello", "line\n", "рус", "123", " space ", strings.Repeat("x", 10), "a\nb", "!"} {
		var b bytes.Buffer
		n, err := WriteString(&b, input)
		if err != nil || n != len([]byte(input)) || b.String() != input {
			t.Fatalf("input %q n=%d value=%q err=%v", input, n, b.String(), err)
		}
	}
	if _, err := WriteString(errWriter{}, "x"); err == nil {
		t.Fatal("expected write error")
	}
}

func TestWriteLines(t *testing.T) {
	tests := []struct {
		lines []string
		want  string
	}{
		{nil, ""},
		{[]string{}, ""},
		{[]string{"a"}, "a\n"},
		{[]string{"a", "b"}, "a\nb\n"},
		{[]string{""}, "\n"},
		{[]string{"рус"}, "рус\n"},
		{[]string{"a b"}, "a b\n"},
		{[]string{"1", "2", "3"}, "1\n2\n3\n"},
		{[]string{"line\ninside"}, "line\ninside\n"},
		{[]string{"!"}, "!\n"}}
	for _, tt := range tests {
		var b bytes.Buffer
		if err := WriteLines(&b, tt.lines); err != nil || b.String() != tt.want {
			t.Fatalf("got %q/%v want %q", b.String(), err, tt.want)
		}
	}
	if err := WriteLines(errWriter{}, []string{"x"}); err == nil {
		t.Fatal("expected write error")
	}
}

func TestCopyAll(t *testing.T) {
	for _, input := range []string{"", "a", "hello", "line\n", "рус", "123", " space ", strings.Repeat("x", 10), "a\nb", "!"} {
		var b bytes.Buffer
		n, err := CopyAll(&b, strings.NewReader(input))
		if err != nil || n != int64(len([]byte(input))) || b.String() != input {
			t.Fatalf("input %q n=%d value=%q err=%v", input, n, b.String(), err)
		}
	}
	if _, err := CopyAll(io.Discard, errReader{}); err == nil {
		t.Fatal("expected read error")
	}
}

func TestCopyUpper(t *testing.T) {
	for _, tt := range []struct{ in, want string }{
		{"", ""},
		{"a", "A"},
		{"hello", "HELLO"},
		{"Go Lang", "GO LANG"},
		{"рус", "РУС"},
		{"123", "123"},
		{" space ", " SPACE "},
		{"a-b", "A-B"},
		{"a\nb", "A\nB"},
		{"!", "!"}} {
		var b bytes.Buffer
		n, err := CopyUpper(&b, strings.NewReader(tt.in))
		if err != nil || n != int64(len([]byte(tt.want))) || b.String() != tt.want {
			t.Fatalf("got n=%d value=%q err=%v want %q", n, b.String(), err, tt.want)
		}
	}
	if _, err := CopyUpper(io.Discard, errReader{}); err == nil {
		t.Fatal("expected read error")
	}
}

func TestScanLines(t *testing.T) {
	tests := []struct {
		in   string
		want []string
	}{
		{"", []string{}},
		{"a", []string{"a"}},
		{"a\n", []string{"a"}},
		{"a\nb", []string{"a", "b"}},
		{"\n", []string{""}},
		{"a\n\n", []string{"a", ""}},
		{"рус\nтекст", []string{"рус", "текст"}},
		{" space ", []string{" space "}},
		{"1\n2\n3", []string{"1", "2", "3"}},
		{"!", []string{"!"}}}
	for _, tt := range tests {
		got, err := ScanLines(strings.NewReader(tt.in))
		if err != nil || !reflect.DeepEqual(got, tt.want) {
			t.Fatalf("got %#v/%v want %#v", got, err, tt.want)
		}
	}
	if _, err := ScanLines(errReader{}); err == nil {
		t.Fatal("expected read error")
	}
}

func TestCountNonEmptyLines(t *testing.T) {
	for _, tt := range []struct {
		in   string
		want int
	}{
		{"", 0},
		{"a", 1},
		{"a\n", 1},
		{"a\nb", 2},
		{"\n", 0},
		{"a\n\n", 1},
		{"рус\nтекст", 2},
		{" space ", 1},
		{"1\n2\n3", 3},
		{" \n\t\n!", 1}} {
		got, err := CountNonEmptyLines(strings.NewReader(tt.in))
		if err != nil || got != tt.want {
			t.Fatalf("got %d/%v want %d", got, err, tt.want)
		}
	}
	if _, err := CountNonEmptyLines(errReader{}); err == nil {
		t.Fatal("expected read error")
	}
}

func TestFirstLine(t *testing.T) {
	for _, tt := range []struct{ in, want string }{
		{"", ""},
		{"a", "a"},
		{"a\n", "a"},
		{"a\nb", "a"},
		{"\n", ""},
		{"\nb", ""},
		{"рус\nтекст", "рус"},
		{" space ", " space "},
		{"1\n2\n3", "1"},
		{"!", "!"}} {
		got, err := FirstLine(strings.NewReader(tt.in))
		if err != nil || got != tt.want {
			t.Fatalf("got %q/%v want %q", got, err, tt.want)
		}
	}
	if _, err := FirstLine(errReader{}); err == nil {
		t.Fatal("expected read error")
	}
}

func TestReadCSVLike(t *testing.T) {
	tests := []struct {
		in   string
		want [][]string
	}{
		{"", [][]string{}},
		{"a", [][]string{{"a"}}},
		{"a,b", [][]string{{"a", "b"}}},
		{"a, b", [][]string{{"a", "b"}}},
		{"a,b\nc,d", [][]string{{"a", "b"},
			{"c", "d"}}},
		{",", [][]string{{"", ""}}},
		{" a , b ", [][]string{{"a", "b"}}},
		{"рус,текст", [][]string{{"рус", "текст"}}},
		{"1,2,3", [][]string{{"1", "2", "3"}}},
		{"!", [][]string{{"!"}}}}
	for _, tt := range tests {
		got, err := ReadCSVLike(strings.NewReader(tt.in))
		if err != nil || !reflect.DeepEqual(got, tt.want) {
			t.Fatalf("got %#v/%v want %#v", got, err, tt.want)
		}
	}
	if _, err := ReadCSVLike(errReader{}); err == nil {
		t.Fatal("expected read error")
	}
}

func TestLimitRead(t *testing.T) {
	for _, tt := range []struct {
		in    string
		limit int
		want  string
	}{
		{"abc", 0, ""},
		{"abc", 1, "a"},
		{"abc", 2, "ab"},
		{"abc", 3, "abc"},
		{"abc", 10, "abc"},
		{"", 10, ""},
		{"abc", -1, ""},
		{"рус", 2, "ру"[:2]},
		{"line\n", 4, "line"},
		{"!", 1, "!"}} {
		got, err := LimitRead(strings.NewReader(tt.in), tt.limit)
		if err != nil || got != tt.want {
			t.Fatalf("got %q/%v want %q", got, err, tt.want)
		}
	}
	if _, err := LimitRead(errReader{}, 10); err == nil {
		t.Fatal("expected read error")
	}
}

func TestRepeatToWriter(t *testing.T) {
	for _, tt := range []struct {
		text  string
		count int
		want  string
	}{
		{"a", 0, ""},
		{"a", 1, "a"},
		{"a", 2, "aa"},
		{"ab", 3, "ababab"},
		{"", 5, ""},
		{"рус", 2, "русрус"},
		{"x", -1, ""},
		{"!", 4, "!!!!"},
		{" ", 3, "   "},
		{"go", 1, "go"}} {
		var b bytes.Buffer
		if err := RepeatToWriter(&b, tt.text, tt.count); err != nil || b.String() != tt.want {
			t.Fatalf("got %q/%v want %q", b.String(), err, tt.want)
		}
	}
	if err := RepeatToWriter(errWriter{}, "x", 1); err == nil {
		t.Fatal("expected write error")
	}
}

func TestBufferReport(t *testing.T) {
	tests := []struct {
		lines []string
		want  string
	}{
		{nil, ""},
		{[]string{}, ""},
		{[]string{"a"}, "0:a\n"},
		{[]string{"a", "b"}, "0:a\n1:b\n"},
		{[]string{""}, "0:\n"},
		{[]string{"рус"}, "0:рус\n"},
		{[]string{"a b"}, "0:a b\n"},
		{[]string{"1", "2", "3"}, "0:1\n1:2\n2:3\n"},
		{[]string{"line\ninside"}, "0:line\ninside\n"},
		{[]string{"!"}, "0:!\n"}}
	for _, tt := range tests {
		if got := BufferReport(tt.lines); got != tt.want {
			t.Fatalf("got %q want %q", got, tt.want)
		}
	}
}

func TestReadAndTrim(t *testing.T) {
	for _, tt := range []struct{ in, want string }{
		{"", ""},
		{" a ", "a"},
		{"\nhello\n", "hello"},
		{"рус ", "рус"},
		{" 123 ", "123"},
		{" space inside ", "space inside"},
		{"\tgo\t", "go"},
		{"!", "!"},
		{" a\nb ", "a\nb"},
		{"   ", ""}} {
		got, err := ReadAndTrim(strings.NewReader(tt.in))
		if err != nil || got != tt.want {
			t.Fatalf("got %q/%v want %q", got, err, tt.want)
		}
	}
	if _, err := ReadAndTrim(errReader{}); err == nil {
		t.Fatal("expected read error")
	}
}

func TestWriteKeyValues(t *testing.T) {
	tests := []struct {
		values map[string]string
		want   string
	}{
		{nil, ""},
		{map[string]string{}, ""},
		{map[string]string{"a": "1"}, "a=1\n"},
		{map[string]string{"b": "2", "a": "1"}, "a=1\nb=2\n"},
		{map[string]string{"": "x"}, "=x\n"},
		{map[string]string{"рус": "текст"}, "рус=текст\n"},
		{map[string]string{"a": ""}, "a=\n"},
		{map[string]string{"10": "x", "2": "y"}, "10=x\n2=y\n"},
		{map[string]string{"a": "1", "b": "", "c": "3"}, "a=1\nb=\nc=3\n"},
		{map[string]string{"!": "?"}, "!=?\n"}}
	for _, tt := range tests {
		var b bytes.Buffer
		if err := WriteKeyValues(&b, tt.values); err != nil || b.String() != tt.want {
			t.Fatalf("got %q/%v want %q", b.String(), err, tt.want)
		}
	}
	if err := WriteKeyValues(errWriter{}, map[string]string{"a": "1"}); err == nil {
		t.Fatal("expected write error")
	}
}

func TestExample(t *testing.T) {
	want := "read: hello\nnon-empty: 0:2"
	if got := strings.TrimSpace(Example()); got != want {
		t.Fatalf("got %q want %q", got, want)
	}
}
