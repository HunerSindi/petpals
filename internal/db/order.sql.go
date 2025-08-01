// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.29.0
// source: order.sql

package db

import (
	"context"
	"database/sql"
)

const createOrder = `-- name: CreateOrder :one
INSERT INTO orders (
    user_id, total_amount, status, delivery_address
) VALUES (
    $1, $2, $3, $4
) RETURNING id, user_id, total_amount, status, delivery_address, order_date, delivered_at
`

type CreateOrderParams struct {
	UserID          sql.NullInt64   `json:"user_id"`
	TotalAmount     sql.NullString  `json:"total_amount"`
	Status          NullOrderStatus `json:"status"`
	DeliveryAddress sql.NullString  `json:"delivery_address"`
}

func (q *Queries) CreateOrder(ctx context.Context, arg CreateOrderParams) (Order, error) {
	row := q.db.QueryRowContext(ctx, createOrder,
		arg.UserID,
		arg.TotalAmount,
		arg.Status,
		arg.DeliveryAddress,
	)
	var i Order
	err := row.Scan(
		&i.ID,
		&i.UserID,
		&i.TotalAmount,
		&i.Status,
		&i.DeliveryAddress,
		&i.OrderDate,
		&i.DeliveredAt,
	)
	return i, err
}

const getOrderByID = `-- name: GetOrderByID :one
SELECT id, user_id, total_amount, status, delivery_address, order_date, delivered_at FROM orders
WHERE id = $1 LIMIT 1
`

func (q *Queries) GetOrderByID(ctx context.Context, id int64) (Order, error) {
	row := q.db.QueryRowContext(ctx, getOrderByID, id)
	var i Order
	err := row.Scan(
		&i.ID,
		&i.UserID,
		&i.TotalAmount,
		&i.Status,
		&i.DeliveryAddress,
		&i.OrderDate,
		&i.DeliveredAt,
	)
	return i, err
}

const listOrdersByUserID = `-- name: ListOrdersByUserID :many
SELECT id, user_id, total_amount, status, delivery_address, order_date, delivered_at FROM orders
WHERE user_id = $1
ORDER BY order_date DESC
`

func (q *Queries) ListOrdersByUserID(ctx context.Context, userID sql.NullInt64) ([]Order, error) {
	rows, err := q.db.QueryContext(ctx, listOrdersByUserID, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []Order{}
	for rows.Next() {
		var i Order
		if err := rows.Scan(
			&i.ID,
			&i.UserID,
			&i.TotalAmount,
			&i.Status,
			&i.DeliveryAddress,
			&i.OrderDate,
			&i.DeliveredAt,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const updateOrderStatus = `-- name: UpdateOrderStatus :one
UPDATE orders
SET
    status = $2
WHERE
    id = $1
RETURNING id, user_id, total_amount, status, delivery_address, order_date, delivered_at
`

type UpdateOrderStatusParams struct {
	ID     int64           `json:"id"`
	Status NullOrderStatus `json:"status"`
}

func (q *Queries) UpdateOrderStatus(ctx context.Context, arg UpdateOrderStatusParams) (Order, error) {
	row := q.db.QueryRowContext(ctx, updateOrderStatus, arg.ID, arg.Status)
	var i Order
	err := row.Scan(
		&i.ID,
		&i.UserID,
		&i.TotalAmount,
		&i.Status,
		&i.DeliveryAddress,
		&i.OrderDate,
		&i.DeliveredAt,
	)
	return i, err
}
