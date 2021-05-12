package bpf

import (
	"os"
	"os/signal"

	"github.com/iovisor/gobpf/bcc"
)

type BPF struct {
	module  *bcc.Module
	perfMap *bcc.PerfMap
	signal  chan os.Signal
	channel chan []byte
}

func New(source string) BPF {
	return BPF{
		module:  bcc.NewModule(source, []string{}),
		channel: make(chan []byte),
	}
}

func (b *BPF) Close() {
	b.module.Close()
}

func (b *BPF) AttachKprobe(event, fn string) error {
	fName, err := b.module.LoadKprobe(fn)
	if err != nil {
		return err
	}
	return b.module.AttachKprobe(
		bcc.GetSyscallFnName(event),
		fName, -1,
	)
}

func (b *BPF) NewEventListener(event string) error {
	var err error
	b.signal = make(chan os.Signal, 1)
	b.perfMap, err = bcc.InitPerfMap(
		bcc.NewTable(b.module.TableId(event), b.module),
		b.channel, nil,
	)
	if err != nil {
		return err
	}

	signal.Notify(b.signal, os.Interrupt, os.Kill)
	return err
}

func (b *BPF) Run(fn func(data <-chan []byte)) {
	go fn(b.channel)

	b.perfMap.Start()
	<-b.signal
	b.perfMap.Stop()
}
