package myslice

import (
    "fmt"
    "strconv"
)

// Has Function to check if a slice contains an element
func Has[T comparable](slice []T, element T) bool {
    for _, v := range slice {
        if v == element {
            return true
        }
    }
    return false
}

// HasLike Function to check if a slice contains an element satisfying a predicate
func HasLike[T any](slice []T, predicate func(T) bool) int {
    for i, v := range slice {
        if predicate(v) {
            return i
        }
    }
    return -1
}

// GetIndex Find the index of an element
func GetIndex[T comparable](element T, slice []T) int {
    index := -1
    for i, player := range slice {
        if player == element {
            index = i
            break
        }
    }
    return index
}

func StringsToInt64s(strings []string) ([]int64, error) {
    ints := make([]int64, len(strings))
    for i, s := range strings {
        num, err := strconv.ParseInt(s, 10, 64)
        if err != nil {
            return nil, fmt.Errorf("error converting %s to int64: %w", s, err)
        }
        ints[i] = num
    }
    return ints, nil
}
