package utils

import (
	"fmt"
	"strconv"

	"github.com/andromaril/gophermmart/internal/errormart"
	"github.com/theplant/luhn"
	log "github.com/sirupsen/logrus"
)

func ValidLuhn(number string) (bool, error) {
	number2, err := strconv.Atoi(number)
	if err != nil {
		e := errormart.NewMartError(err)
		log.Error(e.Error())
		return false, fmt.Errorf("error %q", e.Error())
	}
	validnumber := luhn.Valid(number2)
	return validnumber, nil
}
