package repository

import (
	"fmt"
	"sync"
)

const MAX_USERS = 100

// Currently logged users container
type UserRepository struct {
	m sync.RWMutex
	users map[string]string
}

// Init map data structure
var repo = UserRepository{
	users: make(map[string]string, MAX_USERS),
}

// Registery user as surrently lo
func UserLogin(id string) error {
	repo.m.Lock()
	defer repo.m.Unlock()
	
	repo.users[id] = id
	
	return nil
}

// Check if the user is already loggerd
func UserExist(id string) bool {
	repo.m.Lock()
	defer repo.m.Unlock()
	
	_, exists := repo.users[id]
	
	return exists
}

// Deregister user from the currently logged list
func UserLogoff(id string) error {
	repo.m.Lock()
	defer repo.m.Unlock()
	
	_, exists := repo.users[id]
	if !exists  {
		return fmt.Errorf("User not registered: %s", id)
	}
	
	delete(repo.users, id)
	
	return nil
}
