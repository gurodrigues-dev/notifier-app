package slicecommon

import (
	"unicode"

	"github.com/gurodrigues-dev/notifier-app/pkg/stringcommon"
)

func Partition(input []string) (nums, words []string) {
	for _, s := range input {
		if stringcommon.Empty(s) {
			continue
		}
		isNum := true
		for _, r := range s {
			if !unicode.IsDigit(r) {
				isNum = false
				break
			}
		}
		if isNum {
			nums = append(nums, s)
		} else {
			words = append(words, s)
		}
	}
	return
}

func Contains[T comparable](list []T, target T) bool {
	for _, item := range list {
		if item == target {
			return true
		}
	}
	return false
}
