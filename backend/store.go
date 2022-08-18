package main

import (
	"errors"
	"fmt"
	"github.com/deta/deta-go/deta"
	"github.com/deta/deta-go/service/base"
)

var (
	ErrNilValue = errors.New("nil value")
)

// Code xx
type Code struct{
	Key string `json:"key"`
	Content string `json:"content"`
}

// CodeStore xx
type CodeStore struct {
	base *base.Base
}

// NewCodeStore xx
func NewCodeStore(name string, projectKey *string) (*CodeStore, error) {
	var d *deta.Deta
	var err error

	if projectKey != nil {
		d, err = deta.New(deta.WithProjectKey(*projectKey))

	} else {
		d, err = deta.New()
	}
	if err != nil {
		return nil, fmt.Errorf("failed to create new deta instance: %w", err)
	}
	b, err := base.New(d, name)
	if err != nil {
		return nil, fmt.Errorf("failed to create new base instance: %w", err)
	}
	return &CodeStore{
		base: b,
	}, nil
}

// Get xx
func (s *CodeStore) Get(key string) (*Code, error) {
	var c Code
	if err := s.base.Get(key, &c); err != nil {
		if errors.Is(err, deta.ErrNotFound) {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to get code with key %s: %w", key, err)
	}
	return &c, nil
}

// Put xx
func (s *CodeStore) Put(c *Code) error {
	if c == nil {
		return ErrNilValue
	}
	if _, err := s.base.Put(c); err != nil {
		return fmt.Errorf("failed to put code with key %s: %w", c.Key, err)
	}
	return nil
}

// Delete xx
func (s *CodeStore) Delete(key string) error {
	if err := s.base.Delete(key); err != nil {
		return fmt.Errorf("failed to delete code with key %s: %w", key, err)
	}
	return nil
}

// List xx
func (s *CodeStore) List() ([]*Code, error) {
	results := make([]*Code, 0)
	page := make([]*Code, 0)
	lastKey := ""

	i := &base.FetchInput{
		Q:     nil,
		Dest:  &page,
	}
	lastKey, err := s.base.Fetch(i)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch codes: %w", err)
	}
	results = append(results, page...)

	// get all pages
	for lastKey != "" {
		i.LastKey = lastKey
		lastKey, err = s.base.Fetch(i)
		if err != nil {
			return nil, fmt.Errorf("failed to fetch codes: %w", err)
		}
		results = append(results, page...)
	}
	return results, nil
}