package redis

import (
	"work/pkg/errmsg"
)

func PutFollowList(uid string, followList *[]string) error {
	exist, err := redisDBRelationInfo.Exists(`s:` + uid).Result()
	if err != nil {
		return errmsg.RedisError
	}
	pipe := redisDBRelationInfo.TxPipeline()
	if exist != 0 {
		pipe.Del(`s:` + uid)
	}
	for _, item := range *followList {
		pipe.SAdd(`s:`+uid, item)
	}
	if _, err := pipe.Exec(); err != nil {
		return errmsg.RedisError
	}
	return nil
}

func PutFollowerList(uid string, followerList *[]string) error {
	exist, err := redisDBRelationInfo.Exists(`f:` + uid).Result()
	if err != nil {
		return errmsg.RedisError
	}
	pipe := redisDBRelationInfo.TxPipeline()
	if exist != 0 {
		pipe.Del(`f:` + uid)
	}
	for _, item := range *followerList {
		pipe.SAdd(`f:`+uid, item)
	}
	if _, err := pipe.Exec(); err != nil {
		return errmsg.RedisError
	}
	return nil
}

func AppendFollow(uid, followID string) error {
	if _, err := redisDBRelationInfo.SAdd(`s:`+uid, followID).Result(); err != nil {
		return errmsg.RedisError
	}
	return nil
}

func AppendFollower(uid, followerID string) error {
	if _, err := redisDBRelationInfo.SAdd(`f:`+uid, followerID).Result(); err != nil {
		return errmsg.RedisError
	}
	return nil
}

func RemoveFollow(uid, followID string) error {
	exist, err := redisDBRelationInfo.SIsMember(`s:`+uid, followID).Result()
	if err != nil {
		return errmsg.RedisError
	}
	if !exist {
		return errmsg.ServiceError
	}
	if _, err := redisDBRelationInfo.SRem(`s:`+uid, followID).Result(); err != nil {
		return errmsg.RedisError
	}
	return nil
}

func RemoveFollower(uid, followerID string) error {
	exist, err := redisDBRelationInfo.SIsMember(`f:`+uid, followerID).Result()
	if err != nil {
		return errmsg.RedisError
	}
	if !exist {
		return errmsg.ServiceError
	}
	if _, err := redisDBRelationInfo.SRem(`f:`+uid, followerID).Result(); err != nil {
		return errmsg.RedisError
	}
	return nil
}

func GetFollowList(uid string) (*[]string, error) {
	list, err := redisDBRelationInfo.SMembers(`s:` + uid).Result()
	if err != nil {
		return nil, errmsg.RedisError
	}
	return &list, nil
}

func GetFollowerList(uid string) (*[]string, error) {
	list, err := redisDBRelationInfo.SMembers(`f:` + uid).Result()
	if err != nil {
		return nil, errmsg.RedisError
	}
	return &list, nil
}

func GetFriendList(uid string) (*[]string, error) {
	list, err := redisDBRelationInfo.SInter(`f:`+uid, `s:`+uid).Result()
	if err != nil {
		return nil, errmsg.RedisError
	}
	return &list, nil
}

func DoesUserFollow(uid, userFollowId string) (bool, error) {
	exist, err := redisDBRelationInfo.SIsMember(`s:`+uid, userFollowId).Result()
	if err != nil {
		return true, errmsg.RedisError
	}
	return exist, nil
}

func IsUserFollower(uid, userFollowerId string) (bool, error) {
	exist, err := redisDBRelationInfo.SIsMember(`f:`+uid, userFollowerId).Result()
	if err != nil {
		return false, errmsg.RedisError
	}
	return exist, nil
}

func GetFollowCount(uid string) (int64, error) {
	count, err := redisDBRelationInfo.SCard(`s:` + uid).Result()
	if err != nil {
		return -1, errmsg.RedisError
	}
	return count, nil
}

func GetFollowerCount(uid string) (int64, error) {
	count, err := redisDBRelationInfo.SCard(`f:` + uid).Result()
	if err != nil {
		return -1, errmsg.RedisError
	}
	return count, nil
}
