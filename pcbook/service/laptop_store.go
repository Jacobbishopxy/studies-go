package service

import (
	"errors"
	"fmt"
	"pcbook/pb"
	"sync"

	"github.com/jinzhu/copier"
)

// laptop store 接口
type LaptopStore interface {
	Save(laptop *pb.Laptop) error

	Find(id string) (*pb.Laptop, error)
}

var _ LaptopStore = &InMemoryLaptopStore{}

// 内存数据库
type InMemoryLaptopStore struct {
	mutex sync.RWMutex
	data  map[string]*pb.Laptop
}

func NewInMemoryLaptopStore() *InMemoryLaptopStore {
	return &InMemoryLaptopStore{
		data: make(map[string]*pb.Laptop),
	}
}

var ErrAlreadyExists = errors.New("record already exists")

func (store *InMemoryLaptopStore) Save(laptop *pb.Laptop) error {
	store.mutex.Lock()
	defer store.mutex.Unlock()

	if store.data[laptop.Id] != nil {
		return ErrAlreadyExists
	}

	orther, err := deepCopy(laptop)
	if err != nil {
		return err
	}

	store.data[orther.Id] = orther
	return nil
}

func deepCopy(laptop *pb.Laptop) (*pb.Laptop, error) {
	other := &pb.Laptop{}

	err := copier.Copy(other, laptop)
	if err != nil {
		return nil, fmt.Errorf("cannot copy laptop data: %w", err)
	}

	return other, nil
}

func (store *InMemoryLaptopStore) Find(id string) (*pb.Laptop, error) {
	store.mutex.RLock()
	defer store.mutex.RUnlock()

	laptop := store.data[id]
	if laptop == nil {
		return nil, nil
	}

	return deepCopy(laptop)
}
