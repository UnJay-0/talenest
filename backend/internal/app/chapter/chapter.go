package chapter

import "talenest/backend/internal/app/tales"

type Chapter struct {
	id        int
	content   string
	sentiment float64
	tale      tales.Tale
}
