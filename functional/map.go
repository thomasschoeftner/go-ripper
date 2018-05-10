package functional

//type MapFunc func(interface{}) (interface{}, error)
//func (si *iterator) Map(f MapFunc) (*iterator, error) {
//	result := GenericSlice{}
//	for si.HasNext() {
//		val, _ := si.Next()
//		mapped, error := f(val)
//
//		if error != nil {
//			return nil, error
//		} else {
//			result = append(result, mapped)
//		}
//	}
//	return NewIterator(result), nil
//}
//
//func Map(seq ISequence, f MapFunc) (ISequence, error)  {
//	result, error := NewIterator(seq).Map(f)
//	if error != nil {
//		return nil, error
//	} else {
//		return result.ToSlice(), nil
//	}
//}
//
