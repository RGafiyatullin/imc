package config

import (
	"fmt"
	"os"
	"strconv"
)

const StorageRingSizeDefault = 32

type StorageConfig interface {
	ResetToDefaults()
	RingSize() uint
}

type storageConfig struct {
	ringSize uint
}

func (this *storageConfig) RingSize() uint { return this.ringSize }
func (this *storageConfig) ResetToDefaults() {
	this.ringSize = StorageRingSizeDefault
}

func (this *storageConfig) ReadFromOSEnv() {
	ringSizeStr := os.Getenv("IMCD_STORAGE_RING_SIZE")
	if ringSizeStr != "" {
		ringSizeParsed, ringSizeParseErr := strconv.ParseInt(ringSizeStr, 10, 64)
		if ringSizeParseErr != nil || ringSizeParsed > 256 || ringSizeParsed < 1 {
			os.Stderr.WriteString(
				fmt.Sprintf(
					"Invalid value for IMCD_STORAGE_RING_SIZE: should be an integer [1...256]; using default value: %d",
					StorageRingSizeDefault))
		} else {
			this.ringSize = uint(ringSizeParsed)
		}
	}

}
