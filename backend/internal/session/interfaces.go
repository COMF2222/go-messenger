package session

type SessionManager interface {
	SetOnline(userID int) error
	SetOffline(userID int) error
	IsOnline(userID int) (bool, error)
	GetAllOnline() ([]int, error)
}
