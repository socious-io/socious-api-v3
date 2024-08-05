package database

import (
	"context"
	"database/sql"
	"fmt"
	"os"
	"path/filepath"
	"reflect"
	"sync"

	"github.com/jmoiron/sqlx"
)

var (
	queryCache = make(map[string]string)
	cacheLock  sync.RWMutex
)

// Create inserts a new record into the database
func Create(ctx context.Context, queryName string, model Model) (sql.Result, error) {
	return ExecuteQuery(ctx, queryName, extractFields(model)...)
}

// Get retrieves multiple records from the database with pagination
func Get(ctx context.Context, queryName string, dest interface{}, ids ...string) error {
	_, err := cb.Execute(func() (interface{}, error) {
		q, err := LoadQuery(queryName)
		if err != nil {
			return nil, fmt.Errorf("could not load query: %v", err)
		}

		db := GetDB()

		if db == nil {
			return nil, fmt.Errorf("database not connected")
		}

		query, args, err := sqlx.In(q, ids)
		if err != nil {
			return nil, err
		}

		query = db.Rebind(query)
		destType := reflect.TypeOf(dest)
		if destType.Kind() == reflect.Ptr && destType.Elem().Kind() == reflect.Slice {
			if err := db.Select(dest, query, args...); err != nil {
				return nil, err
			}
		} else {
			if err := db.Get(dest, query, args...); err != nil {
				return nil, err
			}
		}
		return nil, nil
	})
	return err
}

// Update updates an existing record in the database
func Update(ctx context.Context, queryName string, model Model, args ...interface{}) (sql.Result, error) {
	return ExecuteQuery(ctx, queryName, append(extractFields(model), args...)...)
}

// LoadQuery reads the SQL query from the file and caches it.
func LoadQuery(queryName string) (string, error) {
	cacheLock.RLock()
	if query, found := queryCache[queryName]; found {
		cacheLock.RUnlock()
		return query, nil
	}
	cacheLock.RUnlock()

	// Load query from file
	filePath := filepath.Join(sqlDir, queryName+".sql")
	queryBytes, err := os.ReadFile(filePath)
	if err != nil {
		return "", err
	}
	query := string(queryBytes)

	// Cache the query
	cacheLock.Lock()
	queryCache[queryName] = query
	cacheLock.Unlock()

	return query, nil
}

// ExecuteQuery is a general function to execute a write operation (INSERT, UPDATE, DELETE) with a circuit breaker
func ExecuteQuery(ctx context.Context, queryName string, args ...interface{}) (sql.Result, error) {
	var result sql.Result
	_, err := cb.Execute(func() (interface{}, error) {
		query, err := LoadQuery(queryName)
		if err != nil {
			return nil, fmt.Errorf("could not load query: %v", err)
		}

		db := GetDB()
		if db == nil {
			return nil, fmt.Errorf("database not connected")
		}

		result, err = db.ExecContext(ctx, query, args...)
		if err != nil {
			return nil, err
		}

		return result, nil
	})

	return result, err
}

// QueryRows is a general function to execute a read operation that returns multiple rows with a circuit breaker
func QueryRows(ctx context.Context, queryName string, args ...interface{}) (*sqlx.Rows, error) {
	rows, err := cb.Execute(func() (interface{}, error) {
		query, err := LoadQuery(queryName)
		if err != nil {
			return nil, fmt.Errorf("could not load query: %v", err)
		}

		db := GetDB()

		if db == nil {
			return nil, fmt.Errorf("database not connected")
		}
		return db.QueryxContext(ctx, query, args...)
	})

	if err != nil {
		return nil, err
	}

	return rows.(*sqlx.Rows), nil
}

// QuerSelect is a general function to execute a read operation that returns multiple rows with a circuit breaker
func QuerySelect(queryName string, dest interface{}, args ...interface{}) error {
	_, err := cb.Execute(func() (interface{}, error) {
		query, err := LoadQuery(queryName)
		if err != nil {
			return nil, fmt.Errorf("could not load query: %v", err)
		}

		db := GetDB()

		if db == nil {
			return nil, fmt.Errorf("database not connected")
		}
		return nil, db.Select(dest, query, args...)
	})
	return err
}

// extractFields extracts the fields from a struct and returns them as a slice of interface{}
func extractFields(model interface{}) []interface{} {
	v := reflect.ValueOf(model).Elem()
	numFields := v.NumField()
	fields := make([]interface{}, numFields)

	for i := 0; i < numFields; i++ {
		fields[i] = v.Field(i).Interface()
	}

	return fields
}
