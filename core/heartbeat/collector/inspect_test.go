package collector

import (
	"bufio"
	"bytes"
	"testing"
)

func TestGetCPUInfo(t *testing.T) {

	r := bytes.NewBufferString(`processor       : 0
BogoMIPS        : 48.00
Features        : fp asimd evtstrm aes pmull sha1 sha2 crc32
CPU implementer : 0x41
CPU architecture: 8
CPU variant     : 0x0
CPU part        : 0xd03
CPU revision    : 4

processor       : 1
BogoMIPS        : 48.00
Features        : fp asimd evtstrm aes pmull sha1 sha2 crc32
CPU implementer : 0x41
CPU architecture: 8
CPU variant     : 0x0
CPU part        : 0xd03
CPU revision    : 4

processor       : 2
BogoMIPS        : 48.00
Features        : fp asimd evtstrm aes pmull sha1 sha2 crc32
CPU implementer : 0x41
CPU architecture: 8
CPU variant     : 0x0
CPU part        : 0xd03
CPU revision    : 4

processor       : 3
BogoMIPS        : 48.00
Features        : fp asimd evtstrm aes pmull sha1 sha2 crc32
CPU implementer : 0x41
CPU architecture: 8
CPU variant     : 0x0
CPU part        : 0xd03
CPU revision    : 4

processor       : 4
BogoMIPS        : 48.00
Features        : fp asimd evtstrm aes pmull sha1 sha2 crc32
CPU implementer : 0x41
CPU architecture: 8
CPU variant     : 0x0
CPU part        : 0xd08
CPU revision    : 2

processor       : 5
BogoMIPS        : 48.00
Features        : fp asimd evtstrm aes pmull sha1 sha2 crc32
CPU implementer : 0x41
CPU architecture: 8
CPU variant     : 0x0
CPU part        : 0xd08
CPU revision    : 2

Serial          : 01d656a0002a3ebf`)
	scanner := bufio.NewScanner(r)

	info := cpuInfoFromScanner(scanner)
	infos := info.Processors

	if len(infos) != 6 {
		t.Error(info)
	}

	t.Log(info.Serial)
	for _, it := range infos {
		t.Log(it)
	}
}

func TestMemoryInfoFromScanner(t *testing.T) {

	r := bytes.NewBufferString(`MemTotal:        3934472 kB
MemFree:          411592 kB
MemAvailable:    2804084 kB
Buffers:          131276 kB
Cached:          2218008 kB
SwapCached:            0 kB
Active:          2180672 kB
Inactive:        1064416 kB
Active(anon):     896900 kB
Inactive(anon):    15712 kB
Active(file):    1283772 kB
Inactive(file):  1048704 kB
Unevictable:           0 kB
Mlocked:               0 kB
SwapTotal:             0 kB
SwapFree:              0 kB
Dirty:              1216 kB
Writeback:            16 kB
AnonPages:        895852 kB
Mapped:           439128 kB
Shmem:             16792 kB
Slab:             150208 kB
SReclaimable:      91736 kB
SUnreclaim:        58472 kB
KernelStack:       31216 kB
PageTables:         9428 kB
NFS_Unstable:          0 kB
Bounce:                0 kB
WritebackTmp:          0 kB
CommitLimit:     1967236 kB
Committed_AS:    9465808 kB
VmallocTotal:   258867136 kB
VmallocUsed:           0 kB
VmallocChunk:          0 kB`)
	scanner := bufio.NewScanner(r)

	info := memoryInfoFromScanner(scanner)

	t.Log(info)

}
