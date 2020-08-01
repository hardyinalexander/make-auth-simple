#!/bin/bash

start:
	@go build -o bin/make-auth-simple .
	@bin/make-auth-simple