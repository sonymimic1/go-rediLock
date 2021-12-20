#/bin/bash

NAME = $(shell basename $(shell pwd -P))
OUTPUT = $(NAME)
TARGET = ./main.go
LINUX_OS = linux
LINUX_ARCH = amd64


docker-redis:
	docker run -itd --name redis-test -p 6379:6379 redis




