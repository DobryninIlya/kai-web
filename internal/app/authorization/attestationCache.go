package authorization

import (
	"main/internal/app/model"
	"sync"
)

type AttestationCache struct {
	mu sync.Mutex
	m  map[string][]model.Discipline
}

func (r *AttestationCache) SetAttestationList(key string, list []model.Discipline) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.m[key] = list
}

func (r *AttestationCache) GetAttestationList(key string) ([]model.Discipline, bool) {
	r.mu.Lock()
	defer r.mu.Unlock()
	value, ok := r.m[key]
	return value, ok
}

func NewAttestationCache() *AttestationCache {
	return &AttestationCache{
		m: make(map[string][]model.Discipline),
	}
}
