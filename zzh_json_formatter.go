package logrus

import (
	"bytes"
	"encoding/json"
	"fmt"
	//	"github.com/jimlawless/whereami"
	"net"
	"strings"
)

var (
	ip string
)

func init() {
	localIp()
}

type ZzhJSONFormatter JSONFormatter

// Format renders a single log entry
func (f *ZzhJSONFormatter) Format(entry *Entry) ([]byte, error) {
	data := make(Fields, len(entry.Data))
	var serialized []byte
	var err error
	for k, v := range entry.Data {
		switch v := v.(type) {
		case error:
			// Otherwise errors are ignored by `encoding/json`
			// https://github.com/sirupsen/logrus/issues/137
			data[k] = v.Error()
		default:
			data[k] = v
		}
	}
	if entry.Message != "" {
		serialized = []byte(entry.Message)
	} else {
		serialized, err = json.Marshal(data)
		if err != nil {
			return nil, fmt.Errorf("Failed to marshal fields to JSON, %v", err)
		}
	}
	timestampFormat := f.TimestampFormat
	if timestampFormat == "" {
		timestampFormat = defaultTimestampFormat
	}
	dt := entry.Time.Format(timestampFormat)
	level := entry.Level.String()

	return bytes.Join([][]byte{[]byte(dt), []byte(level), []byte(ip), []byte(""), serialized, []byte("\n")}, []byte("|")), err
	//return bytes.Join([][]byte{[]byte(dt), []byte(level), []byte(ip), []byte(whereami.WhereAmI(7)), serialized, []byte("\n")}, []byte("|")), err //for debug
}

func localIp() {
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		panic(err)
	}
	for _, addr := range addrs {
		a := addr.String()
		if strings.HasPrefix(a, "10.") {
			ip = a
			return
		}
	}
}
