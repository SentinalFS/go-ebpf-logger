package utils

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"strings"
	"time"

	"go-ebp-logger/models"
)



func PrintBaseEvent(data []byte) {
    var event models.BaseEventData
    if err := binary.Read(bytes.NewReader(data), binary.LittleEndian, &event); err != nil {
        fmt.Printf("Failed to decode event: %v\n", err)
        return
    }
	ts := time.Unix(0, int64(event.Timestamp))
	fmt.Printf("[%s] PID=%d, UID=%d, CgroupId=%d, Operation=%s, File=%s, Process=%s\n",
		ts.Format("2006-01-02 15:04:05.000"),
		event.Pid,
		event.Uid,
		event.CgroupId,
		strings.TrimRight(string(event.Otype[:]), "\x00"),
		strings.TrimRight(string(event.Filename[:]), "\x00"),
		strings.TrimRight(string(event.Comm[:]), "\x00"),
	)
}

func PrintRenameEvent(data []byte) {
    var event models.RenameData
    if err := binary.Read(bytes.NewReader(data), binary.LittleEndian, &event); err != nil {
        fmt.Printf("Failed to decode rename event: %v\n", err)
        return
    }
	ts := time.Unix(0, int64(event.Timestamp))
    fmt.Printf("[%s] PID=%d, UID=%d, CgroupId=%d, OldFile=%s, NewFile=%s, Process=%s\n",
		ts.Format("2006-01-02 15:04:05.000"),
        event.Pid,
        event.Uid,
        event.CgroupId,
        strings.TrimRight(string(event.OldFileName[:]), "\x00"),
        strings.TrimRight(string(event.NewFileName[:]), "\x00"),
        strings.TrimRight(string(event.Comm[:]), "\x00"),
    )
}