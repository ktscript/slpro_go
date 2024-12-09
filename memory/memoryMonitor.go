package memory

import (
	"log"
	"runtime"
	"time"
)

func MonitorMemoryUsage() {
	for {
		var m runtime.MemStats
		runtime.ReadMemStats(&m)

		if m.Alloc > 1000*1024*1024 { // 1GB test threshold
			log.Printf("High memory usage detected: %v bytes, triggering garbage collection", m.Alloc)
			runtime.GC()
		}

		time.Sleep(30 * time.Second)
	}
}
 