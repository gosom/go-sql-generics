package main

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	_ "github.com/lib/pq"
)

func main() {
	ctx := context.Background()
	dsn := "postgres://localhost/demo?sslmode=disable&user=demo&password=demo"
	db, err := getDbCon(ctx, dsn)
	if err != nil {
		panic(err)
	}

	const q = "SELECT id, title, content, created_at  FROM notes"
	notes, err := Query[Note](ctx, db, q, nil)
	if err != nil {
		panic(err)
	}
	for i := range notes {
		fmt.Println(notes[i])
	}

}

type Note struct {
	ID        int
	Title     string
	Content   string
	CreatedAt time.Time
}

func (o *Note) String() string {
	return fmt.Sprintf("id=%d title=%q content=%q createdAt=%s", o.ID, o.Title, o.Content, o.CreatedAt)
}

func (o *Note) DbBind() []any {
	return []any{&o.ID, &o.Title, &o.Content, &o.CreatedAt}
}

// ===============================================================

func getDbCon(ctx context.Context, dsn string) (*sql.DB, error) {
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		return db, err
	}
	if err := db.PingContext(ctx); err != nil {
		return db, err
	}
	return db, nil
}

type DBTx interface {
	QueryContext(ctx context.Context, q string, args ...any) (*sql.Rows, error)
}

type Bindable[T any] interface {
	*T
	DbBind() []any
}

func Query[T any, PT Bindable[T]](ctx context.Context, d DBTx, q string, args ...any) ([]T, error) {
	rows, err := d.QueryContext(ctx, q)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []T
	for rows.Next() {
		var item T
		var ptr PT
		ptr = &item
		if err := rows.Scan(ptr.DbBind()...); err != nil {
			return nil, err
		}
		items = append(items, item)
	}
	return items, rows.Err()
}
