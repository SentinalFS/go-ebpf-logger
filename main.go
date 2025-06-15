package main

import (
	"go-ebp-logger/bpf"
	"go-ebp-logger/constants"
	"go-ebp-logger/utils"
	"os"
	"os/signal"
	"syscall"

	"github.com/cilium/ebpf/ringbuf"

	"fmt"
)

func main() {
	events, renameEvents, cleanup := bpf.SetupBPF(constants.BpfObjName)
	defer cleanup()

	baseEventsRB, err := ringbuf.NewReader(events)
	if err != nil {
		fmt.Printf("Failed to create ring buffer reader: %v", err)
		panic(err)
	}
	defer baseEventsRB.Close()

	renameEventsRB, err := ringbuf.NewReader(renameEvents)
	if err != nil {
		fmt.Printf("Failed to create ring buffer reader: %v", err)
		panic(err)
	}
	defer renameEventsRB.Close()

	done := make(chan struct{})
    fmt.Println("Waiting for events... Press Ctrl+C to stop.")
    sigCh := make(chan os.Signal, 1)
    signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		for {
			select {
			case <-done:
				return
			default:
				record, err := baseEventsRB.Read()
				if err != nil {
					fmt.Printf("Failed to read from ring buffer 1: %v", err)
					continue
				}
				utils.PrintBaseEvent(record.RawSample)
			}
		}
	}()

	go func() {
		for {
			select {
			case <-done:
				return
			default:
				record, err := renameEventsRB.Read()
				if err != nil {
					fmt.Printf("Failed to read from ring buffer 1: %v", err)
					continue
				}
				utils.PrintRenameEvent(record.RawSample)
			}
		}
	}()

	<-sigCh
	close(done)
	fmt.Println("\nExiting...")
}
