package rsem

import (
	"context"
	"fmt"
	"time"

	"github.com/go-redis/redis/v8"
)

type Semaphore struct {
	rdb redis.UniversalClient
}

func NewSemaphore(rdb redis.UniversalClient) *Semaphore {
	return &Semaphore{rdb: rdb}
}

func (s *Semaphore) Acquire(ctx context.Context, key string, limit int) (func() error, error) {
	uniqueID := fmt.Sprintf("%d-%d", time.Now().UnixNano(), time.Now().UnixNano()%1000)
	lockKey := fmt.Sprintf("semaphore:%s", key)

	// Lua-скрипт для захвата семафора
	script := redis.NewScript(`
		-- Удаляем просроченные записи
		redis.call("zremrangebyscore", KEYS[1], "-inf", ARGV[1])

		-- Проверяем, достигнут ли лимит
		if redis.call("zcard", KEYS[1]) < tonumber(ARGV[2]) then
			redis.call("zadd", KEYS[1], ARGV[3], ARGV[4])
			redis.call("pexpire", KEYS[1], ARGV[5])
			return 1
		else
			return 0
		end
	`)

	now := time.Now().UnixNano()
	expiration := int64(5000) // TTL ключа в миллисекундах
	score := now

	for {
		// Пытаемся выполнить Lua-скрипт
		res, err := script.Run(ctx, s.rdb, []string{lockKey}, now, limit, score, uniqueID, expiration).Int()
		if err != nil {
			return nil, err
		}

		if res == 1 {
			// Успешно захватили семафор
			release := func() error {
				_, err := s.rdb.ZRem(ctx, lockKey, uniqueID).Result()
				return err
			}
			return release, nil
		}

		// Если лимит достигнут, ждем и повторяем попытку
		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		case <-time.After(50 * time.Millisecond):
			// Ждем перед повторной попыткой
		}
	}
}
