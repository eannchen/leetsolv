## TODO

### SM2Scheduler

1. Burnout Protection: Modify the SRS tool to limit reviews per day (e.g., max 12 problems). If a review is due but the day is full, shift it to the next day.
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
2. A Decay for Neglected Questions: To better handle questions I forget/ignore to review, you might implement a “penalty decay”
example.
```
If now.Sub(q.NextReview) is large (you missed the review), reduce EaseFactor slightly or lower familiarity.
```
3. The algorithm is entirely personal to my current condition. (Jun, 2025) It is subject to modification for the general public or the change of my condition.