package postgresql

import "errors"

var ErrorInsertion = errors.New("can not insert")
var ErrorExists = errors.New("already exists")
var ErrorNotExists = errors.New("not exists")

var ErrorTagNotExists = errors.New("no such tag")