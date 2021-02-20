package twitch

import "fmt"

func (op *operator) JoinChannel() error {
	if op.conn == nil {
		return fmt.Errorf("failed to join channel; %w", ErrNilConnection)
	}

	op.conn.Write([]byte(fmt.Sprintf("%s oauth:%s\r\n", CmdPrefixPassword, op.creds.token)))
	op.conn.Write([]byte(fmt.Sprintf("%s %s\r\n", CmdPrefixNickname, op.creds.username)))
	op.conn.Write([]byte(fmt.Sprintf("%s #%s\r\n", CmdPrefixJoinChannel, op.channel)))

	return nil
}
