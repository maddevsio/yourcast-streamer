package data

import "container/list"

// StreamLink used for parsing json urls
type StreamLink struct {
	URL string `json:"url"`
}

// Stream struct for parsing WebUIAPi responses
type Stream struct {
	Name            string       `json:"name"`
	ID              int          `json:"id"`
	Slug            string       `json:"slug"`
	Links           []StreamLink `json:"links"`
	Keywords        string       `json:"keywords"`
	Channels        string       `json:"channels"`
	UpdateFrequency int          `json:"update_frequency"`
	VideoLength     int          `json:"video_length"`
	IsNews          bool         `json:"is_news"`
	IsAuto          bool
}

// IsAutoStream returns true if stream is for botService
func (s *Stream) IsAutoStream() bool {
	return s.Keywords != "" || s.Channels != ""
}

// ToStreamItem converts Stream struct into StreamItem
func (s *Stream) ToStreamItem() StreamItem {
	si := StreamItem{
		ID:     s.ID,
		Name:   s.Name,
		Slug:   s.Slug,
		IsAuto: s.IsAuto,
	}
	l := list.New()
	for _, link := range s.Links {
		l.PushBack(link.URL)
	}
	si.Links = l
	return si
}
