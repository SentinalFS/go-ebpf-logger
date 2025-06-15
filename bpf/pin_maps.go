package bpf

import (
	"os"

	"github.com/cilium/ebpf"

	"fmt"

	"go-ebp-logger/constants"
)

func pinMaps(m *ebpf.Map) error {
	if m != nil {
		path := "/sys/fs/bpf/" + constants.InodeMapName
		if err := os.Remove(path); err != nil && !os.IsNotExist(err) {
			fmt.Printf("Warning: failed to remove existing pin at %s: %v", path, err)
		}
		if err := m.Pin(path); err != nil {
            fmt.Printf("Failed to pin map %s to %s: %v", constants.InodeMapName,path, err)
        }
		fmt.Printf("Map %s pinned succesfully", constants.InodeMapName)
		return nil
	}

	fmt.Printf("Map %s not found", constants.InodeMapName)
	return fmt.Errorf("no map was pinned")
}