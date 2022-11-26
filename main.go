package main

import (
	"context"
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"
)

func main() {
	ctx := context.Background()
	dsn := "postgres://localhost/demo?sslmode=disable&user=demo&password=demo"
	db, err := getDbCon(ctx, dsn)
	if err != nil {
		panic(err)
	}

	alltodos, err := selectAll(ctx, db)
	if err != nil {
		panic(err)
	}
	for _, todo := range alltodos {
		fmt.Println(todo)
	}

	fmt.Println("============ v2 =======================")
	alltodos2, err := selectAllV2(ctx, db)
	if err != nil {
		panic(err)
	}
	for _, todo := range alltodos2 {
		fmt.Println(todo)
	}
	fmt.Println("============ v3 =======================")
	alltodos3, err := selectAllV3[Note](ctx, db)
	if err != nil {
		panic(err)
	}
	for _, todo := range alltodos3 {
		fmt.Println(todo)
	}

}

type Note struct {
	ID      int
	Title   string
	Content string
}

func (o *Note) String() string {
	return fmt.Sprintf("id=%d title=%q content=%q", o.ID, o.Title, o.Content)
}

func (o *Note) DbBind() []any {
	return []any{&o.ID, &o.Title, &o.Content}
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

const q = `select id, title, content from notes`

// selectAll selects all notes.
func selectAll(ctx context.Context, db *sql.DB) ([]Note, error) {
	rows, err := db.QueryContext(ctx, q)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []Note
	for rows.Next() {
		var item Note
		if err := rows.Scan(&item.ID, &item.Title, &item.Content); err != nil {
			return nil, err
		}
		items = append(items, item)
	}
	return items, rows.Err()
}

// selectAllV2 is my second attempt to make it more generic
func selectAllV2(ctx context.Context, db *sql.DB) ([]Note, error) {
	rows, err := db.QueryContext(ctx, q)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []Note
	for rows.Next() {
		var item Note
		if err := rows.Scan(item.DbBind()...); err != nil {
			return nil, err
		}
		items = append(items, item)
	}
	return items, rows.Err()
}

type Bindable[T any] interface {
	*T
	DbBind() []any
}

func selectAllV3[T any, PT Bindable[T]](ctx context.Context, db *sql.DB) ([]T, error) {
	rows, err := db.QueryContext(ctx, q)
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
