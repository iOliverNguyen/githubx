package dot

import (
	"strconv"
)

func JoinIDs(ids []ID) string {
	b := make([]byte, 0, 20*len(ids))
	for _, id := range ids {
		b = strconv.AppendInt(b, id.Int64(), 10)
		b = append(b, ',')
	}
	if len(b) > 0 {
		b = b[:len(b)-1]
	}
	return string(b)
}
