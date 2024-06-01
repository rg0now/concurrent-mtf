package main

import (
	"encoding/binary"
	"fmt"
	"hash"
	"math"
	"math/rand"
	"strings"

	"github.com/spaolacci/murmur3"
)

var bs = make([]byte, 8)

type hasher struct {
	hash hash.Hash64
	id   int
	next *hasher
}

// BloomFilter represents a single Bloom filter structure.
type BloomFilter struct {
	bitSet []bool   // The bit array represented as a slice of bool
	n      int      // The expected number of elements to store
	m      int      // The number of bits in the bit set (shortcut for len(bitSet)
	k      int      // The number of hash functions to use (shortcut for len(hashes)
	hashes []hasher // The hash functions to use
	front  *hasher
	rnd    *rand.Rand // For randomizing lookups
	sa     bool       // Self-adjusting?
}

// NewBloomFilterWithHasher creates a new Bloom filter with the given number of elements (n) and
// false positive rate (p).
func NewBloomFilter(n int, p float64) *BloomFilter {
	return newBloomFilter(n, p, false)
}

// NewBloomSelfAdjustingFilter creates a new self-adjusting Bloom filter with the given number of elements (n) and
// false positive rate (p).
func NewSelfAdjustingBloomFilter(n int, p float64) *BloomFilter {
	return newBloomFilter(n, p, true)
}

// NewBloomFilterWithHasher creates a new Bloom filter with the given number of elements (n) and
// false positive rate (p).
func newBloomFilter(n int, p float64, sa bool) *BloomFilter {
	if n == 0 {
		panic("number of elements cannot be 0")
	}
	if p <= 0 || p >= 1 {
		panic("false positive rate must be between 0 and 1")
	}
	m, k := getOptimalParams(n, p)
	hashes := make([]hasher, k)
	for i := 0; i < k; i++ {
		hashes[i].id = i
		hashes[i].hash = murmur3.New64WithSeed(uint32(i)) // use id as hash seed
		if i < k-1 {
			hashes[i].next = &hashes[i+1]
		}
	}

	return &BloomFilter{
		n:      n,
		m:      m,
		k:      k,
		bitSet: make([]bool, m),
		hashes: hashes,
		front:  &hashes[0],
		rnd:    rand.New(rand.NewSource(seed)),
		sa:     sa,
	}
}

// getOptimalParams calculates the optimal parameters for the Bloom filter,
// the number of bits in the bit set (m) and the number of hash functions (k).
func getOptimalParams(n int, p float64) (int, int) {
	m := int(math.Ceil(-1 * float64(n) * math.Log(p) / math.Pow(math.Log(2), 2)))
	if m == 0 {
		m = 1
	}
	k := int(math.Ceil((float64(m) / float64(n)) * math.Log(2)))
	if k == 0 {
		k = 1
	}
	return m, k
}

func (bf *BloomFilter) Add(i Item) {
	binary.LittleEndian.PutUint64(bs, uint64(i.Id()))
	bf.add(bs)
}

// Add adds an item to the Bloom filter.
func (bf *BloomFilter) add(data []byte) {
	for _, hash := range bf.hashes {
		hash.hash.Reset()
		hash.hash.Write(data)
		hashValue := hash.hash.Sum64() % uint64(bf.m)
		bf.bitSet[hashValue] = true
		// fmt.Printf("setting bit %d at hash %d\n", hashValue, hash.id)
	}
}

func (bf *BloomFilter) Find(i Item) int {
	// convert to a []byte
	binary.LittleEndian.PutUint64(bs, uint64(i.Id()))

	// test
	if bf.test(bs) {
		return 1
	}
	return -1
}

func (bf *BloomFilter) test(data []byte) bool {
	for current, prev := bf.front, (*hasher)(nil); current != nil; current, prev = current.next, current {
		hash := current.hash
		hash.Reset()
		hash.Write(data)
		hashValue := hash.Sum64() % uint64(bf.m)
		if !bf.bitSet[hashValue] {
			// fmt.Printf("false with value %d at hash %d\n", hashValue, current.id)
			if bf.sa && current != bf.front {
				// move-to-front
				// remove current from the list
				if prev != nil {
					prev.next = current.next
				}

				// insert to the front
				current.next = bf.front
				bf.front = current
			}
			return false
		}
	}

	return true
}

// Print displays all the nodes from linked list
func (bf *BloomFilter) String() string {
	t := "Bloomfilter"
	if bf.sa {
		t = "SA_Bloomfilter"
	}
	ret := fmt.Sprintf("%s: n(#elems): %d, m(width) %d, k(#hash): %d, hash-order: ", t, bf.n, bf.m, bf.k)
	hs := []string{}
	for current := bf.front; current != nil; current = current.next {
		hs = append(hs, fmt.Sprintf("%d", current.id))
	}
	return ret + strings.Join(hs, ", ")
}
