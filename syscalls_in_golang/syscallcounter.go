package main

import (
	"fmt"
	"sort"

	seccomp "github.com/seccomp/libseccomp-golang"
)

// SyscallCounter tracks statistics about syscalls
type SyscallCounter struct {
	counts map[uint64]int
}

// NewSyscallCounter initializes a new SyscallCounter
func NewSyscallCounter() *SyscallCounter {
	return &SyscallCounter{
		counts: make(map[uint64]int),
	}
}

// Inc increments the count for a syscall
func (s *SyscallCounter) Inc(syscallNum uint64) {
	s.counts[syscallNum]++
}

// GetName returns the name of a syscall
func (s *SyscallCounter) GetName(syscallNum uint64) string {
	name, err := seccomp.ScmpSyscall(syscallNum).GetName()
	if err != nil {
		return fmt.Sprintf("unknown(%d)", syscallNum)
	}
	return name
}

// Print displays statistics about syscall usage
func (s *SyscallCounter) Print() {
	fmt.Println("\nSyscall Statistics:")

	// Convert map to slice for sorting
	type syscallCount struct {
		num   uint64
		count int
		name  string
	}

	var counts []syscallCount
	for num, count := range s.counts {
		counts = append(counts, syscallCount{
			num:   num,
			count: count,
			name:  s.GetName(num),
		})
	}

	// Sort by count (most frequent first)
	sort.Slice(counts, func(i, j int) bool {
		return counts[i].count > counts[j].count
	})

	// Print results
	fmt.Println("Count  | Syscall Name")
	fmt.Println("-------|-------------")
	for _, sc := range counts {
		fmt.Printf("%-6d | %s\n", sc.count, sc.name)
	}
}
