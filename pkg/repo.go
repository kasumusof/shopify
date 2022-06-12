package pkg

import (
	"context"
	"log"
	"os"
	"time"

	sqr "github.com/Masterminds/squirrel"
	"github.com/gofrs/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"github.com/pkg/errors"
)

var (
	db   *sqlx.DB
	psql = sqr.StatementBuilder.PlaceholderFormat(sqr.Dollar)
)

var (
	Port     string
	dbname   string
	user     string
	password string
	dbhost   string
)

func init() {
	err := godotenv.Load()
	if err != nil {
		log.Println(err)
	}

	Port = os.Getenv("PORT")
	if Port == "" {
		Port = "8080"
	}
	dbname = os.Getenv("DBNAME")
	if dbname == "" {
		dbname = "postgres"
	}

	user = os.Getenv("DBUSER")
	if user == "" {
		user = "postgres"
	}

	password = os.Getenv("PASSWORD")
	if password == "" {
		password = "postgres"
	}

	dbhost = os.Getenv("DBHOST")
	if dbhost == "" {
		dbhost = "localhost"
	}
	// log.Println("port: ", Port, "dbname: ", dbname, "user: ", user, "password: ", password, "dbhost: ", dbhost)
	db, err = sqlx.Connect("postgres", "user="+user+" dbname="+dbname+" password="+password+" host="+dbhost+" sslmode=require")
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

func listItems(ctx context.Context, statusParam string) ([]Item, error) {
	var items []Item
	var subQuery = "*"
	if statusParam != "" {
		subQuery = "id, name, quantity, description, created_at, updated_at, (select comment from deleted_comments where item = items.id) as comment "
	}
	q := psql.Select(subQuery).From("items")
	// i would just check if the param was passed (without checking its value) and if it was, add the where clause
	if statusParam != "" {
		q = q.Where(sqr.NotEq{"deleted_at": nil})
	} else {
		q = q.Where(sqr.Eq{"deleted_at": nil})
	}

	q = q.OrderBy("created_at desc")

	toSql, args, err := q.ToSql()
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

func deleteItem(ctx context.Context, id string, comment *DeletedComment) error {
	newTime := time.Now()
	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		return errors.Wrap(err, "failed to begin transaction")
	}

	err = func() error {
		// update the item
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
		result, err := tx.ExecContext(ctx, toSql, args...)
		if err != nil {
			return errors.Wrap(err, "failed to delete item")
		}
		rowsAffected, _ := result.RowsAffected()
		if rowsAffected == 0 {
			return errors.New("nothing to delete")
		}
		// create the deleted comment
		toSql, args, err = psql.Insert("deleted_comments").
			Columns("id", "item", "comment", "updated_at", "created_at").
			Values(uuid.Must(uuid.NewV4()), id, comment.Comment, newTime, newTime).
			ToSql()
		if err != nil {
			return errors.Wrap(err, "failed to build query")
		}
		result, err = tx.ExecContext(ctx, toSql, args...)
		if err != nil {
			return errors.Wrap(err, "failed to create deleted comment")
		}
		rowsAffected, _ = result.RowsAffected()
		if rowsAffected == 0 {
			return errors.New("nothing was created")
		}

		return tx.Commit()
	}()
	if err != nil {
		if err := tx.Rollback(); err != nil {
			return errors.Wrap(err, "failed to rollback transaction")
		}
		return errors.Wrap(err, "failed to delete item")
	}
	return nil
}

func updateItem(ctx context.Context, id string, item *Item) error {
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
		return errors.Wrap(err, "failed to build query")
	}
	result, err := db.ExecContext(ctx, toSql, args...)
	if err != nil {
		return errors.Wrap(err, "failed to update item")
	}
	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		return errors.New("nothing to update")
	}
	return nil
}

func unArchiveItem(ctx context.Context, id string) error {
	newTime := time.Now()
	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		return errors.Wrap(err, "failed to begin transaction")
	}

	err = func() error {
		// update the item
		toSql, args, err := psql.
			Update("items").
			SetMap(map[string]interface{}{
				"deleted_at": nil,
				"updated_at": newTime,
			}).
			Where(sqr.Eq{"id": id}).
			ToSql()
		if err != nil {
			return errors.Wrap(err, "failed to build query")
		}
		result, err := tx.ExecContext(ctx, toSql, args...)
		if err != nil {
			return errors.Wrap(err, "failed to delete item")
		}
		rowsAffected, _ := result.RowsAffected()
		if rowsAffected == 0 {
			return errors.New("nothing to delete")
		}
		// delete the comment associated with the item
		toSql, args, err = psql.Delete("deleted_comments").
			Where(sqr.Eq{"item": id}).ToSql()
		if err != nil {
			return errors.Wrap(err, "failed to build query")
		}
		result, err = tx.ExecContext(ctx, toSql, args...)
		if err != nil {
			return errors.Wrap(err, "failed deleted comment")
		}
		rowsAffected, _ = result.RowsAffected()
		if rowsAffected == 0 {
			return errors.New("nothing was deleted")
		}

		return tx.Commit()
	}()
	if err != nil {
		if err := tx.Rollback(); err != nil {
			return errors.Wrap(err, "failed to rollback transaction")
		}
		return errors.Wrap(err, "failed to delete comment")
	}
	return nil
}
