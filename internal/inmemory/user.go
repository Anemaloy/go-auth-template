package inmemory

import (
	"auth/internal"
	"errors"
	"sync"
)

type UserStorage struct {
	mu     sync.RWMutex
	users  map[internal.UserId]*internal.User
	lastId internal.UserId
}

func NewUserStorage() *UserStorage {
	return &UserStorage{
		users: make(map[internal.UserId]*internal.User),
	}
}

func (a *UserStorage) Create(name string, email string, password string, role internal.Role) (*internal.User, error) {
	a.mu.Lock()
	defer a.mu.Unlock()

	user := &internal.User{
		Id:       a.lastId + 1,
		Name:     name,
		Email:    email,
		Password: password,
		Role:     role,
	}

	a.users[user.Id] = user
	a.lastId = user.Id

	return user, nil
}

func (a *UserStorage) Update(id internal.UserId, name string, email string) (*internal.User, error) {
	a.mu.Lock()
	defer a.mu.Unlock()

	user, ok := a.users[id]
	if !ok {
		return nil, errors.New("user doesn't exist")
	}

	user.Name = name
	user.Email = email

	return user, nil
}

func (a *UserStorage) Get(id internal.UserId) (*internal.User, error) {
	a.mu.RLock()
	defer a.mu.RUnlock()
	if _, ok := a.users[id]; ok {
		return a.users[id], nil
	}

	return nil, errors.New("user not found")
}

func (a *UserStorage) Delete(id internal.UserId) error {
	a.mu.Lock()
	defer a.mu.Unlock()

	delete(a.users, id)

	return nil
}
