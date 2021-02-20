package twitch

import (
	"fmt"
	"net"
	"time"
)

const (
	_retryAttempts int = 5
)

func (op *operator) Connect() error {
	conn, err := net.Dial("tcp", fmt.Sprintf("%s:%s", op.endpoint, op.port))
	if err != nil {
		if conn, err = op.retry(_retryAttempts, 1); err != nil {
			return fmt.Errorf("failed to connect to Twitch server '%s:%s'; %w", op.endpoint, op.port, err)
		}
	}

	op.conn = conn
	return nil
}

func (op *operator) Disconnect() {
	close(op.chatObserver)
	op.conn.Close()
}

func (op *operator) retry(attempts, sleepSec int) (net.Conn, error) {
	conn, err := net.Dial("tcp", fmt.Sprintf("%s:%s", op.endpoint, op.port))
	if err != nil {
		if attempts--; attempts > 0 {
			time.Sleep(time.Duration(sleepSec) * time.Second)
			return op.retry(attempts, sleepSec*2)
		}
	}

	return conn, nil
}
