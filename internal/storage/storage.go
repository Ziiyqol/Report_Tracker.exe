package storage

import (
	"encoding/json"
	"os"
	"report/internal/models"
)

const stateFile = "state.json"

// Storage описывает методы для работы с данными
type Storage interface {
	Save(state models.State) error
	Load() models.State
	Reset() error
}

// FileStorage реализация хранения в JSON файле
type FileStorage struct{}

func NewFileStorage() *FileStorage {
	return &FileStorage{}
}

func (fs *FileStorage) Save(state models.State) error {
	data, err := json.MarshalIndent(state, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(stateFile, data, 0644)
}

func (fs *FileStorage) Load() models.State {
	var state models.State
	data, err := os.ReadFile(stateFile)
	if err != nil {
		return models.State{} // Возвращаем пустую структуру, если файла нет
	}
	if err := json.Unmarshal(data, &state); err != nil {
		return models.State{}
	}
	return state
}

func (fs *FileStorage) Reset() error {
	return os.Remove(stateFile)
}
