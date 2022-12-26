package consts

import "github.com/rs/zerolog/log"

type exception struct {
	Location    string `json:"location"`
	Name        string `json:"name"`
	Description string `json:"description"`
}

func (e *exception) Error() string {
	return e.Description
}
func NewException(location, name, desc string) *exception {
	return &exception{location, name, desc}
}

func DeferHandler() error {
	if x := recover(); x != nil {
		log.Error().Msgf("发生了panic错误:%v", x.(error))
		return x.(error)
	}
	return nil
}
