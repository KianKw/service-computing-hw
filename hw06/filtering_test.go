package rxgo

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestDebounce(t *testing.T) {
	res := []int{}
	ob := Just(100, 200, 300, 400).Map(func(x int) int {
		time.Sleep(20 * time.Millisecond)
		return 2 * x
	}).Debounce(30 * time.Millisecond)
	ob.Subscribe(func(x int) {
		res = append(res, x)
	})
	assert.Equal(t, []int{}, res, "Debounce Test Error!")
}

func TestDistinct(t *testing.T) {
	res := []int{}
	ob := Just(1, 2, 3, 4, 1, 2, 3, 4).Map(func(x int) int {
		return 2 * x
	}).Distinct()
	ob.Subscribe(func(x int) {
		res = append(res, x)
	})
	assert.Equal(t, []int{2, 4, 6, 8}, res, "Distinct Test Error!")
}

func TestElementAt(t *testing.T) {
	res := []int{}
	ob := Just(1, 2, 3, 4, 5, 6, 7).Map(func(x int) int {
		return 2 * x
	}).ElementAt(4)
	ob.Subscribe(func(x int) {
		res = append(res, x)
	})
	assert.Equal(t, []int{8}, res, "ElementAt Test Error!")
}

func TestFirst(t *testing.T) {
	res := []int{}
	ob := Just(1, 2, 3).Map(func(x int) int {
		return 2 * x
	}).First()
	ob.Subscribe(func(x int) {
		res = append(res, x)
	})
	assert.Equal(t, []int{2}, res, "First Test Error!")
}

func TestLast(t *testing.T) {
	res := []int{}
	ob := Just(10, 20, 30).Map(func(x int) int {
		return 2 * x
	}).Last()
	ob.Subscribe(func(x int) {
		res = append(res, x)
	})

	assert.Equal(t, []int{60}, res, "Last Test Error!")
}

func TestSample(t *testing.T) {
	res := []int{}
	Just(1, 2, 3, 4, 3, 1, 2, 4, 3).Map(func(x int) int {
		time.Sleep(20 * time.Millisecond)
		return 2 * x
	}).Sample(15 * time.Millisecond).Subscribe(func(x int) {
		res = append(res, x)
	})
	assert.Equal(t, []int{2, 4, 6, 8, 6, 2, 4, 8, 6}, res, "Sample Test Error!")
}

func TestSkip(t *testing.T) {
	res := []int{}
	ob := Just(1, 2, 3, 4, 5, 6, 7).Map(func(x int) int {
		return 2 * x
	}).Skip(4)
	ob.Subscribe(func(x int) {
		res = append(res, x)
	})
	assert.Equal(t, []int{10, 12, 14}, res, "Skip Test Error!")
}

func TestSkipLast(t *testing.T) {
	res := []int{}
	ob := Just(1, 2, 3, 4, 5, 6, 7).Map(func(x int) int {
		return 2 * x
	}).SkipLast(4)
	ob.Subscribe(func(x int) {
		res = append(res, x)
	})
	assert.Equal(t, []int{2, 4, 6}, res, "SkipLast Test Error!")
}

func TestTake(t *testing.T) {
	res := []int{}
	ob := Just(1, 2, 3, 4, 5, 6, 7).Map(func(x int) int {
		return 2 * x
	}).Take(4)
	ob.Subscribe(func(x int) {
		res = append(res, x)
	})
	assert.Equal(t, []int{2, 4, 6, 8}, res, "Take Test Error!")
}

func TestTakeLast(t *testing.T) {
	res := []int{}
	ob := Just(1, 2, 3, 4, 5, 6, 7).Map(func(x int) int {
		return 2 * x
	}).TakeLast(4)
	ob.Subscribe(func(x int) {
		res = append(res, x)
	})
	assert.Equal(t, []int{8, 10, 12, 14}, res, "TakeLast Test Error!")
}
