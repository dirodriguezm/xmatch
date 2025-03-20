package utils

type Set map[string]struct{}

func (s Set) Add(key string) {
	s[key] = struct{}{}
}

func (s Set) Contains(key string) bool {
	_, ok := s[key]
	return ok
}

func (s Set) ToSlice() []string {
	keys := make([]string, 0, len(s))
	for k := range s {
		keys = append(keys, k)
	}
	return keys
}

func (s Set) Remove(key string) {
	delete(s, key)
}
