package main

type contextKey string

const (
	excerptContextKey   = contextKey("excerpt")
	wordIdContextKey    = contextKey("wid")
	letterPosContextKey = contextKey("lid")
	userContextKey      = contextKey("user")
)
