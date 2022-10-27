package clock

import "time"

// TimeNowFunc функция для определения текущего времени.
var TimeNowFunc = time.Now

// Time структура для работы с датой и временем.
type Time = time.Time

// NewFromTime создает Time на основе time.Time.
func NewFromTime(t time.Time) Time {
	return t
}

// Now возвращает текущее время.
func Now() Time {
	return TimeNowFunc()
}
