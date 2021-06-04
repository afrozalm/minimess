package sets

type Set map[interface{}]struct{}

func (s Set) Insert(x interface{}) {
	s[x] = struct{}{}
}

func (s Set) Remove(x interface{}) {
	delete(s, x)
}

func (s Set) Exists(x interface{}) bool {
	if _, ok := s[x]; ok {
		return true
	}
	return false
}

func (s Set) Size() int {
	return len(s)
}

func (s Set) Empty() bool {
	return len(s) == 0
}
