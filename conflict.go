package main

// Weeks 表示一个教学班在哪些周上课。
type Weeks string

const (
	AllWeeks  Weeks = "all"  // 每周
	OddWeeks  Weeks = "odd"  // 单周
	EvenWeeks Weeks = "even" // 双周
)

// Section 是一个教学班在一周内的一个上课时段。
type Section struct {
	ID      string `json:"id"`
	Weekday int    `json:"weekday"` // 1=周一 … 7=周日
	Start   int    `json:"start"`   // 起始节次，含
	End     int    `json:"end"`     // 结束节次，含
	Weeks   Weeks  `json:"weeks"`
}

// Valid 检查时段本身是否合法。
func (s Section) Valid() bool {
	if s.Weekday < 1 || s.Weekday > 7 {
		return false
	}
	if s.Start < 1 || s.End < s.Start {
		return false
	}
	switch s.Weeks {
	case AllWeeks, OddWeeks, EvenWeeks:
		return true
	}
	return false
}

// weeksOverlap 判断两个班的上课周次是否有交集：
// 任一方每周上课就一定有交集；否则只有单周对单周、双周对双周才有交集。
func weeksOverlap(a, b Weeks) bool {
	if a == AllWeeks || b == AllWeeks {
		return true
	}
	return a == b
}

// Conflicts 判断两个教学班的上课时间是否冲突：
// 同一天、周次有交集，且节次区间有重叠。
func Conflicts(a, b Section) bool {
	if a.Weekday != b.Weekday {
		return false
	}
	if !weeksOverlap(a.Weeks, b.Weeks) {
		return false
	}
	return a.Start <= b.End && b.Start <= a.End
}
