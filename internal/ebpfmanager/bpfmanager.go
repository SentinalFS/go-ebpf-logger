package bpfmanager

import (
	"fmt"
	"log"
	"os"

	"github.com/cilium/ebpf"
	"github.com/cilium/ebpf/link"
	"github.com/cilium/ebpf/ringbuf"
)

type BPFManager struct {
	collection *ebpf.Collection
	reader     *ringbuf.Reader
	links      []link.Link 
}


func NewBPFManager(bpfObjPath string) (*BPFManager, error) {
	spec, err := ebpf.LoadCollectionSpec(bpfObjPath)
	if err != nil {
		return nil, fmt.Errorf("failed to load BPF spec from %s: %w", bpfObjPath, err)
	}
	coll, err := ebpf.NewCollection(spec)
	if err != nil {
		return nil, fmt.Errorf("failed to load BPF collection: %w", err)
	}
	return &BPFManager{
		collection: coll,
		links: make([]link.Link, 0), 
	}, nil
}



func (bm *BPFManager) AttachKprobe(kernelSymbol string, programName string) error {
	prog := bm.collection.Programs[programName]
	if prog == nil {
		return fmt.Errorf("eBPF program '%s' not found in collection", programName)
	}
	kp, err := link.Kprobe(kernelSymbol, prog, nil)
	if err != nil {
		return fmt.Errorf("failed to attach kprobe to '%s' using program '%s': %w", kernelSymbol, programName, err)
	}
	bm.links = append(bm.links, kp) 
	log.Printf("Successfully attached kprobe: %s -> %s", kernelSymbol, programName)
	return nil
}



func (bm *BPFManager) InitRingBufferReader(mapName string) error {
	eventsMap := bm.collection.Maps[mapName]
	if eventsMap == nil {
		return fmt.Errorf("eBPF map '%s' not found in collection", mapName)
	}
	rd, err := ringbuf.NewReader(eventsMap)
	if err != nil {
		return fmt.Errorf("failed to create ring buffer reader for map '%s': %w", mapName, err)
	}
	bm.reader = rd 
	log.Printf("Initialized ring buffer reader for map: %s", mapName)
	return nil
}



func (bm *BPFManager) GetEventReader() *ringbuf.Reader {
	return bm.reader
}



func (bm *BPFManager) PinMap(mapName string, pinPath string) error {
	bpfMap := bm.collection.Maps[mapName]
	if bpfMap == nil {
		log.Printf("Map '%s' not found in collection. Skipping pinning.", mapName)
		return nil
	}
	if err := os.Remove(pinPath); err != nil && !os.IsNotExist(err) {
		log.Printf("Warning: failed to remove existing pin at %s: %v", pinPath, err)
	}
	if err := bpfMap.Pin(pinPath); err != nil {
		return fmt.Errorf("failed to pin map '%s' to %s: %w", mapName, pinPath, err)
	}
	log.Printf("Map '%s' pinned to %s", mapName, pinPath)
	return nil
}


func (bm *BPFManager) Close() {
	for _, l := range bm.links {
		if err := l.Close(); err != nil {
			log.Printf("Error closing BPF link: %v", err)
		}
	}
	if bm.reader != nil {
		if err := bm.reader.Close(); err != nil {
			log.Printf("Error closing ring buffer reader: %v", err)
		}
	}
	if bm.collection != nil {
		if err := bm.collection.Close(); err != nil {
			log.Printf("Error closing BPF collection: %v", err)
		}
	}
	log.Println("All eBPF resources released.")
}