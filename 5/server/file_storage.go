package server

import (
	"context"
	"encoding/json"
	"os"
	"time"
)

type FileStorage struct {
	FilePath string
}

func NewFileStorage(filePath string) *FileStorage {
	return &FileStorage{
		FilePath: filePath,
	}
}

type Task struct {
	TaskID   int32  `json:"taskId"`
	TaskName string `json:"taskName"`
}

type fileStorageError struct {
	message string
}

func (e *fileStorageError) Error() string {
	return e.message
}

func (fs *FileStorage) Save(ctx context.Context, task Task) error {
	select {
	case <-ctx.Done():
		return context.Canceled
	default:
		time.Sleep(2 * time.Second)

		data, err := json.Marshal(task)
		if err != nil {
			return &fileStorageError{message: "Error encoding task"}
		}

		err = os.WriteFile(fs.FilePath, data, 0644)
		if err != nil {
			return &fileStorageError{message: "Error writing to file"}
		}

		return nil
	}
}

func (fs *FileStorage) Get(ctx context.Context) (Task, error) {
	select {
	case <-ctx.Done():
		return Task{}, context.Canceled
	default:
		data, err := os.ReadFile(fs.FilePath)
		if err != nil {
			if os.IsNotExist(err) {
				return Task{}, &fileStorageError{message: "File not found"}
			}
			return Task{}, &fileStorageError{message: "Error reading from file"}
		}

		var task Task
		err = json.Unmarshal(data, &task)
		if err != nil {
			return Task{}, &fileStorageError{message: "Error decoding task"}
		}

		return task, nil
	}
}
