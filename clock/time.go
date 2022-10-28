package clock

import (
	"sync"
	"time"
)

var mu sync.Mutex

// timeNowFunc функция для определения текущего времени.
var timeNowFunc = time.Now

// ResetTimeNowFunc сбрасывает функцию для определния текущего времени.
func ResetTimeNowFunc() {
	mu.Lock()
	defer mu.Unlock()

	timeNowFunc = time.Now
}

// SetTimeNowFunc устанавливает функцию для определения текущего времени.
func SetTimeNowFunc(f func() time.Time) {
	mu.Lock()
	defer mu.Unlock()

	timeNowFunc = f
}

// Time структура для работы с датой и временем.
type Time = time.Time

// NewFromTime создает Time на основе time.Time.
func NewFromTime(t time.Time) Time {
	return t
}

// Now возвращает текущее время.
func Now() Time {
	mu.Lock()
	defer mu.Unlock()

	return timeNowFunc()
}
