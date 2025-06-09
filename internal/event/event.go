package event

import (
	"fmt"
	"strings"
	"time"
	"models/models"
)


func Print(e models.Data) {
	ts := time.Unix(0, int64(e.Timestamp))
	fmt.Printf("[%s] PID=%d, UID=%d, Operation=%s, File=%s, Process=%s\n",
		ts.Format("15:04:05"),
		e.Pid,
		e.Uid,
		strings.TrimRight(string(e.Otype[:]), "\x00"),
		strings.TrimRight(string(e.Filename[:]), "\x00"),
		strings.TrimRight(string(e.Comm[:]), "\x00"),
	)
}
