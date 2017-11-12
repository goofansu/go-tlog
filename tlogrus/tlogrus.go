package tlogrus

import (
	"io"
	"strconv"
	"sync"

	validator "gopkg.in/go-playground/validator.v9"

	"github.com/fatih/structs"
	"github.com/sirupsen/logrus"
)

var (
	mu       = &sync.Mutex{}
	log      = logrus.New()
	validate = validator.New()
)

func init() {
	log.Formatter = new(Formatter)
}

// SetOutput sets the log output.
func SetOutput(out io.Writer) {
	mu.Lock()
	defer mu.Unlock()
	log.Out = out
}

// AddHook adds a hook to the log.
func AddHook(hook logrus.Hook) {
	mu.Lock()
	defer mu.Unlock()
	log.Hooks.Add(hook)
}

// Log event in TLOG format.
func Log(event interface{}) error {
	if err := validate.Struct(event); err != nil {
		return err
	}
	name := name(event)
	fields := fields(event)
	log.WithFields(fields).Info(name)
	return nil
}

func name(event interface{}) string {
	return structs.Name(event)
}

func fields(event interface{}) logrus.Fields {
	fieldsInStruct := structs.Fields(event)
	fields := make(logrus.Fields, len(fieldsInStruct))
	for i, f := range fieldsInStruct {
		if f.Tag("tlog") == "" {
			continue
		}
		fields[strconv.Itoa(i)] = f.Value()
	}
	return fields
}
