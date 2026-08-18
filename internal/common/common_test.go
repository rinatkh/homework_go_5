package common

import (
	"errors"
	"os"
	"path/filepath"
	"testing"
)

func TestReadFileOrDefault(t *testing.T) {
	dir := t.TempDir()
	for _, tt := range []struct {
		name, content, def, want string
		create                   bool
	}{
		{"a.txt", "hello", "def", "hello", true}, {"empty.txt", "", "def", "", true}, {"missing.txt", "", "def", "def", false}, {"rus.txt", "рус", "def", "рус", true}, {"lines.txt", "a\nb", "def", "a\nb", true}, {"space.txt", " space ", "def", " space ", true}, {"zero.txt", "0", "def", "0", true}, {"bang.txt", "!", "def", "!", true}, {"missing2.txt", "", "", "", false}, {"long.txt", "abcdef", "x", "abcdef", true},
	} {
		path := filepath.Join(dir, tt.name)
		if tt.create {
			if err := os.WriteFile(path, []byte(tt.content), 0o644); err != nil {
				t.Fatal(err)
			}
		}
		got, err := ReadFileOrDefault(path, tt.def)
		if err != nil || got != tt.want {
			t.Fatalf("%s got %q/%v want %q", tt.name, got, err, tt.want)
		}
	}
	_, err := ReadFileOrDefault(dir, "def")
	if err == nil {
		t.Fatal("reading directory must return error")
	}
}

func TestWriteFileLines(t *testing.T) {
	dir := t.TempDir()
	for i, tt := range []struct {
		lines []string
		want  string
	}{{nil, ""}, {[]string{}, ""}, {[]string{"a"}, "a\n"}, {[]string{"a", "b"}, "a\nb\n"}, {[]string{""}, "\n"}, {[]string{"рус"}, "рус\n"}, {[]string{"a b"}, "a b\n"}, {[]string{"1", "2", "3"}, "1\n2\n3\n"}, {[]string{"line\ninside"}, "line\ninside\n"}, {[]string{"!"}, "!\n"}} {
		path := filepath.Join(dir, "file"+string(rune('a'+i))+".txt")
		if err := WriteFileLines(path, tt.lines); err != nil {
			t.Fatal(err)
		}
		data, err := os.ReadFile(path)
		if err != nil || string(data) != tt.want {
			t.Fatalf("got %q/%v want %q", string(data), err, tt.want)
		}
	}
	if err := WriteFileLines(filepath.Join(dir, "missing", "x.txt"), []string{"x"}); err == nil {
		t.Fatal("expected write error")
	}
}

func TestCountFileLines(t *testing.T) {
	dir := t.TempDir()
	for _, tt := range []struct {
		name, content string
		want          int
	}{{"empty.txt", "", 0}, {"one.txt", "a", 1}, {"one-nl.txt", "a\n", 1}, {"two.txt", "a\nb", 2}, {"blank.txt", "\n", 1}, {"two-blanks.txt", "a\n\n", 2}, {"rus.txt", "рус\nтекст", 2}, {"space.txt", " space ", 1}, {"three.txt", "1\n2\n3", 3}, {"bang.txt", "!", 1}} {
		path := filepath.Join(dir, tt.name)
		if err := os.WriteFile(path, []byte(tt.content), 0o644); err != nil {
			t.Fatal(err)
		}
		got, err := CountFileLines(path)
		if err != nil || got != tt.want {
			t.Fatalf("%s got %d/%v want %d", tt.name, got, err, tt.want)
		}
	}
	_, err := CountFileLines(filepath.Join(dir, "missing.txt"))
	if !errors.Is(err, os.ErrNotExist) {
		t.Fatalf("expected os.ErrNotExist, got %v", err)
	}
}

func TestExample(t *testing.T) {
	if got := Example(); got != "common file tasks" {
		t.Fatalf("got %q", got)
	}
}
