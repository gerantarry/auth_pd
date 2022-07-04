package mysql

import (
	"context"
	"database/sql"
	"fmt"
	"testing"
)

const (
	login    = "user1_login"
	password = "so61ty872nst"
	name     = "Carla"
)

func createDbClient() *sql.DB {
	db, _ := sql.Open(mysqlDriver, dataSourceName)
	return db
}

func TestNewUserStorage(t *testing.T) {
	db := createDbClient()
	if db == nil {
		t.Fatal()
	}
	storage := NewUserStorage(db)
	if storage == nil {
		t.FailNow()
	}
}

func TestUserStorage_Get(t *testing.T) {
	storage := NewUserStorage(createDbClient())
	user, err := storage.Get(context.Background(), login, password)
	if err != nil {
		t.Fatal(err)
	}
	if user.FirstName != name {
		fmt.Errorf("данные в записи не совпадают")
		t.FailNow()
	}

}
