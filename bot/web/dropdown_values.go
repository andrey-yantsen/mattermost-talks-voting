package web

import "time"

func getViewingDays() []dropdownValueInt {
	return []dropdownValueInt{
		{1, "Monday"},
		{2, "Tuesday"},
		{3, "Wednesday"},
		{4, "Thursday"},
		{5, "Friday"},
	}
}

func getViewingTimes() (ret []dropdownValue) {
	ret = make([]dropdownValue, 0)

	start, _ := time.Parse("15:04", "10:00")

	// Iterate from 10:00 to 18:00
	for d := start; d.Hour() < 18 || d.Minute() == 0; d = d.Add(30 * time.Minute) {
		formattedTime := d.Format("15:04")
		ret = append(ret, dropdownValue{formattedTime, formattedTime})
	}

	return
}

func getTimezones() []dropdownValue {
	return []dropdownValue{{"Europe/London", "London"}, {"Europe/Moscow", "Moscow"}}
}
