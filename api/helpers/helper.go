package helpers

import (
	"reflect"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

func PageOfSlice(s interface{}, page, pageSize int) reflect.Value {
	v := reflect.ValueOf(s)
	l := v.Len()
	i := (page - 1) * pageSize
	j := i + pageSize
	if i >= l {
		return reflect.New(reflect.TypeOf(s)).Elem()
	}
	if j > l {
		j = l
	}
	return v.Slice(i, j)
}

func GenerateColumnLetter(n int) string {
	if n <= 0 {
		return ""
	}
	const alphabetSize = 26

	firstLetter := string(rune('A' + (n-1)%alphabetSize))

	remainingPart := GenerateColumnLetter((n - 1) / alphabetSize)

	return remainingPart + firstLetter
}

// GetQueryInt retrieves an integer query parameter, returning a default value if not present or invalid.
func GetQueryInt(c *gin.Context, key string, defaultValue int) int {
	if valueStr := c.Query(key); valueStr != "" {
		if value, err := strconv.Atoi(valueStr); err == nil {
			return value
		}
	}
	return defaultValue
}

func ParseRFC3339(timestamp string) (time.Time, error) {
	return time.Parse(time.RFC3339, timestamp)
}
