package eventstream

import (
	"log"
	"sync"
)

type mutex struct {
	m sync.RWMutex
	i int
}

func (m *mutex) Lock() {
	log.Println("Lock    ")
	m.m.Lock()
	m.i++
	log.Println("Lock    ", m.i)
}

func (m *mutex) Unlock() {
	log.Println("Unlock  ", m.i)
	m.m.Unlock()
}

func (m *mutex) RLock() {
	log.Println("RLock   ")
	m.m.RLock()
	m.i++
	log.Println("RLock   ", m.i)
}

func (m *mutex) RUnlock() {
	log.Println("RUnlock ", m.i)
	m.m.RUnlock()
}
