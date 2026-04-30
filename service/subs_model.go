package service

import (
	"errors"
	"time"

	"github.com/google/uuid"
)

type SubscriptionRequest struct {
	ServiceName string    `json:"service_name" binding:"required"`
	Price       int       `json:"price" binding:"gte=0"`
	UserId      uuid.UUID `json:"user_id" binding:"required"`
	StartDate   string    `json:"start_date" binding:"required"`
	EndDate     string    `json:"end_date"`
}

type SubscriptionResponse struct {
	Id          uuid.UUID `json:"id"`
	ServiceName string    `json:"service_name"`
	Price       int       `json:"price"`
	UserId      uuid.UUID `json:"user_id"`
	StartDate   string    `json:"start_date"`
	EndDate     *string   `json:"end_date,omitempty"`
}

type SubscriptionUpdateRequest struct {
	ServiceName string `json:"service_name" binding:"required"`
	Price       int    `json:"price" binding:"gte=0"`
	StartDate   string `json:"start_date" binding:"required"`
	EndDate     string `json:"end_date"`
}

type SubscriptionUpdate struct {
	Id          uuid.UUID  `db:"id"`
	ServiceName string     `json:"service_name"`
	Price       int        `json:"price"`
	StartDate   time.Time  `json:"start_date"`
	EndDate     *time.Time `json:"end_date"`
}

type Subscription struct {
	Id          uuid.UUID  `db:"id"`
	ServiceName string     `db:"service_name"`
	Price       int        `db:"price"`
	UserId      uuid.UUID  `db:"user_id"`
	StartDate   time.Time  `db:"start_date"`
	EndDate     *time.Time `db:"end_date"`
}

func parseDates(startStr, endStr string) (time.Time, *time.Time, error) {
	start, err := time.Parse("01-2006", startStr)
	if err != nil {
		return time.Time{}, nil, errors.New("incorrect date format")
	}

	var end *time.Time

	if endStr != "" {
		t, err := time.Parse("01-2006", endStr)
		if err != nil {
			return time.Time{}, nil, errors.New("incorrect date format")
		}

		if t.Before(start) {
			return time.Time{}, nil, errors.New("end_date cannot be before start_date")
		}

		end = &t
	}

	return start, end, nil
}

func NewSubscription(s SubscriptionRequest) (Subscription, error) {
	start, end, err := parseDates(s.StartDate, s.EndDate)
	if err != nil {
		return Subscription{}, err
	}

	return Subscription{
		Id:          uuid.New(),
		ServiceName: s.ServiceName,
		Price:       s.Price,
		UserId:      s.UserId,
		StartDate:   start,
		EndDate:     end,
	}, nil
}

func UpdateSubscription(s SubscriptionUpdateRequest, subId uuid.UUID) (SubscriptionUpdate, error) {
	start, end, err := parseDates(s.StartDate, s.EndDate)
	if err != nil {
		return SubscriptionUpdate{}, err
	}
	return SubscriptionUpdate{
		Id:          subId,
		ServiceName: s.ServiceName,
		Price:       s.Price,
		StartDate:   start,
		EndDate:     end,
	}, nil
}

func ToResponse(subs Subscription) SubscriptionResponse {
	var end *string
	start := subs.StartDate.Format("01-2006")
	if subs.EndDate != nil {
		e := subs.EndDate.Format("01-2006")
		end = &e
	}

	return SubscriptionResponse{
		Id:          subs.Id,
		ServiceName: subs.ServiceName,
		Price:       subs.Price,
		UserId:      subs.UserId,
		StartDate:   start,
		EndDate:     end,
	}
}

func ToResponseList(sub []Subscription) []SubscriptionResponse {
	response := make([]SubscriptionResponse, 0, len(sub))
	for _, v := range sub {
		response = append(response, ToResponse(v))
	}
	return response
}

func ParsePeriod(fromStr, toStr string) (time.Time, time.Time, error) {
	from, err := time.Parse("01-2006", fromStr)
	if err != nil {
		return time.Time{}, time.Time{}, errors.New("invalid from date")
	}

	to, err := time.Parse("01-2006", toStr)
	if err != nil {
		return time.Time{}, time.Time{}, errors.New("invalid to date")
	}

	if to.Before(from) {
		return time.Time{}, time.Time{}, errors.New("to cannot be before from")
	}

	return from, to, nil
}

func countMonths(subStart time.Time, subEnd *time.Time, from, to time.Time) int {
	start := subStart
	if start.Before(from) {
		start = from
	}

	end := to
	if subEnd != nil && subEnd.Before(to) {
		end = *subEnd
	}

	if start.After(end) {
		return 0
	}

	year1, month1, _ := start.Date()
	year2, month2, _ := end.Date()

	return (year2-year1)*12 + int(month2-month1) + 1
}

func CountTotal(sub []Subscription, from, to time.Time) int {
	var total int
	for _, v := range sub {
		months := countMonths(v.StartDate, v.EndDate, from, to)
		total += months * v.Price
	}
	return total
}
