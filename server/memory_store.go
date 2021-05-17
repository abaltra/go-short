package main

import (
	"log"
)

type Store struct {
	current_idx int
	data        []string
}

// var data_store Store = Store{
// 	current_idx: -1,
// 	data:        []entry{},
// }

func NewStore() *Store {
	s := new(Store)
	s.current_idx = 0
	s.data = []string{}

	return s
}

func (s *Store) InsertURL(url string) int {
	idx := s.current_idx

	log.Printf("Inserting new entry in index %d", idx)
	s.data = append(s.data, url)
	s.current_idx += 1
	log.Print(s)
	return idx
}

func (s *Store) GetURL(idx int) string {
	log.Printf("Trying to retrieve data from index %d", idx)
	if idx < 0 || idx > s.current_idx {
		return ""
	}

	return s.data[idx]
}
