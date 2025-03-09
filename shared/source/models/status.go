package models

import (
	"encoding/json"
	"fmt"
	"log/slog"
)

type StatusStr string

const (
	StatusOk StatusStr = "OK"
)

type StatusData struct {
	Status StatusStr `json:"status"`
	Data   any       `json:"data"`
}

func NewStatusData(data any, err error) StatusData {
	res := StatusData{
		Status: StatusOk,
		Data:   data,
	}
	if err != nil {
		res.Status = StatusStr(err.Error()) //TODO replace to some unique codes
		res.Data = err
	}
	return res
}

func ParseStatusData[T any](sd StatusData) (data T, err error) {
	if sd.Status == StatusOk {
		data = parseAs[T](sd.Data)
		return
	}
	// By unique code decide type of errors
	err = fmt.Errorf("%v", sd.Status)
	return
}

func parseAs[T any](data any) (out T) {
	outBytes, err := json.Marshal(data)
	if err != nil {
		slog.Error("Can't marshal data", "data", data, "error", err)
		return
	}

	err = json.Unmarshal(outBytes, &out)
	if err != nil {
		slog.Error("Can't unmarshsa data", "data", outBytes, "error", err, "type", fmt.Sprintf("%T", out))
		return
	}
	return
}
