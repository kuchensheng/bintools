package store

import "time"

type Options struct {

	//Expiration allows to specify an expiration time when setting a value
	Expiration time.Duration
}
