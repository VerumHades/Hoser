package util

import (
	"encoding/base64"
	"encoding/json"
)

func EncodeCursor[CursorType any](cursor CursorType) (string, error) {
	data, err := json.Marshal(cursor)
	if err != nil {
		return "", err
	}
	return base64.URLEncoding.EncodeToString(data), nil
}

func DecodeCursor[CursorType any](s string) (CursorType, error) {
	var cursor CursorType
	data, err := base64.URLEncoding.DecodeString(s)
	if err != nil {
		return cursor, err
	}
	err = json.Unmarshal(data, &cursor)
	return cursor, err
}
