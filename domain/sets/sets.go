package sets

type Set map[interface{}]struct{}

func (s Set) Insert(x interface{}) {
	s[x] = struct{}{}
}

func (s Set) Remove(x interface{}) {
	delete(s, x)
}
