package internal

type UserId int64

type User struct {
	Id       UserId
	Name     string
	Email    string
	Password string
	Role     Role
}

type Role int

const (
	RoleUser Role = iota
	RoleAdmin
)

type UserStorage interface {
	// Create создает статью
	Create(name string, email string, password string, role Role) (*User, error)
	// Update обновляет статью
	Update(id UserId, name string, email string) (*User, error)
	// Get получает статью
	Get(id UserId) (*User, error)
	// Delete удаляет статьи
	Delete(id UserId) error
}
