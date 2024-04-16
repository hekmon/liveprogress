package main

import (
	"crypto/sha256"
	"hash"
	"io"
	"sync"
)

type SHA256Progress struct {
	source    *readerCounter
	hasher    *hasherReporter
	hash      []byte
	computing sync.Mutex
}

func New(reader io.Reader, totalBytes int) *SHA256Progress {
	if reader == nil || totalBytes <= 0 {
		return nil
	}
	return &SHA256Progress{
		source: &readerCounter{
			source:  reader,
			maxRead: totalBytes,
		},
	}
}

func (hp *SHA256Progress) ComputeHash(reportWritten func(uint64)) (err error) {
	defer hp.computing.Unlock()
	hp.computing.Lock()
	if hp.hash != nil {
		return
	}
	// Prepare the hasher
	hp.hasher = &hasherReporter{
		dest:        sha256.New(),
		writeReport: reportWritten,
	}
	// Start copy
	if _, err = io.Copy(hp.hasher, hp.source); err != nil {
		return
	}
	hp.hash = hp.hasher.GetCurrentHash()
	hp.hasher = nil
	return
}

func (hp *SHA256Progress) GetCurrentHash() []byte {
	if hp.hash != nil {
		return hp.hash
	}
	return hp.hasher.GetCurrentHash()
}

type readerCounter struct {
	source  io.Reader
	read    int
	maxRead int
}

func (rc *readerCounter) Read(p []byte) (n int, err error) {
	if rc.read >= rc.maxRead {
		err = io.EOF
		return
	}
	n, err = rc.source.Read(p)
	rc.read += n
	return
}

type hasherReporter struct {
	dest        hash.Hash
	writeReport func(uint64)
	access      sync.Mutex
}

func (hc *hasherReporter) Write(p []byte) (n int, err error) {
	hc.access.Lock()
	n, err = hc.dest.Write(p)
	hc.writeReport(uint64(n))
	hc.access.Unlock()
	return
}

func (hc *hasherReporter) GetCurrentHash() (hash []byte) {
	hc.access.Lock()
	hash = hc.dest.Sum(nil)
	hc.access.Unlock()
	return
}
