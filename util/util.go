package util

// helper function to convert a string slice to a map.
func SliceToMap(s []string) map[string]bool {
	v := map[string]bool{}
	for _, ss := range s {
		v[ss] = true
	}
	return v
}
