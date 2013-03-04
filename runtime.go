package gmx

// pkg/runtime instrumentation

import "runtime"

var memstats runtime.MemStats

func init() {
	reg := Registry("runtime")
	reg("gomaxprocs", runtimeGOMAXPROCS)
	reg("numcgocall", runtimeNumCgoCall)
	reg("numcpu", runtimeNumCPU)
	reg("numgoroutine", runtimeNumGoroutine)
	reg("version", runtimeVersion)
	reg("memstats", runtimeMemStats)
}

func runtimeGOMAXPROCS() interface{} {
	return runtime.GOMAXPROCS(0)
}

func runtimeNumCgoCall() interface{} {
	return runtime.NumCgoCall()
}

func runtimeNumCPU() interface{} {
	return runtime.NumCPU()
}

func runtimeNumGoroutine() interface{} {
	return runtime.NumGoroutine()
}

func runtimeVersion() interface{} {
	return runtime.Version()
}

func runtimeMemStats() interface{} {
	runtime.ReadMemStats(&memstats)
	return memstats
}
