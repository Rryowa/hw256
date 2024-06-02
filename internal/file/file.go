package file

import (
	"encoding/json"
	"errors"
	"fmt"
	"homework-1/internal/entities"
	"os"
)

// FileService defines file operations
type FileService interface {
	CreateFile() error
	CheckFile() error
	Write(map[string]entities.Order) error
	Read() (map[string]entities.Order, []string, error)
}

type file struct {
	fileName string
}

func NewFileService(fileName string) FileService {
	return &file{fileName: fileName}
}

func (fs *file) CheckFile() error {
	if _, err := os.Stat(fs.fileName); errors.Is(err, os.ErrNotExist) {
		return os.ErrNotExist
	} else {
		return nil
	}
}

func (fs *file) CreateFile() error {
	file, err := os.Create(fs.fileName)
	if err != nil {
		return fmt.Errorf("failed to create file: %w", err)
	}
	defer file.Close()

	return nil
}

func (fs *file) Write(orders map[string]entities.Order) error {
	bWrite, err := json.MarshalIndent(orders, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(fs.fileName, bWrite, 0666)
}

func (fs *file) Read() (map[string]entities.Order, []string, error) {
	file, err := os.Open(fs.fileName)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to open file: %w", err)
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
