package main

type contextKey string

const (
	idContextKey      = contextKey("id")
	excerptContextKey = contextKey("excerpt")
)
