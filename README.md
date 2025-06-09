# TODO

## SM2Scheduler

### Burnout Protection
1. Modify the SRS tool to limit reviews per day (e.g., max 12 problems). If a review is due but the day is full, shift it to the next day.
example.
```go
func (s *SM2Scheduler) ScheduleWithLimit(q *Question, grade Familiarity, dailyLimit int, currentLoad int) bool {
    if currentLoad >= dailyLimit {
        return false // Defer to next day
    }
    s.Schedule(q, grade)
    return true
}
```
2. A Priority Queue for showing due questions
example.
```go
const (
    importanceWeight  = 0.5
    familiarityWeight = 0.3
    overdueWeight     = 0.2
)
priority := importanceWeight*float64(q.Importance) +
			familiarityWeight*float64(5-q.Familiarity) +
			overdueWeight*overdueDays
```