package config

type storageConfig struct {
	ringSize uint
}

func (this *storageConfig) RingSize() uint { return this.ringSize }
func (this *storageConfig) ResetToDefaults() {
	this.ringSize = 32
}
