package version

import (
	"strconv"
	"strings"
)

type Version string

func (me Version) compareTo(other Version) int {
	var (
		meTab    = strings.Split(string(me), ".")
		otherTab = strings.Split(string(other), ".")
	)

	max := len(meTab)
	if len(otherTab) > max {
		max = len(otherTab)
	}
	for i := 0; i < max; i++ {
		var meInt, otherInt int

		if len(meTab) > i {
			meInt, _ = strconv.Atoi(meTab[i])
		}
		if len(otherTab) > i {
			otherInt, _ = strconv.Atoi(otherTab[i])
		}
		if meInt > otherInt {
			return 1
		}
		if otherInt > meInt {
			return -1
		}
	}
	return 0
}

func (me Version) LessThan(other Version) bool {
	return me.compareTo(other) == -1
}

func (me Version) LessThanOrEqualTo(other Version) bool {
	return me.compareTo(other) <= 0
}

func (me Version) GreaterThan(other Version) bool {
	return me.compareTo(other) == 1
}

func (me Version) GreaterThanOrEqualTo(other Version) bool {
	return me.compareTo(other) >= 0
}

func (me Version) Equal(other Version) bool {
	return me.compareTo(other) == 0
}
