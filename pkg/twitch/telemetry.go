package twitch

import "fmt"

func (op *operator) Ping() error {
	if _, err := op.conn.Write([]byte(fmt.Sprintf("PONG :%s \r\n", op.endpoint))); err != nil {
		return fmt.Errorf("%s; %w", ErrPing, err)
	}

	return nil
}
