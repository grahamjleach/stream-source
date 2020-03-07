package domain

type ScheduleItems []ScheduleItem

func (is ScheduleItems) Len() int {
	return len(is)
}

func (is ScheduleItems) Less(i, j int) bool {
	return is[i].StartOffset() < is[j].StartOffset()
}

func (is ScheduleItems) Swap(i, j int) {
	is[i], is[j] = is[j], is[i]
}
