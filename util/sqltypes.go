package util

import (
	"database/sql"
	"encoding/json"
)

type NullString struct {
	sql.NullString
}

func (ns *NullString) MarshalJSON() ([]byte, error) {
	if !ns.Valid {
		return []byte("null"), nil
	}
	return json.Marshal(ns.String)
}

func (ns *NullString) UnmarshalJSON(b []byte) error {
	err := json.Unmarshal(b, &ns.String)
	ns.Valid = err == nil
	return err
}

func NewNullStringFromString(s string) (NullString, error) {
	res := NullString{}
	if err := res.Scan(s); err != nil {
		return res, err
	}
	return res, nil
}
