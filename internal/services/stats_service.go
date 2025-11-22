package services

import (
	"fmt"
	"os"
	"report/internal/models"
	"report/internal/storage"
)

type StatsService struct {
	state   models.State
	history []models.HistoryItem
	storage storage.Storage
}

func NewStatsService(store storage.Storage) *StatsService {
	// При запуске загружаем сохраненные данные
	existingState := store.Load()
	return &StatsService{
		state:   existingState,
		storage: store,
		history: make([]models.HistoryItem, 0),
	}
}

// GetStatsText формирует текст для UI
func (s *StatsService) GetStatsText() string {
	total := s.state.ProcessedOld + s.state.ProcessedNew
	return fmt.Sprintf(
		"Обработано %d (старых %d, новых %d)\nЗаписано %d (в резерв %d)\nДумают %d\nНе подходят %d\nНе дозвонился %d",
		total, s.state.ProcessedOld, s.state.ProcessedNew,
		s.state.Recorded, s.state.Reserved,
		s.state.Thinking, s.state.Rejected, s.state.NoAnswer,
	)
}

// saveState - вспомогательная функция для сохранения после каждого чиха
func (s *StatsService) saveState() {
	_ = s.storage.Save(s.state)
}

func (s *StatsService) AddOld() {
	s.state.ProcessedOld++
	s.history = append(s.history, models.ActionOld)
	s.saveState()
}

func (s *StatsService) AddNew() {
	s.state.ProcessedNew++
	s.history = append(s.history, models.ActionNew)
	s.saveState()
}

func (s *StatsService) AddRecorded() {
	s.state.Recorded++
	s.history = append(s.history, models.ActionRecorded)
	s.saveState()
}

func (s *StatsService) AddReserved() {
	s.state.Reserved++
	s.state.Recorded++
	s.history = append(s.history, models.ActionReserved)
	s.saveState()
}

func (s *StatsService) AddThinking() {
	s.state.Thinking++
	s.history = append(s.history, models.ActionThinking)
	s.saveState()
}

func (s *StatsService) AddRejected() {
	s.state.Rejected++
	s.history = append(s.history, models.ActionRejected)
	s.saveState()
}

func (s *StatsService) AddNoAnswer() {
	s.state.NoAnswer++
	s.history = append(s.history, models.ActionNoAnswer)
	s.saveState()
}

func (s *StatsService) UndoLast() {
	if len(s.history) == 0 {
		return
	}
	last := s.history[len(s.history)-1]
	s.history = s.history[:len(s.history)-1]

	switch last {
	case models.ActionOld:
		s.state.ProcessedOld--
	case models.ActionNew:
		s.state.ProcessedNew--
	case models.ActionRecorded:
		s.state.Recorded--
	case models.ActionReserved:
		s.state.Reserved--
	case models.ActionThinking:
		s.state.Thinking--
	case models.ActionRejected:
		s.state.Rejected--
	case models.ActionNoAnswer:
		s.state.NoAnswer--
	}
	s.saveState()
}

func (s *StatsService) Reset() {
	s.state = models.State{}
	s.history = nil
	_ = s.storage.Reset()
}

func (s *StatsService) SaveReportToFile() error {
	data := s.GetStatsText()
	return os.WriteFile("report.txt", []byte(data), 0644)
}
