package utils

import "time"

func GetSchoolYear() int {
	now := time.Now()
	month := int(now.Month())
	if 1 <= month && month <= 3 {
		return now.Year() - 1
	}
	return now.Year()
}

func GetUniqueString(str []string) []string {
	m := make(map[string]bool)
	uniq := []string{}

	for _, ele := range str {
		if !m[ele] {
			m[ele] = true
			uniq = append(uniq, ele)
		}
	}
	return uniq
}
