package constant

import "fmt"

var (
	ChatNotReplyYetError    = fmt.Errorf("gpt not reply yet")
	ConcurrencyWriteDBError = fmt.Errorf("concurrency write db")
)
