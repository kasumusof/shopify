package pkg

import (
	"context"
	"time"

	sqr "github.com/Masterminds/squirrel"
	"github.com/gofrs/uuid"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"github.com/pkg/errors"
)

var (
	db   *sqlx.DB
	err  error
	psql = sqr.StatementBuilder.PlaceholderFormat(sqr.Dollar)
)

func init() {
	db, err = sqlx.Connect("postgres", "user=postgres dbname=shopify_development sslmode=disable")
	if err != nil {
		panic(err)
	}
}

func getItem(ctx context.Context, id string) (*Item, error) {
	var item Item
	toSql, args, err := psql.Select("*").From("items").Where(sqr.Eq{"id": id}).ToSql()

	if err != nil {
		return nil, errors.Wrap(err, "failed to build query")
	}

	err = db.GetContext(ctx, &item, toSql, args...)

	if err != nil {
		return nil, errors.Wrap(err, "failed to get item")
	}

	return &item, nil
}

func listItems(ctx context.Context) ([]Item, error) {
	var items []Item
	toSql, args, err := psql.Select("*").From("items").ToSql()
	if err != nil {
		return nil, errors.Wrap(err, "failed to build query")
	}

	err = db.SelectContext(ctx, &items, toSql, args...)
	if err != nil {
		return nil, errors.Wrap(err, "failed to get items")
	}

	return items, nil
}

func createItem(ctx context.Context, item *Item) (*Item, error) {
	newTime := time.Now()
	item.ID = uuid.Must(uuid.NewV4())
	item.CreatedAt = newTime
	item.UpdatedAt = newTime

	toSql, args, err := psql.
		Insert("items").
		Columns("id", "name", "quantity", "description", "created_at", "updated_at").
		Values(item.ID, item.Name, item.Quantity, item.Description, item.CreatedAt, item.UpdatedAt).
		ToSql()
	if err != nil {
		return nil, errors.Wrap(err, "failed to build query")
	}

	rows, err := db.ExecContext(ctx, toSql, args...)
	if err != nil {
		return nil, errors.Wrap(err, "failed to create item")
	}

	rowsAffected, _ := rows.RowsAffected()
	if rowsAffected == 0 {
		return nil, errors.New("nothing to create")
	}

	return item, nil
}

func deleteItem(ctx context.Context, id string) error {
	newTime := time.Now()
	toSql, args, err := psql.
		Update("items").
		SetMap(map[string]interface{}{
			"deleted_at": newTime,
			"updated_at": newTime,
		}).
		Where(sqr.Eq{"id": id}).
		ToSql()
	if err != nil {
		return errors.Wrap(err, "failed to build query")
	}
	result, err := db.ExecContext(ctx, toSql, args...)
	if err != nil {
		return errors.Wrap(err, "failed to delete item")
	}
	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		return errors.New("nothing to delete")
	}
	return nil
}

func updatedItem(ctx context.Context, id string, item *Item) (*Item, error) {
	newTime := time.Now()
	toSql, args, err := psql.
		Update("items").
		SetMap(map[string]interface{}{
			"name":        item.Name,
			"quantity":    item.Quantity,
			"description": item.Description,
			"updated_at":  newTime,
		}).
		Where(sqr.Eq{"id": id}).
		ToSql()
	if err != nil {
		return nil, errors.Wrap(err, "failed to build query")
	}
	result, err := db.ExecContext(ctx, toSql, args...)
	if err != nil {
		return nil, errors.Wrap(err, "failed to update item")
	}
	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		return nil, errors.New("nothing to update")
	}
	return nil, nil
}
