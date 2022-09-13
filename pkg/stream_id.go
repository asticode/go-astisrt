package astisrt

import (
	"errors"
	"sort"
	"strings"
)

// https://github.com/Haivision/srt/blob/master/docs/features/access-control.md

type StreamIDKey string

const (
	StreamIDKeyHostName      StreamIDKey = "h"
	StreamIDKeyMode          StreamIDKey = "m"
	StreamIDKeyRessourceName StreamIDKey = "r"
	StreamIDKeySessionID     StreamIDKey = "s"
	StreamIDKeyType          StreamIDKey = "t"
	StreamIDKeyUserName      StreamIDKey = "u"
)

type StreamID map[StreamIDKey]StreamIDItem

type StreamIDItem struct {
	Children map[StreamIDKey]StreamIDItem
	Value    string
}

const streamIDPrefix = "#!:"

func ParseStreamID(str string) (s StreamID, err error) {
	// Create stream id
	s = make(StreamID)

	// Handle prefix
	if !strings.HasPrefix(str, streamIDPrefix) {
		err = errors.New("astisrt: missing prefix")
		return
	}
	str = strings.TrimPrefix(str, streamIDPrefix)

	// Check format
	switch str[0] {
	// Comma-separated key-value pairs with no nesting
	case ':':
		// Loop through key/value pairs
		for _, pair := range strings.Split(strings.TrimPrefix(str, ":"), ",") {
			// Split on =
			split := strings.Split(pair, "=")
			if len(split) < 2 {
				continue
			}

			// Add child
			s[StreamIDKey(split[0])] = StreamIDItem{Value: split[1]}
		}
	// Nested block with one or several key-value pairs
	case '{':
		idx := 1
		parseStreamIDNestedFunc(str, &idx, s)
	// Invalid
	default:
		err = errors.New("astisrt: invalid format")
		return
	}
	return
}

func parseStreamIDNestedFunc(str string, idx *int, m map[StreamIDKey]StreamIDItem) {
	var key string
	var i StreamIDItem
	isKey := true
	for ; *idx < len(str); *idx++ {
		c := str[*idx]
		switch c {
		case '{':
			*idx++
			i.Children = make(map[StreamIDKey]StreamIDItem)
			parseStreamIDNestedFunc(str, idx, i.Children)
		case '=':
			isKey = false
		case ',', '}':
			isKey = true
			m[StreamIDKey(key)] = i
			if c == ',' {
				key = ""
				i = StreamIDItem{}
			} else {
				return
			}
		default:
			if isKey {
				key += string(c)
			} else {
				i.Value += string(c)
			}
		}
	}
}

func (s StreamID) String() (str string) {
	// Add prefix
	str += streamIDPrefix

	// Nothing to do
	if len(s) == 0 {
		return
	}

	// Get format
	format := ':'
	for _, c := range s {
		if len(c.Children) > 0 {
			format = '{'
			break
		}
	}

	// Switch on format
	switch format {
	case '{':
		str += stringStreamIDNested(s)
	default:
		var ss []string
		for k, v := range s {
			ss = append(ss, string(k)+"="+v.Value)
		}
		str += ":" + strings.Join(ss, ",")
	}
	return
}

func stringStreamIDNested(m map[StreamIDKey]StreamIDItem) string {
	// Sort keys
	var ks []string
	for k := range m {
		ks = append(ks, string(k))
	}
	sort.Strings(ks)

	// Loop through keys
	var ss []string
	for _, k := range ks {
		s := k + "="
		v := m[StreamIDKey(k)]
		if len(v.Children) > 0 {
			s += stringStreamIDNested(v.Children)
		} else {
			s += v.Value
		}
		ss = append(ss, s)
	}
	return "{" + strings.Join(ss, ",") + "}"
}

func (s StreamID) MarshalText() ([]byte, error) {
	return []byte(s.String()), nil
}

func (s *StreamID) UnmarshalText(b []byte) (err error) {
	*s, err = ParseStreamID(string(b))
	return err
}
