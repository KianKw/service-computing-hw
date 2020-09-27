package main

import "fmt"

func BinarySearch(nums[]int, target int) int {
	left, right := 0, len(nums)
	for left < right {
		mid := left + (right - left) / 2
		if nums[mid] == target {
			right = mid
		} else if nums[mid] < target {
			left = mid + 1
		} else {
			right = mid
		}
	}
	if nums[left] == target {
		return left
	} else {
		return -1
	}
}

func main() {
	nums := []int{1,2,3,4,5,6,7,8,9}
	fmt.Println(BinarySearch(nums, 7))
}
