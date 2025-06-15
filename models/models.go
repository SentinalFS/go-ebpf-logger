package models

type BaseEventData struct {
    Pid       uint32
    Uid       uint32
    Filename  [144]byte
    Comm      [32]byte
    Timestamp uint64
    CgroupId  uint64
    Otype     [16]byte
}

type RenameData struct {
    Pid         uint32
    Uid         uint32
    OldFileName [144]byte
    NewFileName [144]byte
    Comm        [32]byte
    Timestamp   uint64
    CgroupId    uint64
    Otype       [16]byte
}