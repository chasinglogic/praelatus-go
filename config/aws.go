package config

func AWSRegion() string {
	return Cfg.AWS.Region
}

func AWSURL() *string {
	return Cfg.AWS.BaseURL
}
