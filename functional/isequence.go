package functional
//
//type ISequence interface {
//	Len() int
//	Get(int) interface{}
//	AsSlice() []interface{}
//}
//
//
//type GenericSlice []interface{}
//
//func ToISequence(slice interface{}) GenericSlice {
//	return slice.(GenericSlice)
//}
//
//func (gs GenericSlice) Len() int {
//	return len(gs)
//}
//
//func (gs GenericSlice) Get(idx int) interface{} {
//	return gs[idx]
//}
//
//func (gs GenericSlice) AsSlice() []interface{} {
//	return gs
//}