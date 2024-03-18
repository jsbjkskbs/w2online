package redis

func IsAvatarUploading(uid string) (bool, error) {
	exist, err := redisDBAvatarUpload.Exists(uid).Result()
	if err != nil {
		return true, err
	}
	return (exist != 0), nil
}

func AvatarSetUploadUncompleted(uid string) error {
	_, err := redisDBAvatarUpload.Set(uid, 1, 0).Result()
	if err != nil {
		return err
	}
	return nil
}

func AvatarSetUploadCompleted(uid string) error {
	_, err := redisDBAvatarUpload.Del(uid).Result()
	if err != nil {
		return err
	}
	return nil
}
