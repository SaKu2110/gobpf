package main

import (
	"bytes"
	"dev/hello/v1/bpf"
	"encoding/binary"
	"fmt"
	"os"
	"unsafe"
)

import "C"

type events struct {
	Ts   uint64
	Task [16]byte
	Pid  uint32
	Type uint32
	Argv [128]byte
}

func handler(d <-chan []byte) {
	var event events
	var args string = ""

	for {
		// data := <-channel
		data := <-d
		err := binary.Read(bytes.NewBuffer(data), binary.LittleEndian, &event)
		if err != nil {
			fmt.Printf("failed to decode received data: %s\n", err)
			continue
		}

		if event.Type == 1 {
			argv := C.GoString((*C.char)(unsafe.Pointer(&event.Argv)))
			args = args + argv
			args = args + " "
		} else {
			task := C.GoString((*C.char)(unsafe.Pointer(&event.Task)))
			fmt.Printf("%-18d %-16s %-6d %s\n", (event.Ts)/1000000000, task, event.Pid, args)

			args = ""
		}
	}
}

func main() {
	bpf := bpf.New(Source)
	defer bpf.Close()

	if err := bpf.AttachKprobe("execve", "syscall__execve"); err != nil {
		fmt.Fprintf(os.Stderr, "ERROR: %s\n", err)
		os.Exit(1)
	}

	fmt.Println("Tracing strlen()... hit Ctrl-C to end.")
	fmt.Printf("%-18s %-16s %-6s %s\n", "TIME(s)", "COMM", "PID", "MESSAGE")

	if err := bpf.NewEventListener("events"); err != nil {
		fmt.Fprintf(os.Stderr, "ERROR: %s\n", err)
		os.Exit(1)
	}

	bpf.Run(handler)
}
