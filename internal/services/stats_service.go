package service

import (
	"fmt"
	"os"
)

type Stats struct {
	ProcessedOld int
	ProcessedNew int
	Recorded     int
	Thinking     int
	Rejected     int
	NoAnswer     int
}

type StatsService struct {
	stats   Stats
	history []string
}

func NewStatsService(_ interface{}) *StatsService {
	return &StatsService{}
}

func (s *StatsService) GetStatsText() string {
	total := s.stats.ProcessedOld + s.stats.ProcessedNew
	return fmt.Sprintf(
		"Обработано %d (старых %d, новых %d)\nЗаписано %d\nДумают %d\nНе подходят %d\nНе дозвонился %d",
		total, s.stats.ProcessedOld, s.stats.ProcessedNew,
		s.stats.Recorded, s.stats.Thinking, s.stats.Rejected, s.stats.NoAnswer,
	)
}

func (s *StatsService) AddOld()      { s.stats.ProcessedOld++; s.history = append(s.history, "old") }
func (s *StatsService) AddNew()      { s.stats.ProcessedNew++; s.history = append(s.history, "new") }
func (s *StatsService) AddRecorded() { s.stats.Recorded++; s.history = append(s.history, "recorded") }
func (s *StatsService) AddThinking() { s.stats.Thinking++; s.history = append(s.history, "thinking") }
func (s *StatsService) AddRejected() { s.stats.Rejected++; s.history = append(s.history, "rejected") }
func (s *StatsService) AddNoAnswer() { s.stats.NoAnswer++; s.history = append(s.history, "noanswer") }

func (s *StatsService) UndoLast() {
	if len(s.history) == 0 {
		return
	}
	last := s.history[len(s.history)-1]
	s.history = s.history[:len(s.history)-1]

	switch last {
	case "old":
		s.stats.ProcessedOld--
	case "new":
		s.stats.ProcessedNew--
	case "recorded":
		s.stats.Recorded--
	case "thinking":
		s.stats.Thinking--
	case "rejected":
		s.stats.Rejected--
	case "noanswer":
		s.stats.NoAnswer--
	}
}

func (s *StatsService) Reset() {
	s.stats = Stats{}
	s.history = nil
	os.Remove("state.json")
}

func (s *StatsService) SaveReport() error {
	data := s.GetStatsText()
	return os.WriteFile("report.txt", []byte(data), 0644)
}
