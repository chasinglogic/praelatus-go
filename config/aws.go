package config

func AWSBucket() string {
	return Cfg.AWS.Bucket
}

func AWSRegion() string {
	return Cfg.AWS.Region
}

func AWSBaseURL() *string {
	return Cfg.AWS.BaseURL
}

func AWSAccessKeyID() string {
	return Cfg.AWS.AccessKeyID
}

func AWSSecretKey() string {
	return Cfg.AWS.SecretKey
}
