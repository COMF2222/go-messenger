package session

type RedisSessionManager struct{}

func NewRedisSessionManager() *RedisSessionManager {
	return &RedisSessionManager{}
}

func (r *RedisSessionManager) SetOnline(userID int) error {
	return SetUserOnline(userID)
}

func (r *RedisSessionManager) SetOffline(userID int) error {
	return SetUserOffline(userID)
}

func (r *RedisSessionManager) IsOnline(userID int) (bool, error) {
	return IsUserOnline(userID)
}

func (r *RedisSessionManager) GetAllOnline() ([]int, error) {
	return GetAllOnlineUsers()
}
