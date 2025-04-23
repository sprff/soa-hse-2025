package kafkaevents

import (
	"encoding/json"
	"fmt"
)

type Event struct {
	Event  string
	UserID string
	PostID string
}

func (u Event) String() string {
	v, err := json.Marshal(u)
	if err != nil {
		panic(fmt.Sprintf("can't marshal: %v", err))
	}
	return string(v)
}
func (u *Event) FromString(s string) {

	err := json.Unmarshal([]byte(s), u)
	if err != nil {
		panic(fmt.Sprintf("can't unmarshal: %v", err))
	}
}
