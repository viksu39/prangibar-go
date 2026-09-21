package services

import "errors"

var (
	ErrEmailExists   = errors.New("email sudah digunakan")
	ErrNotFound      = errors.New("data tidak ditemukan")
	ErrInvalidToken  = errors.New("token tidak valid")
	ErrTokenNotFound = errors.New("token tidak ditemukan")
)
