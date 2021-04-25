package utils

import "sort"

func InArrayString(target string, str_array []string) bool {
	sort.Strings(str_array)
	index := sort.SearchStrings(str_array, target)
	if index < len(str_array) && str_array[index] == target {
		return true
	}
	return false
}

func InArryInt(target int,int64_array  []int) bool {
	sort.Ints(int64_array)
	index := sort.SearchInts(int64_array, target)
	if index < len(int64_array) && int64_array[index] == target {
		return true
	}
	return false
}
