package main

type contextKey string

const (
	excerptContextKey   = contextKey("excerpt")
	wordIndexContextKey = contextKey("wordIndex")
)
