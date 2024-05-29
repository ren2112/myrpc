package Provider

import "sort"

type SortServiceImpl struct {
}

func (s *SortServiceImpl) QuickSort(arr []int) []int {
	sort.Ints(arr)
	return arr
}
