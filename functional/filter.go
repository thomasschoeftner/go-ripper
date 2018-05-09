package functional


type FilterFunc func(interface{}) bool
func (si *iterator) Filter(f FilterFunc) *iterator {
	result := GenericSlice{}
	for si.HasNext() {
		val, _ := si.Next()
		if f(val) {
			result = append(result, val)
		}
	}
	return NewIterator(result)
}

func Filter(seq ISequence, f FilterFunc) ISequence  {
	return NewIterator(seq).Filter(f).ToSlice()
}
