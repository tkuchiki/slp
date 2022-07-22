package helper

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

func CompileQueryMatchingGroups(groups []string) ([]*regexp.Regexp, error) {
	queryMatchingGroups := make([]*regexp.Regexp, 0, len(groups))
	for _, pattern := range groups {
		re, err := regexp.Compile(pattern)
		if err != nil {
			return nil, err
		}
		queryMatchingGroups = append(queryMatchingGroups, re)
	}

	return queryMatchingGroups, nil
}

func StringToFloat64(val string) (float64, error) {
	return strconv.ParseFloat(val, 64)
}

func StringToInt(val string) (int, error) {
	return strconv.Atoi(val)
}

func StringToInt64(val string) (int64, error) {
	return strconv.ParseInt(val, 10, 64)
}

func SplitCSV(val string) []string {
	strs := strings.Split(val, ",")
	if len(strs) == 1 && strs[0] == "" {
		return []string{}
	}

	trimedStrs := make([]string, 0, len(strs))

	for _, s := range strs {
		trimedStrs = append(trimedStrs, strings.Trim(s, " "))
	}

	return trimedStrs
}

func SplitCSVIntoInts(val string) ([]int, error) {
	strs := strings.Split(val, ",")
	if len(strs) == 1 && strs[0] == "" {
		return []int{}, nil
	}

	trimedInts := make([]int, 0, len(strs))

	for _, s := range strs {
		i, err := strconv.Atoi(strings.Trim(s, " "))
		if err != nil {
			return []int{}, err
		}
		trimedInts = append(trimedInts, i)
	}

	for _, i := range trimedInts {
		if i < 1 && i > 100 {
			return []int{}, fmt.Errorf(``)
		}
	}

	return trimedInts, nil
}

func ValidatePercentiles(percentiles []int) error {
	if len(percentiles) == 0 {
		return nil
	}

	for _, i := range percentiles {
		if i < 0 && i > 100 {
			return fmt.Errorf(`percentiles allowed 0 to 100`)
		}
	}

	return nil
}
