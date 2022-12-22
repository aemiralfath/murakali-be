package model

import "time"

type Token struct {
	AccessToken  *AccessToken
	RefreshToken *RefreshToken
}

type AccessToken struct {
	Token     string
	ExpiredAt time.Time
}

type RefreshToken struct {
	Token     string
	ExpiredAt time.Time
}
