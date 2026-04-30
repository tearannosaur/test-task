package repository

import (
	"app/service"
	"database/sql"
	"time"

	"github.com/google/uuid"
)

func (r *Repository) SubscriptionSave(sub service.Subscription) error {
	query := `INSERT INTO subscriptions (id,service_name,price,user_id,start_date,end_date)
	VALUES ($1,$2,$3,$4,$5,$6)`
	_, err := r.db.Exec(query, sub.Id, sub.ServiceName, sub.Price, sub.UserId, sub.StartDate, sub.EndDate)
	if err != nil {
		return err
	}
	return nil
}

func (r *Repository) GetSubscriptionsList() ([]service.Subscription, error) {
	var subscriptions []service.Subscription
	query := `
	SELECT id, service_name, price, user_id, start_date, end_date
	FROM subscriptions`
	err := r.db.Select(&subscriptions, query)
	if err != nil {
		return []service.Subscription{}, err
	}
	return subscriptions, nil
}

func (r *Repository) GetSubscriptionById(id uuid.UUID) (service.Subscription, error) {
	var subscription service.Subscription
	query := `SELECT id, service_name, price, user_id, start_date, end_date
	FROM subscriptions
	WHERE id=$1`
	err := r.db.Get(&subscription, query, id)
	if err != nil {
		return service.Subscription{}, err
	}
	return subscription, nil
}

func (r *Repository) DeleteSubscriptionById(id uuid.UUID) error {
	query := `DELETE FROM subscriptions
	WHERE id=$1`

	result, err := r.db.Exec(query, id)
	if err != nil {
		return err
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rows == 0 {
		return sql.ErrNoRows
	}

	return nil
}

func (r *Repository) UpdateSubscriptionById(sub service.SubscriptionUpdate) error {
	query := `UPDATE subscriptions
		SET service_name=$1,price=$2,start_date=$3,end_date=$4
		WHERE id=$5`
	res, err := r.db.Exec(query, sub.ServiceName, sub.Price, sub.StartDate, sub.EndDate, sub.Id)
	if err != nil {
		return err
	}
	rows, err := res.RowsAffected()
	if err != nil {
		return err
	}

	if rows == 0 {
		return sql.ErrNoRows
	}
	return nil
}

func (r *Repository) GetSubscriptionForPeriod(userId *uuid.UUID, serviceName string, from time.Time, to time.Time) ([]service.Subscription, error) {
	var subscriptions []service.Subscription
	query := `SELECT id, service_name, price, user_id, start_date, end_date FROM subscriptions
	WHERE start_date <= $2
	AND (end_date IS NULL OR end_date >= $1)
	AND ($3::uuid IS NULL OR user_id=$3)
	AND ($4='' OR service_name=$4)
	`
	err := r.db.Select(&subscriptions, query, from, to, userId, serviceName)
	if err != nil {
		return []service.Subscription{}, err
	}
	return subscriptions, nil
}
