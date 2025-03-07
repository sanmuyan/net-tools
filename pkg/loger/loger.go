package loger

import (
	"bytes"
	"fmt"
	"github.com/sirupsen/logrus"
)

type SimpleFormatter struct {
}

func (f *SimpleFormatter) Format(entry *logrus.Entry) ([]byte, error) {
	var b bytes.Buffer
	b.WriteString(fmt.Sprintf("%s\n", entry.Message))
	return b.Bytes(), nil
}

var S *logrus.Logger

func init() {
	S = logrus.New()
	S.SetFormatter(&SimpleFormatter{})
}
