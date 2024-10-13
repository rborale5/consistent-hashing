package consistent

import (
	"encoding/binary"
	"errors"
	"fmt"
	"log"
	"sort"
	"sync"

	hashfunc "github.com/minio/blake2b-simd"
)

type Host struct {
	Name string
	Load int64
}

type Consistent struct {
	hosts         map[uint64]string
	sortedSet     []uint64
	loadMap       map[string]*Host
	replicaFactor uint16
	totalLoad     int64

	sync.RWMutex
}

// func initFunc() {
// 	router := mux.NewRouter()
// 	router.Path("/prometheus").Handler(promhttp.Handler())
// 	err := http.ListenAndServe(":9000", router)
// 	log.Fatal(err)
// }

func New() *Consistent {
	// initFunc()
	return &Consistent{
		hosts:         map[uint64]string{},
		sortedSet:     []uint64{},
		loadMap:       map[string]*Host{},
		replicaFactor: 20,
	}
}

func (c *Consistent) Add(host string) {
	c.Lock()
	defer c.Unlock()

	c.loadMap[host] = &Host{Name: host, Load: 0}
	for i := 0; i < int(c.replicaFactor); i++ {
		h := c.hash(fmt.Sprintf("%s%i", host, i))
		c.hosts[h] = host
		c.sortedSet = append(c.sortedSet, h)
		log.Printf("%d", h)
	}
	sort.Slice(c.sortedSet, func(i, j int) bool {
		if c.sortedSet[i] < c.sortedSet[j] {
			return true
		}
		return false
	})
}

func (c *Consistent) Get(key string) (string, error) {
	c.RLock()
	defer c.RUnlock()

	if len(c.hosts) == 0 {
		return "", errors.New("no hosts added")
	}
	h := c.hash(key)
	idx := c.search(h)

	return c.hosts[c.sortedSet[idx]], nil
}

func (c *Consistent) search(key uint64) int {
	idx := sort.Search(len(c.hosts), func(i int) bool {
		return c.sortedSet[i] >= key
	})
	if idx >= len(c.sortedSet) {
		idx = 0
	}
	return idx
}

func (c *Consistent) hash(key string) uint64 {
	out := hashfunc.Sum512([]byte(key))
	return binary.LittleEndian.Uint64(out[:])
}
