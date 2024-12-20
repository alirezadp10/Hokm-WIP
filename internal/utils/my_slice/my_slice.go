package my_slice

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

func Remove[T comparable](element T, slice []T) []T {
    for i, v := range slice {
        if v == element {
            slice = append(slice[:i], slice[i+1:]...)
        }
    }
    return slice
}
