package utils

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	"strconv"
	"time"
)

type Cursor struct {
	Time time.Time
	ID   uuid.UUID
}

type Pagination struct {
	Previous *string `json:"previous"`
	Next     *string `json:"next"`
	Results  any     `json:"results"`
	Count    int     `json:"count"`
}

func MakePaginatedResp(result []map[string]any, baseUrl string, apiUrl string, limit int) Pagination {
	var Next *string
	var Previous *string
	if len(result) == 0 {
		return Pagination{}
	}
	if len(result) < limit {
		Next = nil
	} else {
		i, err := strconv.ParseInt(fmt.Sprintf("%v", result[limit-2]["created_at"]), 10, 64)
		if err != nil {
			GetAppLogger().Errorf("%v", err)
			return Pagination{}
		}
		tm := time.UnixMilli(i)
		nextDate := tm
		if err != nil {
			GetAppLogger().Errorf("%v", err)
			return Pagination{}
		}
		id, err := uuid.Parse(fmt.Sprintf("%v", result[limit-2]["id"]))
		if err != nil {
			GetAppLogger().Errorf("%v", err)
			return Pagination{}
		}
		// Format uses the same formatting style
		// as parse, or we can use a pre-made constant
		nextCursor := Cursor{
			Time: nextDate,
			ID:   id,
		}
		nextByteArray, err := json.Marshal(nextCursor)
		if err != nil {
			GetAppLogger().Errorf("%v", err)
			return Pagination{}
		}
		nextEnc := base64.StdEncoding.EncodeToString(nextByteArray)
		nextStr := baseUrl + apiUrl + "?cursor=" + nextEnc
		Next = &nextStr
		result = result[:limit-1]
	}
	prevI, err := strconv.ParseInt(fmt.Sprintf("%v", result[0]["created_at"]), 10, 64)
	if err != nil {
		GetAppLogger().Errorf("%v", err)
		return Pagination{}
	}
	prevTm := time.UnixMilli(prevI)
	prevDate := prevTm
	if err != nil {
		GetAppLogger().Errorf("%v", err)
		return Pagination{}
	}
	id, err := uuid.Parse(fmt.Sprintf("%v", result[0]["id"]))
	if err != nil {
		GetAppLogger().Errorf("%v", err)
		return Pagination{}
	}
	// Format uses the same formatting style
	// as parse, or we can use a pre-made constant
	prevCursor := Cursor{
		Time: prevDate,
		ID:   id,
	}
	prevByteArray, err := json.Marshal(prevCursor)
	if err != nil {
		GetAppLogger().Errorf("%v", err)
		return Pagination{}
	}
	prevEnc := base64.StdEncoding.EncodeToString(prevByteArray)

	prevStr := baseUrl + apiUrl + "?cursor=" + prevEnc

	Previous = &prevStr

	return Pagination{
		Previous: Previous,
		Next:     Next,
		Count:    len(result),
		Results:  result,
	}
}

func ParseCursor(cursor string) (*Cursor, error) {
	rawDecodedCursor, err := parseCursor(cursor)
	if err != nil {
		return nil, err
	}

	return rawDecodedCursor, nil
}

func parseCursor(cursor string) (*Cursor, error) {
	var returnCursor Cursor
	if cursor == "0" {
		returnCursor.Time = time.Now()
		returnCursor.ID = uuid.Nil
	} else {
		rawDecodedCursor, err := base64.StdEncoding.DecodeString(cursor)
		if err != nil {
			return nil, err
		}
		err = json.Unmarshal(rawDecodedCursor, &returnCursor)
		if err != nil {
			return nil, err
		}
	}
	return &returnCursor, nil
}
