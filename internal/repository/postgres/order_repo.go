package postgres

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strings"

	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/vuongthanh148/dodongtruongthoi_be/internal/domain"
)

type OrderRepository struct {
	pool *pgxpool.Pool
}

func NewOrderRepository(pool *pgxpool.Pool) *OrderRepository {
	return &OrderRepository{pool: pool}
}

func (r *OrderRepository) List(ctx context.Context, phone *string, status *string, limit int, offset int) ([]domain.Order, error) {
	query := "SELECT id, phone, customer_name, address, note, status, admin_note, total_amount, created_at, updated_at FROM orders"

	args := make([]interface{}, 0)
	argCount := 1

	if phone != nil || status != nil {
		query += " WHERE"
		if phone != nil {
			query += fmt.Sprintf(" phone = $%d", argCount)
			args = append(args, *phone)
			argCount++
			if status != nil {
				query += " AND"
			}
		}
		if status != nil {
			query += fmt.Sprintf(" status = $%d", argCount)
			args = append(args, *status)
			argCount++
		}
	}

	query += " ORDER BY created_at DESC"
	if limit > 0 {
		query += fmt.Sprintf(" LIMIT $%d OFFSET $%d", argCount, argCount+1)
		args = append(args, limit, offset)
	}

	rows, err := r.pool.Query(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var orders []domain.Order
	for rows.Next() {
		var ord domain.Order
		err := rows.Scan(
			&ord.ID, &ord.Phone, &ord.CustomerName, &ord.Address, &ord.Note, &ord.Status, &ord.AdminNote, &ord.TotalAmount, &ord.CreatedAt, &ord.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		orders = append(orders, ord)
	}
	return orders, rows.Err()
}

func (r *OrderRepository) Get(ctx context.Context, id string) (domain.Order, bool, error) {
	query := "SELECT id, phone, customer_name, address, note, status, admin_note, total_amount, created_at, updated_at FROM orders WHERE id = $1"

	var ord domain.Order
	err := r.pool.QueryRow(ctx, query, id).Scan(
		&ord.ID, &ord.Phone, &ord.CustomerName, &ord.Address, &ord.Note, &ord.Status, &ord.AdminNote, &ord.TotalAmount, &ord.CreatedAt, &ord.UpdatedAt,
	)

	if err != nil {
		if strings.Contains(err.Error(), "no rows in result set") {
			return domain.Order{}, false, nil
		}
		return domain.Order{}, false, err
	}
	return ord, true, nil
}

func (r *OrderRepository) Create(ctx context.Context, ord domain.Order) (domain.Order, error) {
	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return domain.Order{}, err
	}
	defer func() {
		_ = tx.Rollback(ctx)
	}()

	query := `INSERT INTO orders (id, phone, customer_name, address, note, status, admin_note, total_amount, created_at, updated_at)
	VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
	RETURNING id, phone, customer_name, address, note, status, admin_note, total_amount, created_at, updated_at`

	err = tx.QueryRow(ctx, query,
		ord.ID, ord.Phone, ord.CustomerName, ord.Address, ord.Note, ord.Status, ord.AdminNote, ord.TotalAmount, ord.CreatedAt, ord.UpdatedAt,
	).Scan(&ord.ID, &ord.Phone, &ord.CustomerName, &ord.Address, &ord.Note, &ord.Status, &ord.AdminNote, &ord.TotalAmount, &ord.CreatedAt, &ord.UpdatedAt)
	if err != nil {
		return domain.Order{}, err
	}

	createdItems := make([]domain.OrderItem, 0, len(ord.Items))
	for _, item := range ord.Items {
		if item.OrderID == "" {
			item.OrderID = ord.ID
		}

		itemQuery := `INSERT INTO order_items (id, order_id, product_id, product_title, product_subtitle, size_code, size_label, selected_attrs, quantity, unit_price, variant_image_url, created_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, now())
		RETURNING id, order_id, product_id, product_title, product_subtitle, size_code, size_label, selected_attrs, quantity, unit_price, variant_image_url`

		var selectedAttrsJSON []byte
		err = tx.QueryRow(ctx, itemQuery,
			item.ID, item.OrderID, item.ProductID, item.ProductTitle, item.ProductSubtitle,
			item.SizeCode, item.SizeLabel, jsonOrNull(item.SelectedAttrs),
			item.Quantity, item.UnitPrice, item.VariantImageURL,
		).Scan(
			&item.ID, &item.OrderID, &item.ProductID, &item.ProductTitle, &item.ProductSubtitle,
			&item.SizeCode, &item.SizeLabel, &selectedAttrsJSON,
			&item.Quantity, &item.UnitPrice, &item.VariantImageURL,
		)
		if len(selectedAttrsJSON) > 0 {
			_ = json.Unmarshal(selectedAttrsJSON, &item.SelectedAttrs)
		}
		if err != nil {
			return domain.Order{}, err
		}

		createdItems = append(createdItems, item)
	}

	if err = tx.Commit(ctx); err != nil {
		return domain.Order{}, err
	}

	ord.Items = createdItems

	return ord, err
}

func (r *OrderRepository) UpdateStatus(ctx context.Context, id string, status string, adminNote *string) error {
	result, err := r.pool.Exec(ctx, "UPDATE orders SET status = $1, admin_note = $2, updated_at = now() WHERE id = $3", status, adminNote, id)
	if err != nil {
		return err
	}
	if result.RowsAffected() == 0 {
		return errors.New("order not found")
	}
	return nil
}

func (r *OrderRepository) AddItem(ctx context.Context, item domain.OrderItem) (domain.OrderItem, error) {
	query := `INSERT INTO order_items (id, order_id, product_id, product_title, product_subtitle, size_code, size_label, selected_attrs, quantity, unit_price, variant_image_url, created_at)
	VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, now())
	RETURNING id, order_id, product_id, product_title, product_subtitle, size_code, size_label, selected_attrs, quantity, unit_price, variant_image_url`

	var selectedAttrsJSON []byte
	err := r.pool.QueryRow(ctx, query,
		item.ID, item.OrderID, item.ProductID, item.ProductTitle, item.ProductSubtitle,
		item.SizeCode, item.SizeLabel, jsonOrNull(item.SelectedAttrs),
		item.Quantity, item.UnitPrice, item.VariantImageURL,
	).Scan(
		&item.ID, &item.OrderID, &item.ProductID, &item.ProductTitle, &item.ProductSubtitle,
		&item.SizeCode, &item.SizeLabel, &selectedAttrsJSON,
		&item.Quantity, &item.UnitPrice, &item.VariantImageURL,
	)
	if len(selectedAttrsJSON) > 0 {
		_ = json.Unmarshal(selectedAttrsJSON, &item.SelectedAttrs)
	}
	return item, err
}

func (r *OrderRepository) GetItems(ctx context.Context, orderID string) ([]domain.OrderItem, error) {
	query := `SELECT id, order_id, product_id, product_title, product_subtitle, size_code, size_label, selected_attrs, quantity, unit_price, variant_image_url
	FROM order_items WHERE order_id = $1`

	rows, err := r.pool.Query(ctx, query, orderID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var items []domain.OrderItem
	for rows.Next() {
		var item domain.OrderItem
		var selectedAttrsJSON []byte
		err := rows.Scan(
			&item.ID, &item.OrderID, &item.ProductID, &item.ProductTitle, &item.ProductSubtitle,
			&item.SizeCode, &item.SizeLabel, &selectedAttrsJSON,
			&item.Quantity, &item.UnitPrice, &item.VariantImageURL,
		)
		if err != nil {
			return nil, err
		}
		if len(selectedAttrsJSON) > 0 {
			_ = json.Unmarshal(selectedAttrsJSON, &item.SelectedAttrs)
		}
		items = append(items, item)
	}
	return items, rows.Err()
}

func jsonOrNull(v any) any {
	if v == nil {
		return nil
	}
	b, _ := json.Marshal(v)
	return b
}
