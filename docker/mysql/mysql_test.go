package mysql_test

import (
	"context"
	"database/sql"
	"fmt"
	"testing"

	"github.com/TDD-all-the-things/learn-by-testing/docker"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestDockerRunMySQL(t *testing.T) {
	t.SkipNow()
	dbname, username, password := "docker", "root", "root123"
	image, port := "mysql:5.7.31", "3306/tcp"
	//env := []string{"MYSQL_ALLOW_EMPTY_PASSWORD=yes"}
	//env := []string{fmt.Sprintf("MYSQL_DATABASE=%s", dbname), fmt.Sprintf("MYSQL_USER=%s", username), fmt.Sprintf("MYSQL_PASSWORD=%s", password)}
	//env := []string{fmt.Sprintf("MYSQL_ROOT_PASSWORD=%s", password)}
	env := []string{fmt.Sprintf("MYSQL_DATABASE=%s", dbname), fmt.Sprintf("MYSQL_ROOT_PASSWORD=%s", password)}
	remove, info, err := docker.Run(image, []string{port}, env)

	defer remove()
	assert.NoError(t, err)
	//time.Sleep(time.Second * 5)

	require.True(t, info.Status)

	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?timeout=60s&charset=utf8mb4&collation=utf8mb4_general_ci",
		username, password, info.Hosts[port][0].IP, info.Hosts[port][0].Port, dbname)
	// assert.Equal(t, "0", dsn)

	db, err := sql.Open("mysql", dsn)
	assert.NoError(t, err)

	//ctx, canel := context.WithTimeout(context.Background(), time.Second*2)
	//defer canel()
	for err := db.Ping(); err != nil; {
		err = db.Ping()
	}
	//require.NoError(t, db.Ping())

	//db.SetConnMaxLifetime(time.Minute)
	//db.SetMaxOpenConns(1)
	//db.SetMaxIdleConns(1)

	conn, err := db.Conn(context.Background())
	require.NoError(t, err)
	defer conn.Close()

	row, err := conn.QueryContext(context.Background(), "show databases")
	assert.NoError(t, err)

	var database string

	for row.Next() {
		err = row.Scan(&database)
		assert.NoError(t, err)
		assert.Zero(t, database)
	}
}
