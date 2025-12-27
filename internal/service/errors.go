package service

import "errors"

// 400 - Bad Request
var ErrInvalidTitle = errors.New("title cannot be empty")
var ErrInvalidID = errors.New("ID cannot be less then 1")

// 404 - not found
var ErrTaskNotFound = errors.New("task not found")
