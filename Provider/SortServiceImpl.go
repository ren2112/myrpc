package Provider

type SortServiceImpl struct {
}

func (s *SortServiceImpl) QuickSort(arr []int) []int {
	//if len(arr) <= 1 {
	//	return arr
	//}
	//
	//pivotIndex := len(arr) / 2
	//pivot := arr[pivotIndex]
	//// 移动pivot到数组开头
	//arr[0], arr[pivotIndex] = arr[pivotIndex], arr[0]
	//
	//// partition操作，将小于pivot的元素放到左边，大于pivot的放到右边
	//left, right := 1, len(arr)-1
	//for left <= right {
	//	for left <= right && arr[left] < pivot {
	//		left++
	//	}
	//	for left <= right && arr[right] > pivot {
	//		right--
	//	}
	//	if left <= right {
	//		arr[left], arr[right] = arr[right], arr[left]
	//		left++
	//		right--
	//	}
	//}
	//
	//// 递归排序左右两边的子数组
	//var sortedLeft []int
	//if left > 1 {
	//	sortedLeft = s.QuickSort(arr[1:left])
	//}
	//var sortedRight []int
	//if right+1 < len(arr) {
	//	sortedRight = s.QuickSort(arr[left:])
	//}
	//
	//// 合并排序后的子数组和中间的pivot
	//result := append(sortedLeft, pivot)
	//result = append(result, sortedRight...)
	result := arr
	result[0] = 10000
	return result
}
