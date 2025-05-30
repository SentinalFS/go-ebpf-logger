package main

import (
	"bytes"
	"encoding/binary"
	"log"
	"os"
	"os/signal"
	"syscall"

	"go-ebp-logger/internal/event"

	"github.com/cilium/ebpf"
	"github.com/cilium/ebpf/link"
	"github.com/cilium/ebpf/ringbuf"
)

func main() {
	bpfObj := "monitor.bpf.o"
	spec, err := ebpf.LoadCollectionSpec(bpfObj)
	if err != nil {
		log.Fatalf("Failed to load BPF spec: %v", err)
	}

	coll, err := ebpf.NewCollection(spec)
	if err != nil {
		log.Fatalf("Failed to load BPF collection: %v", err)
	}
	defer coll.Close()

	prog := coll.Programs["trace_read"]
	if prog == nil {
		log.Fatalf("Program 'trace_read' not found")
	}

	kp, err := link.Kprobe("vfs_read", prog, nil)
	if err != nil {
		log.Fatalf("Failed to attach kprobe: %v", err)
	}
	defer kp.Close()

	events := coll.Maps["events"]
	if events == nil {
		log.Fatalf("Map 'events' not found")
	}

	rd, err := ringbuf.NewReader(events)
	if err != nil {
		log.Fatalf("Failed to create ring buffer reader: %v", err)
	}
	defer rd.Close()

	log.Println("Listening for events... Press Ctrl+C to stop.")
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)

loop:
	for {
		select {
		case <-sigCh:
			break loop
		default:
			record, err := rd.Read()
			if err != nil {
				log.Printf("Failed to read from ring buffer: %v", err)
				continue
			}
			var e event.Data
			if err := binary.Read(bytes.NewReader(record.RawSample), binary.LittleEndian, &e); err == nil {
				event.Print(e)
			}
		}
	}
}
