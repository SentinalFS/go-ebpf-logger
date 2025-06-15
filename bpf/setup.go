package bpf

import (
	"fmt"

	"github.com/cilium/ebpf"
	"github.com/cilium/ebpf/link"

	"go-ebp-logger/constants"
)

func SetupBPF(bpfObj string) (*ebpf.Map, *ebpf.Map, func()) {
    spec, err := ebpf.LoadCollectionSpec(bpfObj)
    if err != nil {
        fmt.Printf("Failed to load BPF collection spec: %v\n", err)
    }

    coll, err := ebpf.NewCollection(spec)
    if err != nil {
        fmt.Printf("Failed to load BPF collection: %v\n", err)
    }

    var links []*link.Link
    for progName, fn := range constants.ProgsToFuncs {
        prog := coll.Programs[progName]
        if prog == nil {
            fmt.Printf("Program '%s' not found\n", progName)
        }
        kp, err := link.Kprobe(fn, prog, nil)
        if err != nil {
            fmt.Printf("Failed to attach kprobe to %s: %v\n", fn, err)
        }
        links = append(links, &kp)
    }

	monitored_inode_map := coll.Maps["monitored_inodes"]
	err = pinMaps(monitored_inode_map)
	if err != nil {
		panic(err)
	}

    events := coll.Maps["events"]
    if events == nil {
        fmt.Printf("Map 'events' not found\n")
    }

	renameEvents := coll.Maps["rename_events"]
	if renameEvents == nil {
		fmt.Printf("Map 'rename_events' not found\n")
	}

    cleanup := func() {
        for _, l := range links {
            (*l).Close()
        }
        coll.Close()
    }

    return events, renameEvents, cleanup
}