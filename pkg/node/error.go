package node

import "errors"

var ErrBadPath = errors.New("bad hierarchy")
var ErrAlreadyExists = errors.New("already exists")
var ErrNotFound = errors.New("not found")
