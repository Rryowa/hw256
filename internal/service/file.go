package service

import (
	"encoding/json"
	"errors"
	"fmt"
	"homework-1/internal/entities"
	"os"
	"sync"
)

// FileService defines storage operations
type FileService interface {
	CheckFile() error
	Write(map[string]entities.Order) error
	Read() (map[string]entities.Order, []string, error)
	IsEmpty() (bool, error)
}

type file struct {
	fileName string
	mu       sync.Mutex
}

func NewFileService(fileName string) FileService {
	return &file{fileName: fileName}
}

func (fs *file) CheckFile() error {
	if _, err := os.Stat(fs.fileName); errors.Is(err, os.ErrNotExist) {
		file, err := os.Create(fs.fileName)
		if err != nil {
			return err
		}
		defer file.Close()
	}
	return nil
}

func (fs *file) Write(orders map[string]entities.Order) error {
	fs.mu.Lock()
	defer fs.mu.Unlock()
	bWrite, err := json.MarshalIndent(orders, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(fs.fileName, bWrite, 0666)
}

func (fs *file) Read() (map[string]entities.Order, []string, error) {
	fs.mu.Lock()
	defer fs.mu.Unlock()
	file, err := os.Open(fs.fileName)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to open storage: %w", err)
	}
	defer file.Close()

	var orders map[string]entities.Order
	decoder := json.NewDecoder(file)
	if err := decoder.Decode(&orders); err != nil {
		return nil, nil, fmt.Errorf("failed to decode JSON: %w", err)
	}

	ids := make([]string, 0, len(orders))
	for id := range orders {
		ids = append(ids, id)
	}

	return orders, ids, nil
}

func (fs *file) IsEmpty() (bool, error) {
	if err := fs.CheckFile(); err != nil {
		return false, err
	}
	file, err := os.Stat(fs.fileName)
	if err != nil {
		return false, err
	}
	if file.Size() != 0 {
		return false, nil
	}
	return true, nil
}
