package functional

import (
	"errors"
)

type iterator struct {
	seq ISequence
	currentPosition int
	step int
}

func NewIterator(sequence ISequence) *iterator {
	return &iterator{
		seq: sequence,
		currentPosition: -1,
		step: 1,
	}
}

func (si *iterator) nextPosition() int {
	return si.currentPosition + si.step
}

func (si *iterator) IsOutOfBounds(position int) bool {
	return position < 0 || position >= si.seq.Len()
}

func (si *iterator) Reverse() *iterator {
	return &iterator{
		seq: si.seq,
		currentPosition: si.seq.Len(),
		step: -si.step,
	}
}

func (si *iterator) Next() (interface{}, error) {
	nextPos := si.nextPosition()
	if si.IsOutOfBounds(nextPos) {
		return nil, errors.New("index out of bounds on slice iterator")
	}

	si.currentPosition = nextPos
	return si.seq.Get(si.currentPosition), nil
}

func (si* iterator) HasNext() bool {
	return !si.IsOutOfBounds(si.nextPosition())
}

func (si* iterator) ToSlice() GenericSlice {
	results := GenericSlice{}
	for si.HasNext() {
		val, _ := si.Next()
		results = append(results, val)
	}
	return results
}
