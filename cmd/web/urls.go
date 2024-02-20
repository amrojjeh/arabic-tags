package main

import "fmt"

type url struct{}

func (u url) index() string {
	return "/"
}

func (u url) login() string {
	return "/login"
}

func (u url) register() string {
	return "/register"
}

func (u url) logout() string {
	return "/logout"
}

func (u url) home() string {
	return "/home"
}

func (u url) createExcerpt() string {
	return "/excerpt"
}

func (u url) excerpt(id any) string {
	return fmt.Sprintf("/excerpt/%v", id)
}

func (u url) excerptTitle(id any) string {
	return fmt.Sprintf("/excerpt/%v/title", id)
}

func (u url) excerptLock(id any) string {
	return fmt.Sprintf("/excerpt/%v/lock", id)
}
