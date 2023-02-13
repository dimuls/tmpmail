package redis

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/go-redis/redis/v8"

	"tmpmail/entity"
)

type Storage struct {
	redis *redis.Client
}

func NewStorage(addr string) *Storage {
	return &Storage{
		redis: redis.NewClient(&redis.Options{
			Addr: addr,
		}),
	}
}

func (s *Storage) Close() error {
	return s.redis.Close()
}

func tokenKey(token string) string {
	return "tkns/" + token
}

func accountKey(username string) string {
	return "accs/" + username
}

func (s *Storage) CreateAccount(token, username string, ttl time.Duration) error {
	tKey := tokenKey(token)
	exists, err := s.redis.Exists(context.Background(), tKey).Result()
	if err != nil {
		return fmt.Errorf("check token exists: %w", err)
	}
	if exists == 1 {
		return fmt.Errorf("token already exists")
	}
	_, err = s.redis.Set(context.Background(), tKey, username, ttl).Result()
	if err != nil {
		return fmt.Errorf("set token")
	}
	aKey := accountKey(username)
	exists, err = s.redis.Exists(context.Background(), aKey).Result()
	if err != nil {
		return fmt.Errorf("check account exists: %w", err)
	}
	if exists == 1 {
		return fmt.Errorf("account already exists")
	}
	_, err = s.redis.LPush(context.Background(), aKey, "-").Result()
	if err != nil {
		return fmt.Errorf("lpush account: %w", err)
	}
	_, err = s.redis.Expire(context.Background(), aKey, ttl).Result()
	if err != nil {
		return fmt.Errorf("expire account: %w", err)
	}
	return nil
}

func (s *Storage) ProlongAccount(token string, ttl time.Duration) error {
	tKey := tokenKey(token)
	_, err := s.redis.Expire(context.Background(), tKey, ttl).Result()
	if err != nil {
		return fmt.Errorf("expire token: %w", err)
	}
	username, err := s.redis.Get(context.Background(), tKey).Result()
	if err != nil {
		return fmt.Errorf("get token username: %w", err)
	}
	_, err = s.redis.Expire(context.Background(), accountKey(username), ttl).Result()
	if err != nil {
		return fmt.Errorf("expire account: %w", err)
	}
	return err
}

func (s *Storage) Account(token string) (entity.Account, error) {
	tKey := tokenKey(token)
	username, err := s.redis.Get(context.Background(), tKey).Result()
	if err != nil {
		if errors.Is(err, redis.Nil) {
			return entity.Account{}, entity.ErrAccountDoesntExists
		}
		return entity.Account{}, fmt.Errorf("get token username: %w", err)
	}
	ttl, err := s.redis.TTL(context.Background(), tKey).Result()
	if err != nil {
		return entity.Account{}, fmt.Errorf("account ttl: %w", err)
	}
	if ttl == 0 {
		return entity.Account{}, entity.ErrAccountDoesntExists
	}
	aKey := accountKey(username)
	emailLen, err := s.redis.LLen(context.Background(), aKey).Result()
	if err != nil {
		return entity.Account{}, fmt.Errorf("llen account: %w", err)
	}
	if emailLen == 0 {
		return entity.Account{}, entity.ErrAccountDoesntExists
	}
	var emails []entity.Email
	if emailLen > 1 {
		emailJSONs, err := s.redis.LRange(context.Background(), aKey, 0, emailLen-2).Result()
		if err != nil {
			return entity.Account{}, fmt.Errorf("lrange account: %w", err)
		}
		for _, emailJSON := range emailJSONs {
			var email entity.Email
			err := json.Unmarshal([]byte(emailJSON), &email)
			if err != nil {
				return entity.Account{}, fmt.Errorf("json unmarshal email json: %w", err)
			}
			emails = append(emails, email)
		}
	}
	return entity.Account{
		Username: username,
		TTL:      ttl.Milliseconds(),
		Emails:   emails,
	}, nil
}

func (s *Storage) RemoveAccount(token string) error {
	tKey := tokenKey(token)
	username, err := s.redis.Get(context.Background(), tKey).Result()
	if err != nil {
		if errors.Is(err, redis.Nil) {
			return entity.ErrAccountDoesntExists
		}
		return fmt.Errorf("get token username: %w", err)
	}
	aKey := accountKey(username)
	_, err = s.redis.Del(context.Background(), aKey).Result()
	if err != nil {
		return fmt.Errorf("remove account: %w", err)
	}
	_, err = s.redis.Del(context.Background(), tKey).Result()
	if err != nil {
		return fmt.Errorf("remove account: %w", err)
	}
	return nil
}

func (s *Storage) AccountExists(username string) (bool, error) {
	exists, err := s.redis.Exists(context.Background(), accountKey(username)).Result()
	if err != nil {
		return false, err
	}
	return exists == 1, nil
}

func (s *Storage) AddEmail(username string, email entity.Email) error {
	key := accountKey(username)
	exists, err := s.redis.Exists(context.Background(), key).Result()
	if err != nil {
		return fmt.Errorf("check account exists: %w", err)
	}
	if exists == 0 {
		return entity.ErrAccountDoesntExists
	}
	emailJSON, err := json.Marshal(email)
	if err != nil {
		return fmt.Errorf("json marshal email: %w", err)
	}
	_, err = s.redis.LPush(context.Background(), key, emailJSON).Result()
	if err != nil {
		return fmt.Errorf("lpush account email: %w", err)
	}
	return nil
}
