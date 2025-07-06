package models

type BaseEventData struct {
	Pid            uint32
	Uid            uint32
	Filename       [176]byte
	ParentFilename [176]byte
	Inode          uint64
	Comm           [32]byte
	Timestamp      uint64
	CgroupId       uint64
	Otype          [16]byte
}

type RenameData struct {
	Pid               uint32
	Uid               uint32
	OldFileName       [176]byte
	OldParentFilename [176]byte
	NewFileName       [176]byte
	NewParentFilename [176]byte
	InodeOld          uint64
	InodeNew          uint64
	Comm              [32]byte
	Timestamp         uint64
	CgroupId          uint64
	Otype             [16]byte
}
