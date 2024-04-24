package errormart

import (
	"fmt"
	"time"
)

type MartError struct {
	Time time.Time
	Err  error
}

func (te *MartError) Error() string {
	return fmt.Sprintf("%v %v", te.Time.Format("2006/01/02 15:04:05"), te.Err)
}

func NewMartError(err error) error {
	return &MartError{
		Time: time.Now(),
		Err:  err,
	}
}
