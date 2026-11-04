package runner

import (
	"fmt"
	"sync"
)

// MaxLogBytes bounds retained output from each subprocess. A truncation marker
// is additional small metadata; stdout/stderr continue draining after the cap.
const MaxLogBytes = 64 << 10

type boundedLog struct {
	mu         sync.Mutex
	head, tail []byte
	total      int64
}

func (b *boundedLog) Write(p []byte) (int, error) {
	b.mu.Lock()
	defer b.mu.Unlock()
	n := len(p)
	b.total += int64(n)
	half := MaxLogBytes / 2
	if room := half - len(b.head); room > 0 {
		used := min(room, len(p))
		b.head = append(b.head, p[:used]...)
		p = p[used:]
	}
	if len(p) >= half {
		b.tail = append(b.tail[:0], p[len(p)-half:]...)
	} else {
		overflow := len(b.tail) + len(p) - half
		if overflow > 0 {
			copy(b.tail, b.tail[overflow:])
			b.tail = b.tail[:len(b.tail)-overflow]
		}
		b.tail = append(b.tail, p...)
	}
	return n, nil
}

func (b *boundedLog) snapshot() (string, bool) {
	b.mu.Lock()
	defer b.mu.Unlock()
	truncated := b.total > MaxLogBytes
	middle := ""
	if truncated {
		middle = fmt.Sprintf("\n[MoBlend logs truncated: %d bytes omitted]\n", b.total-MaxLogBytes)
	}
	return string(b.head) + middle + string(b.tail), truncated
}
