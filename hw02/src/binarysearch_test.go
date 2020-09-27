package main

import "testing"

func TestBinarySearch(t *testing.T) {
	checkBinarySearch := func(t *testing.T, target int, numbers []int, got int) {
		t.Helper()
		want := -1
		size := len(numbers)
		for i := size - 1; i >= 0; i-- {
			if target == numbers[i] {
				want = i
			}
		}
		if got != want {
            t.Errorf("got %d want %d given, %v", got, want, numbers)
		}
	}

	t.Run("the target doesn't exist", func(t *testing.T) {
        numbers := []int{1, 2, 3, 5, 6, 7}
        target := 4

        got := BinarySearch(numbers, target)

		checkBinarySearch(t, target, numbers, got)
	})

	t.Run("the target exists but does not repeat", func(t *testing.T) {
	    numbers := []int{1, 2, 3, 4, 5, 6, 7}
	    target := 2

	    got := BinarySearch(numbers, target)

		checkBinarySearch(t, target, numbers, got)
	})

	t.Run("the target exists and repeats", func(t *testing.T) {
	    numbers := []int{1, 2, 3, 3, 3, 6, 8, 9}
	    target := 3

	    got := BinarySearch(numbers, target)

		checkBinarySearch(t, target, numbers, got)
	})
}
