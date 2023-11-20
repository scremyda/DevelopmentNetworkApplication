package handlers

import "errors"

var (
	idNotFound    = errors.New("param `id` not found")
	idMustBeEmpty = errors.New("param `id` must be empty")

	idCantBeEmpty               = errors.New("param `user_id` cant be empty or less 0")
	autopartNameCannotBeEmpty   = errors.New("autopart name cannot be empty or less 0")
	autopartBrandCannotBeEmpty  = errors.New("autopart name cannot be empty or less 0")
	autopartModelsCannotBeEmpty = errors.New("autopart name cannot be empty or less 0")
	autopartYearCannotBeEmpty   = errors.New("autopart name cannot be empty or less 0")
	autopartPriceCannotBeEmpty  = errors.New("autopart name cannot be empty or less 0")

	headerNotFound            = errors.New("no file uploaded")
	destinationOrCityIsEmpty  = errors.New("destination or city cannot be empty")
	serialNumberCannotBeEmpty = errors.New("param `serial_number` cannot be empty")

	userIsNotModerator = errors.New("user is not moderator")

	countInvalid = errors.New("invalid count param")
)
