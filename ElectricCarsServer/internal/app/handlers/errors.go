package handlers

import "errors"

var (
	idNotFound                  = errors.New("param `id` not found")
	idMustBeEmpty               = errors.New("param `id` must be empty")
	idCantBeEmpty               = errors.New("param `user_id` cant be empty")
	autopartNameCannotBeEmpty   = errors.New("autopart name cannot be empty")
	autopartBrandCannotBeEmpty  = errors.New("autopart name cannot be empty")
	autopartModelsCannotBeEmpty = errors.New("autopart name cannot be empty")
	autopartYearCannotBeEmpty   = errors.New("autopart name cannot be empty")
	autopartPriceCannotBeEmpty  = errors.New("autopart name cannot be empty")
	headerNotFound              = errors.New("no file uploaded")
	destinationOrCityIsEmpty    = errors.New("destination or city cannot be empty")
	serialNumberCannotBeEmpty   = errors.New("param `serial_number` cannot be empty")

	userIsNotModerator = errors.New("user is not moderator")
)
