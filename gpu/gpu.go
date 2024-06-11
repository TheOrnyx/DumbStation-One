package gpu

type Gpu struct {
	gpuStat GpuStat
}

// NewGPU create and return a new gpu
func NewGPU() Gpu {
	return Gpu{gpuStat: NewGPUStat()}
}

