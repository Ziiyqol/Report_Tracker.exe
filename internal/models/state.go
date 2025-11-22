package models

// State хранит текущие показатели счетчиков
type State struct {
	ProcessedOld int `json:"processed_old"`
	ProcessedNew int `json:"processed_new"`
	Recorded     int `json:"recorded"`
	Reserved     int `json:"reserved"`
	Thinking     int `json:"thinking"`
	Rejected     int `json:"rejected"`
	NoAnswer     int `json:"no_answer"`
}

// HistoryItem нужна для функции "Отмена действия"
type HistoryItem string

const (
	ActionOld      HistoryItem = "old"
	ActionNew      HistoryItem = "new"
	ActionRecorded HistoryItem = "recorded"
	ActionReserved HistoryItem = "reserved"
	ActionThinking HistoryItem = "thinking"
	ActionRejected HistoryItem = "rejected"
	ActionNoAnswer HistoryItem = "noanswer"
)
