package data

import (
	"fmt"
	"strconv"
)

type Runtime int32

func (r Runtime) MarshalJSON() ([]byte, error) {

	formatted := fmt.Sprintf("%d min", r)
	quoted := strconv.Quote(formatted)

	return []byte(quoted), nil
}
