package chiserver

import "sync"

type Task struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

type User struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

type DataStore struct {
	Tasks sync.Map
	Users sync.Map
}
