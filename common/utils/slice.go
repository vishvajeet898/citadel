package utils

import "sort"

// Find the elements from slice1 which are not present in slice2
func GetDifferenceBetweenStringSlices(slice1 []string, slice2 []string) []string {
	seen := map[string]bool{}
	for _, element := range slice2 {
		seen[element] = true
	}

	var difference []string
	for _, element := range slice1 {
		if !seen[element] {
			difference = append(difference, element)
		}
	}

	return difference
}

func GetDifferenceBetweenUintSlices(slice1 []uint, slice2 []uint) []uint {
	seen := map[uint]bool{}
	for _, element := range slice2 {
		seen[element] = true
	}

	var difference []uint
	for _, element := range slice1 {
		if !seen[element] {
			difference = append(difference, element)
		}
	}

	return difference
}

func GetCommonElementsBetweenUintSlices(slice1 []uint, slice2 []uint) []uint {
	seen := map[uint]bool{}
	for _, element := range slice2 {
		seen[element] = true
	}

	var common []uint
	for _, element := range slice1 {
		if seen[element] {
			common = append(common, element)
		}
	}

	return common
}

func GetCommonElementsBetweenStringSlices(slice1 []string, slice2 []string) []string {
	seen := map[string]bool{}
	for _, element := range slice2 {
		seen[element] = true
	}

	var common []string
	for _, element := range slice1 {
		if seen[element] {
			common = append(common, element)
		}
	}

	return common
}

func AreEqualUintSlices(slice1, slice2 []uint) bool {
	if len(slice1) != len(slice2) {
		return false
	}

	// Use a map to count occurrences of each element in the first slice
	counts := make(map[uint]int)
	for _, elem := range slice1 {
		counts[elem]++
	}

	// Check if each element in the second slice has a matching count
	for _, elem := range slice2 {
		if counts[elem] == 0 {
			return false
		}
		counts[elem]--
	}

	return true
}

func SliceContainsInt(slice []int, element int) bool {
	for _, sliceElement := range slice {
		if sliceElement == element {
			return true
		}
	}
	return false
}

func SliceContainsUint(slice []uint, element uint) bool {
	for _, sliceElement := range slice {
		if sliceElement == element {
			return true
		}
	}
	return false
}

func SliceContainsString(slice []string, element string) bool {
	for _, sliceElement := range slice {
		if sliceElement == element {
			return true
		}
	}
	return false
}

func RemoveEmptyStringsFromStringSlice(slice []string) []string {
	var nonEmptySlice []string
	for _, element := range slice {
		if element != "" {
			nonEmptySlice = append(nonEmptySlice, element)
		}
	}
	return nonEmptySlice
}

func CreateUniqueSliceString(slice []string) []string {
	seen := map[string]bool{}
	var uniqueSlice []string
	for _, element := range slice {
		if !seen[element] {
			uniqueSlice = append(uniqueSlice, element)
			seen[element] = true
		}
	}
	return uniqueSlice
}

func CreateUniqueSliceUint(slice []uint) []uint {
	seen := map[uint]bool{}
	var uniqueSlice []uint
	for _, element := range slice {
		if !seen[element] {
			uniqueSlice = append(uniqueSlice, element)
			seen[element] = true
		}
	}
	return uniqueSlice
}

func StringSliceHasDuplicates(slice []string) bool {
	seen := map[string]bool{}
	for _, element := range slice {
		if seen[element] {
			return true
		}
		seen[element] = true
	}
	return false
}

func SortUintSlice(slice []uint) []uint {
	sort.Slice(slice, func(i, j int) bool {
		return slice[i] < slice[j]
	})
	return slice
}
