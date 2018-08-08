package smartfeetable

import (
	"math"
	"sort"
)

type sortedFees []FeeWrapper

func (s sortedFees) Len() int {
	return len(s)
}

func (s sortedFees) Less(i, j int) bool {
	return s[i].LowerBound < s[j].LowerBound
}

func (s sortedFees) Swap(i, j int) {
	data := s[i]
	s[i] = s[j]
	s[j] = data
}

func (s sortedFees) Add(entry FeeWrapper) sortedFees {
	result := append(s, entry)
	sort.Sort(result)
	return result
}

func FillFeeGaps(rawFees []FeeWrapper, zeroFee FeeWrapper) []FeeWrapper {
	if len(rawFees) == 0 {
		return nil
	}

	fees := sortedFees(rawFees)
	sort.Sort(fees)

	// check lower bound
	if fees[0].LowerBound != 0 && fees[0].LowerBound != 1 {
		fees = fees.Add(getNewZeroFee(zeroFee, 0, fees[0].LowerBound-1))
	}

	// check upper bound
	if fees[fees.Len()-1].UpperBound != math.MaxInt64 {
		fees = fees.Add(getNewZeroFee(zeroFee, fees[fees.Len()-1].UpperBound+1, math.MaxInt64))
	}

	for i := 0; i < fees.Len()-1; i++ {
		if fees[i].UpperBound == math.MaxInt64 {
			break
		}

		expectedLowerBoundForNextFee := fees[i].UpperBound + 1
		if expectedLowerBoundForNextFee != fees[i+1].LowerBound {
			fees = fees.Add(getNewZeroFee(zeroFee, expectedLowerBoundForNextFee, fees[i+1].LowerBound-1))
		}
	}

	return fees
}

func getNewZeroFee(zeroFee FeeWrapper, lowerBound, upperBound int64) FeeWrapper {
	zeroFee.LowerBound = lowerBound
	zeroFee.UpperBound = upperBound
	zeroFee.NotExist = true
	return zeroFee
}

func SmartFillFeeGaps(rawFees []FeeWrapper, otherFees []FeeWrapper) []FeeWrapper {
	if len(rawFees) == 0 || len(otherFees) == 0 {
		for _, v := range otherFees {
			rawFees = append(rawFees, v)
		}
		return rawFees
	}

	fees := sortedFees(rawFees)
	sort.Sort(fees)
	other := sortedFees(otherFees)
	sort.Sort(other)

	// check lower bound
	for fees[0].LowerBound > other[0].LowerBound {
		newFee, ok := fillGap(otherFees, 0, fees[0].LowerBound-1)
		if ok {
			fees = fees.Add(newFee)
		}
	}

	// check upper bound
	for fees[fees.Len()-1].UpperBound < other[other.Len()-1].UpperBound {
		newFee, ok := fillGap(otherFees, fees[fees.Len()-1].UpperBound+1, math.MaxInt64)
		if ok {
			fees = fees.Add(newFee)
		}
	}

	for i := 0; i < fees.Len()-1; i++ {
		if fees[i].UpperBound == math.MaxInt64 {
			break
		}

		expectedLowerBoundForNextFee := fees[i].UpperBound + 1
		if expectedLowerBoundForNextFee != fees[i+1].LowerBound {
			newFee, ok := fillGap(otherFees, expectedLowerBoundForNextFee, fees[i+1].LowerBound-1)
			if ok {
				fees = fees.Add(newFee)
			}
		}
	}

	return fees
}

func fillGap(fees []FeeWrapper, lowerBound, upperBound int64) (fee FeeWrapper, ok bool) {
	for _, v := range fees {
		if ok, l, b := overlap(lowerBound, upperBound, v.LowerBound, v.UpperBound); ok {
			fee = v
			fee.LowerBound = l
			fee.UpperBound = b
			return fee, true
		}
	}
	return fee, false
}

func overlap(a, b, c, d int64) (ok bool, from int64, to int64) {
	if a > b || c > d {
		return
	}
	if b-c >= 0 && d-a >= 0 {
		return true, max(a, c), min(b, d)
	}

	return false, 0, 0
}

func min(a, b int64) int64 {
	if a > b {
		return b
	}
	return a
}

func max(a, b int64) int64 {
	if a < b {
		return b
	}
	return a
}
