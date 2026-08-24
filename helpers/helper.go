package helpers

var GlobalHelper *Helper

type Helper struct{}

func NewHelper() *Helper {
	GlobalHelper = &Helper{}

	return GlobalHelper
}
