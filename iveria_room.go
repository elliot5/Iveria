package main

const (
	ROOM_SPAWN = 1 << iota
)

type Room struct {
	bounds Rect
	flags  int
}
