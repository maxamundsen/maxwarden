package security

import "testing"

func _BenchmarkDecrypt(b *testing.B) {
	for i := 0; i < b.N; i++ {
		DecryptSecret([]byte("njniYY9+R8kAxUuoI6p+A0AvDfVwtKVKe7FU7q7eW4IlLF1v4hLF14Fwizsddqh54EjiBB2XwD6g07c2Ovd0p8AehEuZgA8vD1N+3zSKKg+ZDVsc/MS+6iNQYK+ARNYHrqreaB2qiJP260Le3YR3xDY/u7n+JN58FxNf2J1DMvBUXD812d7r3ING4TBTkzcCJFXql+TvzUdC1qnhdrz/AOBo919rP2+yodQRTgBsZPiSb0DCZ9nnuwT9t99ORwn8v3AelyzwBOcxiYSlP07WDQE45o962E+GONiA09q8lBIBV6wT5bgZ3GAOdNNJFPrhSUqhblDB8/16Z1NwhS/lHyQUyjGxwt3zsC3axVCNQ6t4AJr8wEyVnoLb"), "password")
	}
}