package methodsets

import "testing"

func TestNewCounter(t *testing.T) {
	for _, tt := range []struct {
		name  string
		value int
	}{{"orders", 1}, {"", 0}, {"x", -1}, {"рус", 5}, {"long name", 10}, {"a:b", 7}, {"zero", 0}, {"neg", -100}, {"big", 1000}, {"space ", 3}} {
		got := NewCounter(tt.name, tt.value)
		if got.Name != tt.name || got.Value != tt.value {
			t.Fatalf("got %#v", got)
		}
	}
}

func TestSnapshot(t *testing.T) {
	for _, value := range []int{0, 1, -1, 10, 42, 100, -100, 7, 999, 5} {
		if got := (Counter{Value: value}).Snapshot(); got != value {
			t.Fatalf("got %d want %d", got, value)
		}
	}
}

func TestLabel(t *testing.T) {
	for _, tt := range []struct {
		c    Counter
		want string
	}{{Counter{"orders", 1}, "orders=1"}, {Counter{"", 0}, "=0"}, {Counter{"x", -1}, "x=-1"}, {Counter{"рус", 5}, "рус=5"}, {Counter{"long name", 10}, "long name=10"}, {Counter{"a:b", 7}, "a:b=7"}, {Counter{"zero", 0}, "zero=0"}, {Counter{"neg", -100}, "neg=-100"}, {Counter{"big", 1000}, "big=1000"}, {Counter{"space ", 3}, "space =3"}} {
		if got := tt.c.Label(); got != tt.want {
			t.Fatalf("got %q want %q", got, tt.want)
		}
	}
}

func TestAdd(t *testing.T) {
	for _, tt := range []struct{ start, delta, want int }{{0, 1, 1}, {1, 1, 2}, {1, -1, 0}, {10, 5, 15}, {-5, 5, 0}, {100, -50, 50}, {0, 0, 0}, {7, 3, 10}, {-10, -5, -15}, {1000, 1, 1001}} {
		c := Counter{Value: tt.start}
		c.Add(tt.delta)
		if c.Value != tt.want {
			t.Fatalf("got %d want %d", c.Value, tt.want)
		}
	}
	var c *Counter
	c.Add(10)
}

func TestReset(t *testing.T) {
	for _, start := range []int{0, 1, -1, 10, 42, 100, -100, 7, 999, 5} {
		c := Counter{Value: start}
		c.Reset()
		if c.Value != 0 {
			t.Fatalf("got %d", c.Value)
		}
	}
	var c *Counter
	c.Reset()
}

func TestUseSnapshoter(t *testing.T) {
	cases := []Snapshoter{Counter{Value: 1}, &Counter{Value: 2}, Counter{Value: 0}, &Counter{Value: -1}, Counter{Value: 10}, &Counter{Value: 42}, Counter{Value: -100}, &Counter{Value: 7}, Counter{Value: 999}, nil}
	wants := []int{1, 2, 0, -1, 10, 42, -100, 7, 999, 0}
	for i, c := range cases {
		if got := UseSnapshoter(c); got != wants[i] {
			t.Fatalf("case %d got %d want %d", i, got, wants[i])
		}
	}
}

func TestUseLabeler(t *testing.T) {
	cases := []Labeler{Counter{Name: "a", Value: 1}, &Counter{Name: "b", Value: 2}, Counter{}, &Counter{Name: "x", Value: -1}, Counter{Name: "рус", Value: 5}, &Counter{Name: "long", Value: 10}, Counter{Name: "zero", Value: 0}, &Counter{Name: "neg", Value: -10}, Counter{Name: "big", Value: 999}, nil}
	wants := []string{"a=1", "b=2", "=0", "x=-1", "рус=5", "long=10", "zero=0", "neg=-10", "big=999", ""}
	for i, c := range cases {
		if got := UseLabeler(c); got != wants[i] {
			t.Fatalf("case %d got %q want %q", i, got, wants[i])
		}
	}
}

func TestUseAdder(t *testing.T) {
	for _, tt := range []struct{ start, delta, want int }{{0, 1, 1}, {1, 1, 2}, {1, -1, 0}, {10, 5, 15}, {-5, 5, 0}, {100, -50, 50}, {0, 0, 0}, {7, 3, 10}, {-10, -5, -15}, {1000, 1, 1001}} {
		c := &Counter{Value: tt.start}
		if got := UseAdder(c, tt.delta); got != tt.want || c.Value != tt.want {
			t.Fatalf("got %d value %d want %d", got, c.Value, tt.want)
		}
	}
}

func TestIsSnapshoter(t *testing.T) {
	for _, tt := range []struct {
		v    any
		want bool
	}{{Counter{}, true}, {&Counter{}, true}, {Profile{}, false}, {&Profile{}, false}, {nil, false}, {42, false}, {"x", false}, {struct{}{}, false}, {Counter{Value: 1}, true}, {(*Counter)(nil), true}} {
		if got := IsSnapshoter(tt.v); got != tt.want {
			t.Fatalf("got %t want %t for %T", got, tt.want, tt.v)
		}
	}
}

func TestIsAdder(t *testing.T) {
	for _, tt := range []struct {
		v    any
		want bool
	}{{Counter{}, false}, {&Counter{}, true}, {Profile{}, false}, {&Profile{}, false}, {nil, false}, {42, false}, {"x", false}, {struct{}{}, false}, {Counter{Value: 1}, false}, {(*Counter)(nil), true}} {
		if got := IsAdder(tt.v); got != tt.want {
			t.Fatalf("got %t want %t for %T", got, tt.want, tt.v)
		}
	}
}

func TestCloneAndAdd(t *testing.T) {
	for _, tt := range []struct{ start, delta, want int }{{0, 1, 1}, {1, 1, 2}, {1, -1, 0}, {10, 5, 15}, {-5, 5, 0}, {100, -50, 50}, {0, 0, 0}, {7, 3, 10}, {-10, -5, -15}, {1000, 1, 1001}} {
		original := Counter{Name: "c", Value: tt.start}
		got := CloneAndAdd(original, tt.delta)
		if got.Value != tt.want || original.Value != tt.start || got.Name != "c" {
			t.Fatalf("got %#v original %#v", got, original)
		}
	}
}

func TestAddInPlace(t *testing.T) {
	for _, tt := range []struct{ start, delta, want int }{{0, 1, 1}, {1, 1, 2}, {1, -1, 0}, {10, 5, 15}, {-5, 5, 0}, {100, -50, 50}, {0, 0, 0}, {7, 3, 10}, {-10, -5, -15}, {1000, 1, 1001}} {
		c := &Counter{Value: tt.start}
		if got := AddInPlace(c, tt.delta); got != tt.want || c.Value != tt.want {
			t.Fatalf("got %d value %d want %d", got, c.Value, tt.want)
		}
	}
	if got := AddInPlace(nil, 10); got != 0 {
		t.Fatalf("nil got %d", got)
	}
}

func TestProfileDisplay(t *testing.T) {
	for _, tt := range []struct {
		p    Profile
		want string
	}{{Profile{"Maria", 25}, "Maria(25)"}, {Profile{"", 0}, "(0)"}, {Profile{"A", -1}, "A(-1)"}, {Profile{"Рус", 5}, "Рус(5)"}, {Profile{"Long Name", 10}, "Long Name(10)"}, {Profile{"a:b", 7}, "a:b(7)"}, {Profile{"zero", 0}, "zero(0)"}, {Profile{"neg", -100}, "neg(-100)"}, {Profile{"big", 1000}, "big(1000)"}, {Profile{"space ", 3}, "space (3)"}} {
		if got := tt.p.Display(); got != tt.want {
			t.Fatalf("got %q want %q", got, tt.want)
		}
	}
}

func TestProfileRename(t *testing.T) {
	for _, tt := range []struct{ start, name string }{{"A", "B"}, {"", "B"}, {"A", ""}, {"Рус", "Имя"}, {"Long", "Longer"}, {"a:b", "c:d"}, {"zero", "0"}, {"same", "same"}, {"space", " space "}, {"old", "new"}} {
		p := &Profile{Name: tt.start}
		p.Rename(tt.name)
		if p.Name != tt.name {
			t.Fatalf("got %q want %q", p.Name, tt.name)
		}
	}
	var p *Profile
	p.Rename("ignored")
}

func TestIsRenamer(t *testing.T) {
	for _, tt := range []struct {
		v    any
		want bool
	}{{Profile{}, false}, {&Profile{}, true}, {Counter{}, false}, {&Counter{}, false}, {nil, false}, {42, false}, {"x", false}, {struct{}{}, false}, {&Profile{Name: "A"}, true}, {(*Profile)(nil), true}} {
		if got := IsRenamer(tt.v); got != tt.want {
			t.Fatalf("got %t want %t for %T", got, tt.want, tt.v)
		}
	}
}

func TestExample(t *testing.T) {
	want := "value: orders=10\n" + "pointer before: orders=10\n" + "pointer after: orders=15\n" + "profile: Masha(25)"
	if got := Example(); got != want {
		t.Fatalf("got %q want %q", got, want)
	}
}
