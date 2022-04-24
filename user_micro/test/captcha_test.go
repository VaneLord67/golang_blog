package test

import (
	"testing"
	"user_micro/service"
)

func TestVerifyCaptcha(t *testing.T) {
	service.VerifyCaptcha("a", "a", "a")
}
