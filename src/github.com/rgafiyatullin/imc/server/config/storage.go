package config

type StorageConfig interface {
	ResetToDefaults()
	RingSize() uint
}

type storageConfig struct {
	ringSize uint
}

func (this *storageConfig) RingSize() uint { return this.ringSize }
func (this *storageConfig) ResetToDefaults() {
	this.ringSize = 32
}
