package utils

import (
	"math"
	"work/pkg/constants"
)

func SlicePage(pageNum, pageSize, len int) (sliceStart, sliceEnd int) {
	if pageNum <= 0 {
		pageNum = 1
	}
	if pageSize <= 0 {
		pageSize = constants.DefaultPageSize
	}
	if pageSize > len {
		return 0, len
	}
	pageCount := int(math.Ceil(float64(len) / float64(pageSize)))
	if pageNum > pageCount {
		return 0, 0
	}
	sliceStart = (pageNum - 1) * pageSize
	sliceEnd = sliceStart + pageSize
	if sliceEnd > len {
		sliceEnd = len
	}
	return sliceStart, sliceEnd
}
