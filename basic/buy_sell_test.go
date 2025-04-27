package main_test

import "testing"

func BuySellStock(prices []int) int {
	l, r := 0, 1
	max := 0
	for r < len(prices) {
		if prices[l] < prices[r] {
			profit := prices[r] - prices[l]
			if profit > max {
				max = profit
			}
		} else {
			l++
		}
		r++
	}
	return max
}

func TestBuySell(t *testing.T) {

	tests := []struct {
		prices   []int
		expected int
	}{
		{[]int{7, 1, 5, 3, 6, 4}, 5},
		{[]int{7, 6, 4, 3, 1}, 0},
		{[]int{1, 2, 3, 4, 5}, 4},
		{[]int{7, 6, 4, 3, 1}, 0},
	}
	for _, test := range tests {
		t.Run("", func(t *testing.T) {
			actual := BuySellStock(test.prices)
			if actual != test.expected {
				t.Errorf("BuySellStock(%v) = %v, want %v", test.prices, actual, test.expected)
			}
		})
	}
}
func BuySellStock2(prices []int) int {
	if len(prices) == 0 {
		return 0
	}

	minPrice := prices[0]
	maxProfit := 0

	for _, price := range prices[1:] {
		if price < minPrice {
			minPrice = price
		} else if profit := price - minPrice; profit > maxProfit {
			maxProfit = profit
		}
	}

	return maxProfit
}

func BenchmarkBuySellStock(b *testing.B) {
	prices := []int{7, 1, 5, 3, 6, 4}
	for i := 0; i < b.N; i++ {
		BuySellStock(prices)
	}
}

func BenchmarkBuySellStock2(b *testing.B) {
	prices := []int{7, 1, 5, 3, 6, 4}
	for i := 0; i < b.N; i++ {
		BuySellStock2(prices)
	}
}
