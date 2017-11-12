package tlogrus

import (
	"bytes"
	"fmt"
	"reflect"
	"sort"
	"strconv"
	"time"

	"github.com/sirupsen/logrus"
)

// TimestampFormat for TLOG
const TimestampFormat = "2006-01-02 15:04:05"

// Formatter for TLOG
type Formatter struct{}

// Format logrus entries to msg|v1|v2|v3|...
func (f *Formatter) Format(entry *logrus.Entry) ([]byte, error) {
	var b *bytes.Buffer
	var keys []string

	for k := range entry.Data {
		keys = append(keys, k)
	}

	if entry.Buffer != nil {
		b = entry.Buffer
	} else {
		b = &bytes.Buffer{}
	}

	b.WriteString(entry.Message)

	sort.Slice(keys, func(i, j int) bool {
		vi, _ := strconv.Atoi(keys[i])
		vj, _ := strconv.Atoi(keys[j])
		return vi < vj
	})

	for _, key := range keys {
		b.WriteByte('|')
		f.appendValue(b, entry.Data[key])
	}

	b.WriteByte('\n')
	return b.Bytes(), nil
}

func (f *Formatter) appendValue(b *bytes.Buffer, data interface{}) {
	rv := reflect.ValueOf(data)
	switch v := data.(type) {
	case string:
		if rv.String() == "" {
			fmt.Fprint(b, "NULL")
		} else {
			fmt.Fprint(b, rv.String())
		}
	case int, int64:
		fmt.Fprint(b, rv.Int())
	case float64:
		fmt.Fprint(b, rv.Float())
	case *time.Time:
		fmt.Fprint(b, v.Format(TimestampFormat))
	default:
		fmt.Fprint(b, rv.Elem())
	}
}
