package mymutex

import (
	"context"
	"database/sql"
	"fmt"
	"os"
	"testing"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/vgarvardt/go-pg-adapter/sqladapter"
)

func getURI() string {
	uri := os.Getenv("MY_URI")
	if uri == "" {
		fmt.Println("Env variable MY_URI is required to run the tests")
		os.Exit(1)
	}

	return uri
}

func TestNew(t *testing.T) {
	db1, err := sql.Open("mysql", getURI())
	require.NoError(t, err)
	defer func() {
		assert.NoError(t, db1.Close())
	}()

	db2, err := sql.Open("mysql", getURI())
	require.NoError(t, err)
	defer func() {
		assert.NoError(t, db2.Close())
	}()

	conn1, err := db1.Conn(context.Background())
	require.NoError(t, err)
	defer func() {
		assert.NoError(t, conn1.Close())
	}()

	conn2, err := db2.Conn(context.Background())
	require.NoError(t, err)
	defer func() {
		assert.NoError(t, conn2.Close())
	}()

	adapter1 := sqladapter.NewConn(conn1)
	adapter2 := sqladapter.NewConn(conn2)

	lockName := fmt.Sprintf("lock_%d", time.Now().UnixNano())

	m1, err := New(adapter1)
	require.NoError(t, err)

	m2, err := New(adapter2)
	require.NoError(t, err)

	err = m1.Lock(lockName)
	require.NoError(t, err)

	success, err := m2.TryLock(lockName)
	require.NoError(t, err)
	assert.False(t, success)

	err = m1.Unlock(lockName)
	require.NoError(t, err)

	success, err = m2.TryLock(lockName)
	require.NoError(t, err)
	assert.True(t, success)

	success, err = m1.TryLock(lockName)
	require.NoError(t, err)
	assert.False(t, success)

	var m1AcquiredLock bool
	go func() {
		err := m1.Lock(lockName)
		require.NoError(t, err)
		m1AcquiredLock = true
	}()

	time.Sleep(time.Second)

	assert.False(t, m1AcquiredLock)

	err = m2.Unlock(lockName)
	require.NoError(t, err)

	time.Sleep(time.Second)

	assert.True(t, m1AcquiredLock)

	err = m1.Unlock(lockName)
	require.NoError(t, err)
}
