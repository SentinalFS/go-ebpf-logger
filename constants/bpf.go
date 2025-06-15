package constants


var ProgsToFuncs = map[string]string{
	"trace_read":   "vfs_read",
	"trace_write":  "vfs_write",
	"trace_rename": "vfs_rename",
	"trace_delete": "vfs_unlink",
}

var InodeMapName = "monitored_inodes"

var BpfObjName = "monitor.bpf.o"