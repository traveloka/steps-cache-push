package main

import (
	"compress/gzip"
	"fmt"
	"io"
	"os"
	"time"

	"github.com/bitrise-io/go-utils/log"
	"github.com/pierrec/lz4/v4"
	"github.com/klauspost/compress/zstd"
)


const (
	maxConcurrency	  = -1
)

type CompressionWriter struct {
	writer		io.Writer
	closer		io.Closer
}


func FastArchiveCompress(cacheArchivePath, compressor string) (int64, error) {
	var compressedArchiveSize int64
	compressStartTime := time.Now()

	in, err := os.Open(cacheArchivePath)
	if err != nil {
		return 0, fmt.Errorf("Fatal error in opening file: ", err.Error())
	}
	defer in.Close()

	compressionWriter, outputFile, err := NewCompressionWriter(cacheArchivePath, compressor, maxConcurrency)
	if err != nil {
		return 0, fmt.Errorf("Error getting compressor writer: ", err.Error())
	}

	_, err = io.Copy(compressionWriter.writer, in)
	if err != nil {
		return 0, fmt.Errorf("Error compressing file:", err.Error())
	}

	defer compressionWriter.closer.Close()

	fileInfo, err := outputFile.Stat()
	if err == nil {
		compressedArchiveSize = fileInfo.Size()
	}

	err = os.Remove(cacheArchivePath)
	if err != nil {
		return 0, fmt.Errorf("Error deleting uncompressed archive file: ", err.Error())
	}

	log.Infof("Done compressing file using %s in %s", compressor, time.Since(compressStartTime))

	return compressedArchiveSize, nil
}

func NewCompressionWriter(cacheArchivePath, compressor string, concurrency int) (*CompressionWriter, *os.File, error) {
	if compressor == "lz4" {
		compressedOutputFile := createCompressedOutputFile(ExtendPathWithCompression(cacheArchivePath, compressor))
		lz4Writer := lz4.NewWriter(compressedOutputFile)
		options := []lz4.Option{
			lz4.CompressionLevelOption(lz4.Level5),
			lz4.ConcurrencyOption(concurrency),
		}
		if err := lz4Writer.Apply(options...); err != nil {
			return nil, compressedOutputFile, err
		}

		return &CompressionWriter{
			writer: lz4Writer,
			closer: lz4Writer,
		}, compressedOutputFile, nil
	} else if compressor == "gzip" {
		compressedOutputFile := createCompressedOutputFile(ExtendPathWithCompression(cacheArchivePath, compressor))
		gzipWriter, err := gzip.NewWriterLevel(compressedOutputFile, gzip.BestCompression)
		if err != nil {
			return nil, compressedOutputFile, err
		}

		return  &CompressionWriter{
			writer: gzipWriter,
			closer: gzipWriter,
		}, compressedOutputFile, nil
	} else if compressor == "zstd" {
		compressedOutputFile := createCompressedOutputFile(ExtendPathWithCompression(cacheArchivePath, compressor))
		zstdWriter, err := zstd.NewWriter(compressedOutputFile)
		if err != nil {
			return nil, compressedOutputFile, err
		}
		
		return  &CompressionWriter{
			writer: zstdWriter,
			closer: zstdWriter,
		}, compressedOutputFile, nil
	}
	
	log.Errorf("Unsupported compressor algorithm in fast-archiver for: ", compressor)
	os.Exit(1)

	return nil, nil, nil
}

func createCompressedOutputFile(path string) (*os.File) {
	compressedOutputFile, err := os.Create(path)

	log.Infof("Compressing file into: ", path)

	if err != nil {
		log.Errorf("Error when creating new compressed file", err.Error())
		os.Exit(1)

		return nil
	}

	return compressedOutputFile
}

func ExtendPathWithCompression(path, compressionAlgorithm string) (string) {
	if compressionAlgorithm == "lz4" {
		return path + ".lz4"
	} else if compressionAlgorithm == "gzip" {
		return path + ".gz"
	} else if compressionAlgorithm == "zstd" {
		return path + ".zst"
	}

	return path
}