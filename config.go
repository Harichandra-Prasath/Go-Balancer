package main

type Config struct {
	Port        int
	Backends    []string
	STATIC_ROOT string
	MEDIA_ROOT  string
	ALGO        string
}
