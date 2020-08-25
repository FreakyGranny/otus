package hw10_program_optimization //nolint:golint,stylecheck

import (
	"strings"
	"testing"
)

func BenchmarkCountDomains(b *testing.B) {
	var us users
	us[1] = User{
		Email: "john@gmail.com",
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		countDomains(us, "com")
	}
}

func BenchmarkGetUsers(b *testing.B) {
	r := strings.NewReader(`{"Id":1,"Name":"Howard Mendoza","Username":"0Oliver","Email":"aliquid_qui_ea@Browsedrive.gov","Phone":"6-866-899-36-79","Password":"InAQJvsq","Address":"Blackbird Place 25"}`)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		getUsers(r)
	}
}
