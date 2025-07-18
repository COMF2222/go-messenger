package session

import (
	"fmt"
	"time"
)

const onlineKeyPrefix = "online_user"

// Установить статус онлайн
func SetUserOnline(userID int) error {
	key := fmt.Sprintf("%s:%d", onlineKeyPrefix, userID)
	return Rdb.Set(Ctx, key, true, 5*time.Minute).Err()
}

// Установить статус оффлайн
func SetUserOffline(userID int) error {
	key := fmt.Sprintf("%s:%d", onlineKeyPrefix, userID)
	return Rdb.Del(Ctx, key).Err()
}

// Проверить, онлайн ли пользователь
func IsUserOnline(userID int) (bool, error) {
	key := fmt.Sprintf("%s:%d", onlineKeyPrefix, userID)
	exists, err := Rdb.Exists(Ctx, key).Result()
	if err != nil {
		return false, err
	}
	return exists == 1, nil
}

// Получить всех онлайн-пользователей
func GetAllOnlineUsers() ([]int, error) {
	pattern := fmt.Sprintf("%s:*", onlineKeyPrefix)
	keys, err := Rdb.Keys(Ctx, pattern).Result()
	if err != nil {
		return nil, err
	}

	var userIDs []int
	for _, key := range keys {
		var id int
		if _, err := fmt.Sscanf(key, onlineKeyPrefix+":%d", &id); err == nil {
			userIDs = append(userIDs, id)
		}
	}

	return userIDs, nil
}
