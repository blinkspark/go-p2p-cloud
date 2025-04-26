package main

import (
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"

	"github.com/klauspost/reedsolomon"
)

func playReedSolomon() {
	dataShards := 4
	parityShards := 2
	totalShards := dataShards + parityShards
	enc, err := reedsolomon.NewStream(dataShards, parityShards)
	if err != nil {
		log.Panic(err)
	}

	iFileName := "LICENSE"
	iFile, err := os.Open(iFileName)
	if err != nil {
		log.Panic(err)
	}
	defer iFile.Close()

	odir := "."
	oFiles := make([]*os.File, totalShards)
	oWriters := make([]io.Writer, dataShards)
	for i := range totalShards {
		oFileName := fmt.Sprintf("%s.%d", iFileName, i)
		oPath := filepath.Join(odir, oFileName)
		var oFile *os.File
		oFile, err = os.Create(oPath)
		if err != nil {
			log.Panic(err)
		}
		oFiles[i] = oFile
	}
	for i := range dataShards {
		oWriters[i] = oFiles[i]
	}
	defer func() {
		for _, oFile := range oFiles {
			if oFile != nil {
				oFile.Close()
			}
		}
	}()
	iFileState, err := iFile.Stat()
	if err != nil {
		log.Panic(err)
	}
	err = enc.Split(iFile, oWriters, iFileState.Size())
	if err != nil {
		log.Panic(err)
	}

	oReaders := make([]io.Reader, dataShards)
	for i := range dataShards {
		oFiles[i].Close()
		ofname := filepath.Join(odir, fmt.Sprintf("%s.%d", iFileName, i))
		oFiles[i], err = os.Open(ofname)
		if err != nil {
			log.Panic(err)
		}
		oReaders[i] = oFiles[i]
	}
	parityWriters := make([]io.Writer, parityShards)
	for i := range parityShards {
		parityWriters[i] = oFiles[dataShards+i]
	}
	err = enc.Encode(oReaders, parityWriters)
	if err != nil {
		log.Panic(err)
	}
}
