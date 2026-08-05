package errorlist

import "errors"

var UsernameCantContainAt = errors.New("username cannot contain '@'")
var UserNotFound = errors.New("user not found")
var EmailAlreadyExists = errors.New("email already exists")
var UsernameAlreadyExists = errors.New("username already exists")
var InvalidBirthDate = errors.New("invalid birth date")
var InvalidCredentials = errors.New("invalid credentials")
