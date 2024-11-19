package test

import (
	"database/sql"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"go.uber.org/zap"

	"github.com/sevings/mindwell-server/utils"
)

func TestPubSub(t *testing.T) {
	logger, err := zap.NewDevelopment()
	require.NoError(t, err)

	config := utils.LoadConfig("../configs/server")
	connString := utils.ConnectionString(config)

	t.Run("NewPubSub", func(t *testing.T) {
		ps := utils.NewPubSub(connString, logger)
		require.NotNil(t, ps)
		ps.Stop()
	})

	t.Run("Subscribe and receive notification", func(t *testing.T) {
		ps := utils.NewPubSub(connString, logger)
		require.NotNil(t, ps)
		defer ps.Stop()

		var wg sync.WaitGroup
		wg.Add(1)

		var receivedMessage []byte
		ps.Subscribe("test_channel", func(msg []byte) {
			receivedMessage = msg
			wg.Done()
		})

		ps.Start()

		sendTestPubSub(t, db, "test_channel", "test message")

		done := make(chan struct{})
		go func() {
			wg.Wait()
			close(done)
		}()

		select {
		case <-done:
			require.Equal(t, []byte("test message"), receivedMessage)
		case <-time.After(5 * time.Second):
			t.Fatal("Timeout waiting for notification")
		}
	})

	t.Run("Multiple subscribers", func(t *testing.T) {
		ps := utils.NewPubSub(connString, logger)
		require.NotNil(t, ps)
		defer ps.Stop()

		var wg sync.WaitGroup
		wg.Add(2)

		messages := make([][]byte, 0, 2)
		var messagesMu sync.Mutex

		handler := func(msg []byte) {
			messagesMu.Lock()
			messages = append(messages, msg)
			messagesMu.Unlock()
			wg.Done()
		}

		ps.Subscribe("test_channel", handler)
		ps.Subscribe("test_channel", handler)

		ps.Start()

		sendTestPubSub(t, db, "test_channel", "test message")

		done := make(chan struct{})
		go func() {
			wg.Wait()
			close(done)
		}()

		select {
		case <-done:
			require.Len(t, messages, 2)
			require.Equal(t, []byte("test message"), messages[0])
			require.Equal(t, []byte("test message"), messages[1])
		case <-time.After(5 * time.Second):
			t.Fatal("Timeout waiting for notifications")
		}
	})

	t.Run("Stop PubSub", func(t *testing.T) {
		ps := utils.NewPubSub(connString, logger)
		require.NotNil(t, ps)

		ps.Start()
		time.Sleep(10 * time.Millisecond)
		ps.Stop()
	})
}

func sendTestPubSub(t *testing.T, db *sql.DB, channel, message string) {
	_, err := db.Exec("SELECT pg_notify($1, $2)", channel, message)
	require.NoError(t, err)
}
