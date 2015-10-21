package dsl

type Filter func(float64) bool

func NumFilter(comparator string, threshold float64) Filter {
	return func(v float64) bool {
		switch comparator {
		case ">":
			return v > threshold
		case "<":
			return v < threshold
		case ">=":
			return v >= threshold
		case "<=":
			return v <= threshold
		case "==":
			return v == threshold
		}
		return false
	}
}
