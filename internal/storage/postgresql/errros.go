package postgresql

import "errors"

var ErrorInsertion = errors.New("can not insert")
var ErrorExists = errors.New("already exists")