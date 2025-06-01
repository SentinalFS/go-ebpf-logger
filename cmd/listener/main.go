// main.go
package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"go-ebp-logger/internal/bpfmanager" 
	"go-ebp-logger/internal/event"     

	"github.com/cilium/ebpf/ringbuf" 
)

func main() {
	
	done := make(chan struct{})
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage: %s <bpf-object-file>\n", os.Args[0])
	}
	flag.Parse()
	if flag.NArg() < 1 {
		flag.Usage()
		os.Exit(1)
	}
	bpfObjPath := flag.Arg(0) 

	bpfManager, err := bpfmanager.NewBPFManager(bpfObjPath)
	if err != nil {
		log.Fatalf("Failed to initialize BPF manager: %v", err)
	}
	defer bpfManager.Close()


	if err := bpfManager.AttachKprobe("vfs_read", "trace_read"); err != nil {
		log.Fatalf("Failed to attach kprobe for vfs_read: %v", err)
	}


	if err := bpfManager.PinMap("monitored_inodes", "/sys/fs/bpf/monitored_inode"); err != nil {
		log.Fatalf("Failed to pin 'monitored_inodes' map: %v", err)
	}


	if err := bpfManager.InitRingBufferReader("events"); err != nil {
		log.Fatalf("Failed to initialize event reader: %v", err)
	}
	rd := bpfManager.GetEventReader()

	fmt.Println("Waiting for events... Press Ctrl+C to stop.")


	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		for {
			select {
			case <-done: 
				return 
			default:
				record, err := rd.Read()
				if err != nil {
					if ringbuf.Is  (err) {
						log.Println("Ring buffer reader closed, exiting event processing goroutine.")
						return 
					}
					log.Printf("Failed to read from ring buffer: %v", err)
					continue 
				}
				event.Print(record.RawSample)
			}
		}
	}()

	<-sigCh
	close(done)
	fmt.Println("\nExiting...")
}