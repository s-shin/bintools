package main

import (
	"bytes"
	"io"
)

type findIndexOptions struct {
	bufferSize      int
	convolutionSize int
}

type FindIndexOptionValue interface {
	setFindIndexOptionValue(opts *findIndexOptions)
}

type BufferSize int

func (s BufferSize) setFindIndexOptionValue(opts *findIndexOptions) {
	opts.bufferSize = int(s)
}

type ConvolutionSize int

func (s ConvolutionSize) setFindIndexOptionValue(opts *findIndexOptions) {
	opts.convolutionSize = int(s)
}

func FindIndexByKey(target io.Reader, key []byte, vals ...FindIndexOptionValue) int {
	vals = append(vals, ConvolutionSize(len(key)))
	return FindIndex(target, func(bs []byte) bool {
		return bytes.Compare(key, bs) == 0
	}, vals...)
}

func FindIndex(target io.Reader, matcher func(bs []byte) bool, vals ...FindIndexOptionValue) int {
	opts := &findIndexOptions{
		bufferSize:      0,
		convolutionSize: 0,
	}
	for _, v := range vals {
		v.setFindIndexOptionValue(opts)
	}

	cnvSize := opts.convolutionSize
	bufSize := opts.bufferSize
	if bufSize <= 0 {
		bufSize = cnvSize * 2
	}
	if bufSize <= cnvSize {
		return -1
	}

	prevBufSize := cnvSize - 1
	nextBufSize := bufSize - prevBufSize

	activeBuf := make([]byte, bufSize)
	nextBuf := make([]byte, nextBufSize)
	tmpBuf := make([]byte, bufSize)

	var numRead int
	var err error
	numRead, err = target.Read(activeBuf)
	activeBuf = activeBuf[:numRead]
	if numRead < cnvSize {
		return -1
	}
	numRead -= (cnvSize - 1)
	var end = err != nil

	var offset int
	for {
		for i := 0; i < numRead; i++ {
			if matcher(activeBuf[i : i+cnvSize]) {
				return offset + i
			}
		}
		if end {
			break
		}
		numRead, err = target.Read(nextBuf)
		end = err != nil
		copy(tmpBuf, activeBuf[nextBufSize:])
		copy(tmpBuf[prevBufSize:], nextBuf[:numRead])
		tmpBuf = tmpBuf[:prevBufSize+numRead]
		t := tmpBuf
		tmpBuf = activeBuf
		activeBuf = t
		offset += nextBufSize
	}

	return -1
}
