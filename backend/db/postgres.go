package db

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	"github.com/jackc/pgx/v4"

	"go.uber.org/zap"

	"github.com/jackc/pgx/v4/pgxpool"
)

const (
	DEFAULT_PAGE_SIZE = 10
)

type PostgresClient struct {
	ctx  context.Context
	Pool *pgxpool.Pool
	Url  string
}

func NewPostgresClient(ctx context.Context, url string) (*PostgresClient, error) {
	poolConfig, err := pgxpool.ParseConfig(url)
	if err != nil {
		return nil, err
	}
	pool, err := pgxpool.ConnectConfig(ctx, poolConfig)
	if err != nil {
		return nil, err
	}
	p := &PostgresClient{
		ctx:  ctx,
		Pool: pool,
		Url:  url,
	}
	return p, nil
}

// CreateTable Creates a table in the database
func (p *PostgresClient) CreateTable(tableName string, columns string) error {

	query := fmt.Sprintf("CREATE TABLE IF NOT EXISTS %s (%s)", tableName, columns)
	zlog.Info("CreateTable", zap.String("query", query))

	conn, err := p.Pool.Acquire(p.ctx)
	if err != nil {
		return err
	}
	defer conn.Release()

	_, err = conn.Exec(p.ctx, query)
	return err
}

// Insert Inserts a record into the database
func (p *PostgresClient) Insert(tableName string, columns []string, args []interface{}) (pgx.Row, error) {

	var argCount string
	for i := range columns {
		if i > 0 {
			argCount += ", "
		}
		argCount += "$" + strconv.Itoa(i+1)
	}

	query := fmt.Sprintf("INSERT INTO %s (%s) VALUES (%s)", tableName, strings.Join(columns, ","), argCount)
	zlog.Info("Insert", zap.String("query", query))

	conn, err := p.Pool.Acquire(p.ctx)
	if err != nil {
		return nil, err
	}
	defer conn.Release()

	return conn.QueryRow(p.ctx, query, args...), nil

}

// SelectOne Fetches one record from the database
func (p *PostgresClient) SelectOne(tableName string, columns string, condition string) pgx.Row {

	query := fmt.Sprintf("SELECT %s FROM %s WHERE %s AND Active = true", columns, tableName, condition)
	zlog.Info("SelectOne", zap.String("query", query))

	conn, err := p.Pool.Acquire(p.ctx)
	if err != nil {
		return nil
	}
	defer conn.Release()

	return conn.QueryRow(p.ctx, query)
}

// Select returns paginated records in the table
func (p *PostgresClient) Select(tableName string, columns string, condition string, page int, pageSize int) (pgx.Rows, error) {
	limit := DEFAULT_PAGE_SIZE
	if pageSize > 0 {
		limit = pageSize
	}

	offset := 0
	if page > 0 {
		offset = page * pageSize
	}

	query := fmt.Sprintf("SELECT %s FROM %s WHERE %s AND active = true LIMIT %d OFFSET %d", columns, tableName, condition, limit, offset)
	zlog.Info("Select", zap.String("query", query))

	conn, err := p.Pool.Acquire(p.ctx)
	if err != nil {
		return nil, err
	}
	defer conn.Release()

	return conn.Query(p.ctx, query)
}

// SelectAll returns all records in the table
func (p *PostgresClient) SelectAll(tableName string, columns string, page int, pageSize int) (pgx.Rows, error) {
	limit := DEFAULT_PAGE_SIZE
	if pageSize > 0 {
		limit = pageSize
	}

	offset := 0
	if page > 0 {
		offset = page * pageSize
	}

	query := fmt.Sprintf("SELECT %s FROM %s WHERE Active = true LIMIT %d OFFSET %d", columns, tableName, limit, offset)
	zlog.Info("SelectAll", zap.Any("query", query))

	conn, err := p.Pool.Acquire(p.ctx)
	if err != nil {
		return nil, err
	}
	defer conn.Release()

	return conn.Query(p.ctx, query)
}

// Update updates an item in the database
func (p *PostgresClient) Update(tableName string, updates string, condition string, args []interface{}) error {
	query := fmt.Sprintf("UPDATE %s SET %s WHERE %s AND Active = true", tableName, updates, condition)
	zlog.Info("Update", zap.Any("query", query))

	conn, err := p.Pool.Acquire(p.ctx)
	if err != nil {
		return err
	}
	defer conn.Release()

	result, err := conn.Exec(p.ctx, query, args...)
	if err != nil {
		return err
	}
	if result.RowsAffected() == 0 {
		return fmt.Errorf("no rows updated")
	}

	return nil
}

// Delete removes an item from the database
func (p *PostgresClient) Delete(tableName string, condition string) error {
	query := fmt.Sprintf("DELETE FROM %s WHERE %s", tableName, condition)
	zlog.Info("Delete", zap.Any("query", query))

	conn, err := p.Pool.Acquire(p.ctx)
	if err != nil {
		return err
	}
	defer conn.Release()

	_, err = conn.Exec(p.ctx, query)
	return err
}
