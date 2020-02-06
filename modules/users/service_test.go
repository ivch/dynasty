package users

import (
	"io/ioutil"
	"os"
	"testing"

	"github.com/rs/zerolog"
)

var defaultLogger *zerolog.Logger

func TestMain(m *testing.M) {
	logger := zerolog.New(ioutil.Discard)
	defaultLogger = &logger
	os.Exit(m.Run())
}
