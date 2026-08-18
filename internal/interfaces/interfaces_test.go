package interfaces

import (
	"errors"
	"fmt"
	"reflect"
	"strings"
	"testing"
)

type textStringer string

func (s textStringer) String() string { return string(s) }

func TestUserString(t *testing.T) {
	tests := []struct {
		user User
		want string
	}{
		{User{ID: 1, Name: "Maria"}, "user#1:Maria"},
		{User{ID: 0, Name: "Zero"}, "user#0:Zero"},
		{User{ID: -1, Name: "Minus"}, "user#-1:Minus"},
		{User{ID: 42, Name: ""}, "user#42:"},
		{User{ID: 7, Name: "Go Dev"}, "user#7:Go Dev"},
		{User{ID: 99, Name: "Аня"}, "user#99:Аня"},
		{User{ID: 100, Name: "A-B"}, "user#100:A-B"},
		{User{ID: 5, Name: "  trim?  "}, "user#5:  trim?  "},
		{User{ID: 13, Name: "x:y"}, "user#13:x:y"},
		{User{ID: 1000, Name: "Long"}, "user#1000:Long"},
	}
	for _, tt := range tests {
		if got := tt.user.String(); got != tt.want {
			t.Fatalf("got %q want %q", got, tt.want)
		}
	}
}

func TestFormatOne(t *testing.T) {
	tests := []struct {
		in   fmt.Stringer
		want string
	}{
		{User{1, "A"}, "user#1:A"},
		{textStringer("x"), "x"},
		{textStringer(""), ""},
		{User{0, ""}, "user#0:"},
		{textStringer(" spaced "), " spaced "},
		{nil, ""},
		{User{-3, "M"}, "user#-3:M"},
		{textStringer("go"), "go"},
		{User{77, "Рус"}, "user#77:Рус"},
		{textStringer("line\n"), "line\n"},
	}
	for _, tt := range tests {
		if got := FormatOne(tt.in); got != tt.want {
			t.Fatalf("got %q want %q", got, tt.want)
		}
	}
}

func TestFormatMany(t *testing.T) {
	tests := []struct {
		in   []fmt.Stringer
		want []string
	}{
		{nil, []string{}},
		{[]fmt.Stringer{}, []string{}},
		{[]fmt.Stringer{User{1, "A"}}, []string{"user#1:A"}},
		{[]fmt.Stringer{textStringer("x"), textStringer("y")}, []string{"x", "y"}},
		{[]fmt.Stringer{nil, textStringer("x")}, []string{"x"}},
		{[]fmt.Stringer{User{1, "A"}, nil, User{2, "B"}}, []string{"user#1:A", "user#2:B"}},
		{[]fmt.Stringer{textStringer("")}, []string{""}},
		{[]fmt.Stringer{textStringer("go"), User{3, "C"}}, []string{"go", "user#3:C"}},
		{[]fmt.Stringer{textStringer("1"), textStringer("2"), textStringer("3")}, []string{"1", "2", "3"}},
		{[]fmt.Stringer{User{-1, "M"}}, []string{"user#-1:M"}},
	}
	for _, tt := range tests {
		if got := FormatMany(tt.in); !reflect.DeepEqual(got, tt.want) {
			t.Fatalf("got %#v want %#v", got, tt.want)
		}
	}
}

func TestEmailNotifierNotify(t *testing.T) {
	for _, tt := range []struct {
		n         EmailNotifier
		msg, want string
	}{
		{EmailNotifier{"a@b"}, "hi", "email:a@b:hi"},
		{EmailNotifier{""}, "hi", "email::hi"},
		{EmailNotifier{"x"}, "", "email:x:"},
		{EmailNotifier{"dev"}, "go", "email:dev:go"},
		{EmailNotifier{"team@example.com"}, "lesson", "email:team@example.com:lesson"},
		{EmailNotifier{"1"}, "2", "email:1:2"},
		{EmailNotifier{"рус"}, "привет", "email:рус:привет"},
		{EmailNotifier{"space addr"}, "space msg", "email:space addr:space msg"},
		{EmailNotifier{"a:b"}, "c:d", "email:a:b:c:d"},
		{EmailNotifier{"long"}, strings.Repeat("x", 3), "email:long:xxx"},
	} {
		if got := tt.n.Notify(tt.msg); got != tt.want {
			t.Fatalf("got %q want %q", got, tt.want)
		}
	}
}

func TestPushNotifierNotify(t *testing.T) {
	for _, tt := range []struct {
		n         PushNotifier
		msg, want string
	}{
		{PushNotifier{"ios"}, "hi", "push:ios:hi"},
		{PushNotifier{""}, "hi", "push::hi"},
		{PushNotifier{"x"}, "", "push:x:"},
		{PushNotifier{"android"}, "go", "push:android:go"},
		{PushNotifier{"dev-1"}, "lesson", "push:dev-1:lesson"},
		{PushNotifier{"1"}, "2", "push:1:2"},
		{PushNotifier{"рус"}, "привет", "push:рус:привет"},
		{PushNotifier{"space id"}, "space msg", "push:space id:space msg"},
		{PushNotifier{"a:b"}, "c:d", "push:a:b:c:d"},
		{PushNotifier{"long"}, strings.Repeat("x", 3), "push:long:xxx"},
	} {
		if got := tt.n.Notify(tt.msg); got != tt.want {
			t.Fatalf("got %q want %q", got, tt.want)
		}
	}
}

func TestSendNotification(t *testing.T) {
	for _, tt := range []struct {
		n         Notifier
		msg, want string
	}{
		{EmailNotifier{"a"}, "m", "email:a:m"},
		{PushNotifier{"p"}, "m", "push:p:m"},
		{nil, "m", ""},
		{EmailNotifier{""}, "", "email::"},
		{PushNotifier{""}, "", "push::"},
		{EmailNotifier{"team"}, "hello", "email:team:hello"},
		{PushNotifier{"dev"}, "hello", "push:dev:hello"},
		{EmailNotifier{"рус"}, "привет", "email:рус:привет"},
		{PushNotifier{"ios"}, "two words", "push:ios:two words"},
		{EmailNotifier{"x:y"}, "z", "email:x:y:z"},
	} {
		if got := SendNotification(tt.n, tt.msg); got != tt.want {
			t.Fatalf("got %q want %q", got, tt.want)
		}
	}
}

func TestSendBatch(t *testing.T) {
	tests := []struct {
		n        Notifier
		messages []string
		want     []string
	}{
		{nil, []string{"a"}, []string{}},
		{EmailNotifier{"e"}, nil, []string{}},
		{EmailNotifier{"e"}, []string{}, []string{}},
		{EmailNotifier{"e"}, []string{"a"}, []string{"email:e:a"}},
		{PushNotifier{"p"}, []string{"a"}, []string{"push:p:a"}},
		{EmailNotifier{"e"}, []string{"a", "b"}, []string{"email:e:a", "email:e:b"}},
		{PushNotifier{"p"}, []string{"", "b"}, []string{"push:p:", "push:p:b"}},
		{EmailNotifier{""}, []string{"x"}, []string{"email::x"}},
		{PushNotifier{"dev"}, []string{"1", "2", "3"}, []string{"push:dev:1", "push:dev:2", "push:dev:3"}},
		{EmailNotifier{"рус"}, []string{"пр"}, []string{"email:рус:пр"}},
	}
	for _, tt := range tests {
		if got := SendBatch(tt.n, tt.messages); !reflect.DeepEqual(got, tt.want) {
			t.Fatalf("got %#v want %#v", got, tt.want)
		}
	}
}

func TestStaticLoaderLoad(t *testing.T) {
	errA := errors.New("a")
	for _, tt := range []StaticLoader{
		{"", nil},
		{"x", nil},
		{"upper", nil},
		{"space value", nil},
		{"рус", nil},
		{"", errA},
		{"x", errA},
		{"line\n", nil},
		{"0", nil},
		{"!", nil}} {
		got, err := tt.Load()
		if got != tt.Value || !errors.Is(err, tt.Err) {
			t.Fatalf("got %q/%v want %q/%v", got, err, tt.Value, tt.Err)
		}
	}
}

func TestLoadUpper(t *testing.T) {
	errA := errors.New("load")
	for _, tt := range []struct {
		l       Loader
		want    string
		wantErr bool
	}{
		{StaticLoader{"go", nil}, "GO", false},
		{StaticLoader{"Go Lang", nil}, "GO LANG", false},
		{StaticLoader{"", nil}, "", false},
		{StaticLoader{"рус", nil}, "РУС", false},
		{StaticLoader{"123", nil}, "123", false},
		{StaticLoader{"a-b", nil}, "A-B", false},
		{StaticLoader{"x", errA}, "", true},
		{StaticLoader{"", errA}, "", true},
		{nil, "", true},
		{StaticLoader{" mixed ", nil}, " MIXED ", false},
	} {
		got, err := LoadUpper(tt.l)
		if got != tt.want || (err != nil) != tt.wantErr {
			t.Fatalf("got %q/%v want %q err=%t", got, err, tt.want, tt.wantErr)
		}
	}
}

func TestDescribe(t *testing.T) {
	errA := errors.New("boom")
	for _, tt := range []struct {
		in   any
		want string
	}{
		{nil, "nil"}, {"go", "string:go"},
		{"", "string:"},
		{User{1, "A"}, "user:user#1:A"},
		{EmailNotifier{"e"}, "notifier:email:e:ping"},
		{PushNotifier{"p"}, "notifier:push:p:ping"},
		{StaticLoader{"x", nil}, "loader:x"},
		{StaticLoader{"", errA}, "loader-error:boom"},
		{errA, "error:boom"},
		{42, "unknown:int"},
	} {
		if got := Describe(tt.in); got != tt.want {
			t.Fatalf("got %q want %q", got, tt.want)
		}
	}
}

func TestOnlyErrors(t *testing.T) {
	e1, e2 := errors.New("one"), errors.New("two")
	tests := []struct {
		in   []error
		want []string
	}{
		{nil, []string{}},
		{[]error{}, []string{}},
		{[]error{nil}, []string{}},
		{[]error{e1}, []string{"one"}},
		{[]error{nil, e1}, []string{"one"}},
		{[]error{e1, nil, e2}, []string{"one", "two"}},
		{[]error{fmt.Errorf("wrap: %w", e1)}, []string{"wrap: one"}},
		{[]error{errors.New("")}, []string{""}},
		{[]error{e1, e1}, []string{"one", "one"}},
		{[]error{nil, nil, e2}, []string{"two"}}}
	for _, tt := range tests {
		if got := OnlyErrors(tt.in); !reflect.DeepEqual(got, tt.want) {
			t.Fatalf("got %#v want %#v", got, tt.want)
		}
	}
}

func TestJoinStringers(t *testing.T) {
	tests := []struct {
		in        []fmt.Stringer
		sep, want string
	}{
		{nil, ",", ""},
		{[]fmt.Stringer{}, ",", ""},
		{[]fmt.Stringer{textStringer("a")}, ",", "a"},
		{[]fmt.Stringer{textStringer("a"), textStringer("b")}, ",", "a,b"},
		{[]fmt.Stringer{textStringer("a"), nil, textStringer("b")}, ",", "a,b"},
		{[]fmt.Stringer{textStringer("")}, ",", ""},
		{[]fmt.Stringer{textStringer("a"), textStringer("")}, "|", "a|"},
		{[]fmt.Stringer{User{1, "A"}, User{2, "B"}}, ";", "user#1:A;user#2:B"},
		{[]fmt.Stringer{textStringer("x"), textStringer("y"), textStringer("z")}, "", "xyz"},
		{[]fmt.Stringer{textStringer("go")}, "---", "go"}}
	for _, tt := range tests {
		if got := JoinStringers(tt.in, tt.sep); got != tt.want {
			t.Fatalf("got %q want %q", got, tt.want)
		}
	}
}

func TestFirstNonEmptyStringer(t *testing.T) {
	tests := []struct {
		in   []fmt.Stringer
		want string
	}{
		{nil, ""},
		{[]fmt.Stringer{}, ""},
		{[]fmt.Stringer{nil}, ""},
		{[]fmt.Stringer{textStringer("")}, ""},
		{[]fmt.Stringer{textStringer("a")}, "a"},
		{[]fmt.Stringer{textStringer(""), textStringer("b")}, "b"},
		{[]fmt.Stringer{nil, User{1, "A"}}, "user#1:A"},
		{[]fmt.Stringer{textStringer("  ")}, "  "},
		{[]fmt.Stringer{textStringer(""), textStringer("0")}, "0"},
		{[]fmt.Stringer{textStringer("first"), textStringer("second")}, "first"}}
	for _, tt := range tests {
		if got := FirstNonEmptyStringer(tt.in); got != tt.want {
			t.Fatalf("got %q want %q", got, tt.want)
		}
	}
}

func TestCountStringers(t *testing.T) {
	tests := []struct {
		in   []any
		want int
	}{
		{nil, 0},
		{[]any{}, 0},
		{[]any{User{1, "A"}}, 1},
		{[]any{"x"}, 0},
		{[]any{textStringer("x")}, 1},
		{[]any{User{}, textStringer("x"), 1}, 2},
		{[]any{nil, User{}}, 1},
		{[]any{errors.New("e")}, 0},
		{[]any{fmt.Stringer(textStringer("x"))}, 1},
		{[]any{User{}, User{}, User{}}, 3}}
	for _, tt := range tests {
		if got := CountStringers(tt.in); got != tt.want {
			t.Fatalf("got %d want %d", got, tt.want)
		}
	}
}

func TestNormalizeAndFormat(t *testing.T) {
	tests := []struct {
		in   fmt.Stringer
		want string
	}{
		{nil, "empty"},
		{textStringer(""), "empty"},
		{textStringer("  "), "empty"},
		{textStringer(" a "), "a"},
		{User{1, "A"}, "user#1:A"},
		{textStringer("go"), "go"},
		{textStringer("\nline\n"), "line"},
		{textStringer("0"), "0"},
		{textStringer(" a b "), "a b"},
		{User{0, ""}, "user#0:"}}
	for _, tt := range tests {
		if got := NormalizeAndFormat(tt.in); got != tt.want {
			t.Fatalf("got %q want %q", got, tt.want)
		}
	}
}

func TestExample(t *testing.T) {
	want := "user: user#7:Maria\n" + "email: email:team@example.com:hello\n" + "push: push:ios-1:hello\n" + "describe: user:user#8:Alex"
	if got := Example(); got != want {
		t.Fatalf("got %q want %q", got, want)
	}
}
