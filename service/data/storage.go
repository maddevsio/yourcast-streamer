package data

import (
	"container/list"
	"sync"
)

// StreamItem storage for stream
type StreamItem struct {
	sync.RWMutex
	ID     int
	Name   string
	Slug   string
	IsAuto bool
	Links  *list.List
}

// StreamStorage storage for multiple StreamItems
type StreamStorage struct {
	sync.RWMutex
	Items map[int]StreamItem
}

// NewStreamStorage initializes a new storage
func NewStreamStorage() *StreamStorage {
	ss := &StreamStorage{}
	i := make(map[int]StreamItem)
	ss.Items = i
	return ss
}
