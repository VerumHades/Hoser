package utils

func MergeMaps[K comparable, V any](m1, m2 map[K]V) map[K]V {
	merged := make(map[K]V)

	for k, v := range m1 {
		merged[k] = v
	}
	for k, v := range m2 {
		merged[k] = v
	}

	return merged
}

func MergeMultipleMaps[K comparable, V any](maps []map[K]V) map[K]V {
	merged := make(map[K]V)
	for i := range maps {
		merged = MergeMaps(merged, maps[i])
	}
	return merged
}
