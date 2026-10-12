package main

import "testing"

func TestConflicts(t *testing.T) {
	mon := func(start, end int, w Weeks) Section {
		return Section{Weekday: 1, Start: start, End: end, Weeks: w}
	}
	cases := []struct {
		name string
		a, b Section
		want bool
	}{
		{"同一时段每周上课", mon(1, 2, AllWeeks), mon(1, 2, AllWeeks), true},
		{"节次部分重叠", mon(1, 3, AllWeeks), mon(3, 4, AllWeeks), true},
		{"节次首尾相接不冲突", mon(1, 2, AllWeeks), mon(3, 4, AllWeeks), false},
		{"不同星期不冲突", mon(1, 2, AllWeeks), Section{Weekday: 2, Start: 1, End: 2, Weeks: AllWeeks}, false},
		{"单周对双周不冲突", mon(1, 2, OddWeeks), mon(1, 2, EvenWeeks), false},
		{"单周对每周冲突", mon(1, 2, OddWeeks), mon(1, 2, AllWeeks), true},
		{"单周对单周冲突", mon(1, 2, OddWeeks), mon(2, 3, OddWeeks), true},
		{"双周对每周冲突", mon(5, 6, EvenWeeks), mon(5, 6, AllWeeks), true},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			if got := Conflicts(c.a, c.b); got != c.want {
				t.Errorf("Conflicts(%+v, %+v) = %v, want %v", c.a, c.b, got, c.want)
			}
			// 冲突关系应当对称
			if got := Conflicts(c.b, c.a); got != c.want {
				t.Errorf("Conflicts 不对称：Conflicts(b, a) = %v, want %v", got, c.want)
			}
		})
	}
}

func TestValid(t *testing.T) {
	good := Section{Weekday: 3, Start: 1, End: 2, Weeks: EvenWeeks}
	if !good.Valid() {
		t.Fatalf("%+v 应当合法", good)
	}
	bad := []Section{
		{Weekday: 0, Start: 1, End: 2, Weeks: AllWeeks},
		{Weekday: 8, Start: 1, End: 2, Weeks: AllWeeks},
		{Weekday: 1, Start: 0, End: 2, Weeks: AllWeeks},
		{Weekday: 1, Start: 3, End: 2, Weeks: AllWeeks},
		{Weekday: 1, Start: 1, End: 2, Weeks: "weekly"},
	}
	for _, s := range bad {
		if s.Valid() {
			t.Errorf("%+v 应当不合法", s)
		}
	}
}
