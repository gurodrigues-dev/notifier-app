package slicecommon

import "unicode"

func Partition(input []string) (nums, words []string) {
	for _, s := range input {
		isNum := true
		for _, r := range s {
			if !unicode.IsDigit(r) {
				isNum = false
				break
			}
		}
		if isNum {
			nums = append(nums, s)
		}

		words = append(words, s)
	}
	return
}
