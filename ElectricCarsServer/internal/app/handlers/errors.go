package handlers

import "errors"

var (
	idNotFound    = errors.New("param `id` not found")
	idMustBeEmpty = errors.New("param `id` must be empty")

	idCantBeEmpty               = errors.New("param `user_id` cant be empty or less 0")
	autopartNameCannotBeEmpty   = errors.New("autopart name cannot be empty")
	autopartBrandCannotBeEmpty  = errors.New("autopart brand cannot be empty")
	autopartModelsCannotBeEmpty = errors.New("autopart model cannot be empty")
	autopartYearCannotBeEmpty   = errors.New("autopart year cannot be empty")
	autopartPriceCannotBeEmpty  = errors.New("autopart price cannot be empty")

	loginCantBeEmpty    = errors.New("param `login` cant be empty")
	passwordCantBeEmpty = errors.New("param `password` cant be empty")

	headerNotFound = errors.New("no file uploaded")

	userIsNotModerator = errors.New("user is not moderator")

	countInvalid = errors.New("invalid count param")
)
