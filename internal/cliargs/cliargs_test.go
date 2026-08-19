package cliargs

import (
	"errors"
	"reflect"
	"strings"
	"testing"
)

func TestUserArgs(t *testing.T) {
	tests := []struct {
		name string
		all  []string
		want []string
	}{
		{name: "nil", all: nil, want: []string{}},
		{name: "empty", all: []string{}, want: []string{}},
		{name: "program only", all: []string{"args"}, want: []string{}},
		{name: "one user argument", all: []string{"args", "Maria"}, want: []string{"Maria"}},
		{name: "two user arguments", all: []string{"args", "Maria", "2"}, want: []string{"Maria", "2"}},
		{name: "keeps blanks", all: []string{"args", ""}, want: []string{""}},
		{name: "keeps spaces", all: []string{"args", "Maria Petrova"}, want: []string{"Maria Petrova"}},
		{name: "keeps unicode", all: []string{"args", "Ринат"}, want: []string{"Ринат"}},
		{name: "keeps flags as values", all: []string{"args", "--name"}, want: []string{"--name"}},
		{name: "returns independent slice", all: []string{"args", "Maria", "2"}, want: []string{"Maria", "2"}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := UserArgs(tt.all)
			if !reflect.DeepEqual(got, tt.want) {
				t.Fatalf("UserArgs(%q) = %q, want %q", tt.all, got, tt.want)
			}
			if tt.name == "returns independent slice" && len(got) > 0 {
				got[0] = "changed"
				if tt.all[1] == "changed" {
					t.Fatal("UserArgs must return a copy, not a view of os.Args")
				}
			}
		})
	}
}

func TestParse(t *testing.T) {
	tests := []struct {
		name    string
		args    []string
		want    Options
		wantErr bool
	}{
		{name: "name only", args: []string{"Maria"}, want: Options{Name: "Maria", Repeat: 1}},
		{name: "name and repeat", args: []string{"Maria", "2"}, want: Options{Name: "Maria", Repeat: 2}},
		{name: "trims name", args: []string{"  Maria  "}, want: Options{Name: "Maria", Repeat: 1}},
		{name: "unicode name", args: []string{"Ринат", "3"}, want: Options{Name: "Ринат", Repeat: 3}},
		{name: "name with space", args: []string{"Maria Petrova"}, want: Options{Name: "Maria Petrova", Repeat: 1}},
		{name: "missing name", args: nil, wantErr: true},
		{name: "blank name", args: []string{"   "}, wantErr: true},
		{name: "too many arguments", args: []string{"Maria", "2", "extra"}, wantErr: true},
		{name: "non numeric repeat", args: []string{"Maria", "two"}, wantErr: true},
		{name: "non positive repeat", args: []string{"Maria", "0"}, wantErr: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Parse(tt.args)
			if tt.wantErr {
				if !errors.Is(err, ErrUsage) {
					t.Fatalf("Parse(%q) error = %v, want ErrUsage", tt.args, err)
				}
				return
			}
			if err != nil {
				t.Fatalf("Parse(%q) error = %v", tt.args, err)
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Fatalf("Parse(%q) = %#v, want %#v", tt.args, got, tt.want)
			}
		})
	}
}

func TestRun(t *testing.T) {
	tests := []struct {
		name    string
		args    []string
		want    string
		wantErr bool
	}{
		{name: "default repeat", args: []string{"Maria"}, want: "hello, Maria\n"},
		{name: "repeat twice", args: []string{"Maria", "2"}, want: "hello, Maria\nhello, Maria\n"},
		{name: "repeat three times", args: []string{"Go", "3"}, want: "hello, Go\nhello, Go\nhello, Go\n"},
		{name: "unicode", args: []string{"Ринат"}, want: "hello, Ринат\n"},
		{name: "space in one argument", args: []string{"Maria Petrova"}, want: "hello, Maria Petrova\n"},
		{name: "missing", args: nil, wantErr: true},
		{name: "blank", args: []string{""}, wantErr: true},
		{name: "bad repeat", args: []string{"Maria", "x"}, wantErr: true},
		{name: "zero repeat", args: []string{"Maria", "0"}, wantErr: true},
		{name: "extra argument", args: []string{"Maria", "2", "extra"}, wantErr: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var out strings.Builder
			err := Run(tt.args, &out)
			if tt.wantErr {
				if !errors.Is(err, ErrUsage) {
					t.Fatalf("Run(%q) error = %v, want ErrUsage", tt.args, err)
				}
				return
			}
			if err != nil {
				t.Fatalf("Run(%q) error = %v", tt.args, err)
			}
			if got := out.String(); got != tt.want {
				t.Fatalf("Run(%q) output = %q, want %q", tt.args, got, tt.want)
			}
		})
	}
}

func TestExample(t *testing.T) {
	if got, want := Example(), "hello, Maria\nhello, Maria\n"; got != want {
		t.Fatalf("Example() = %q, want %q", got, want)
	}
}
