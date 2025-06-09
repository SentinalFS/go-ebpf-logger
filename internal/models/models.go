package models

import (
	"github.com/cilium/ebpf"
	"github.com/cilium/ebpf/link"
	"github.com/cilium/ebpf/ringbuf"
)

type Data struct {
	Pid       uint32        
	Uid       uint32        
	Filename  [128]byte     
	Comm      [16]byte      
	Timestamp uint64        
	Otype     [16]byte      
}


type BPFManager struct {
	Collection *ebpf.Collection      
	Reader     *ringbuf.Reader       
	Links      []link.Link           
}