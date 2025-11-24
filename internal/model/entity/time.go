package entity

import (
	"encoding/json"
	"time"
)

type BizTimeFull time.Time

func (t BizTimeFull) MarshalJSON() ([]byte, error) {
	str := time.Time(t).Format("2006-01-02 15:04:05")
	return json.Marshal(str) // 自动加引号
}

func (t BizTimeFull) String() string {
	return time.Time(t).Format("2006-01-02 15:04:05")
}
