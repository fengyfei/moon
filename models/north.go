package models

import (
	"database/sql"
	"encoding/json"
	"errors"
	"strconv"

	_ "github.com/mattn/go-sqlite3"
)

const (
	reportSummaryMinSize = 3
	reportDetailMinSize  = 10
)

var (
	errInvalidNorthSummarySize = errors.New("[North]: invalid summary size, 3 required")
	errInvalidNorthDetailSize  = errors.New("[North]: invalid detail size, 10 required")
	database                   *sql.DB
)

func init() {
	database, _ = sql.Open("sqlite3", "./north.db")

	stmt, err := database.Prepare(`
		CREATE TABLE IF NOT EXISTS summary (
			date TEXT PRIMARY KEY,
			ins NUMERIC,
			out NUMERIC,
			total NUMERIC
		)
	`)
	if err != nil {
		panic(err)
	}

	if _, err = stmt.Exec(); err != nil {
		panic(err)
	}

	stmt, err = database.Prepare(`
		CREATE TABLE IF NOT EXISTS detail (
			date TEXT,
			code TEXT,
			name TEXT,
			ins NUMERIC,
			out NUMERIC,
			total NUMERIC
		)
	`)
	if err != nil {
		panic(err)
	}

	if _, err = stmt.Exec(); err != nil {
		panic(err)
	}
}

type length interface {
	Size() int
	Save(date string) error
}

// NorthDailyReport -
type NorthDailyReport []length

// NorthSummaryItem -
type NorthSummaryItem struct {
	Total string `json:"total"`
}

// NorthDailySummary -
type NorthDailySummary struct {
	Data []NorthSummaryItem `json:"data"`
}

// Size -
func (n *NorthDailySummary) Size() int {
	return len(n.Data)
}

// Save -
func (n *NorthDailySummary) Save(date string) error {
	stmt, err := database.Prepare(`
		INSERT INTO summary(date, ins, out, total) VALUES(?, ?, ?, ?)
	`)

	in, err := strconv.ParseFloat(n.Data[0].Total, 64)
	if err != nil {
		return err
	}

	out, err := strconv.ParseFloat(n.Data[1].Total, 64)
	if err != nil {
		return err
	}

	total, err := strconv.ParseFloat(n.Data[2].Total, 64)
	if err != nil {
		return err
	}

	_, err = stmt.Exec(date, in, out, total)

	return err
}

// NorthDetailItem -
type NorthDetailItem struct {
	Code  string `json:"zqdm"`
	Name  string `json:"zqjc"`
	In    string `json:"mrjyje"`
	Out   string `json:"mcjyje"`
	Total string `json:"jyje"`
}

// NorthDailyDetail -
type NorthDailyDetail struct {
	Data []NorthDetailItem `json:"data"`
}

// Size -
func (n *NorthDailyDetail) Size() int {
	return len(n.Data)
}

// Save -
func (n *NorthDailyDetail) Save(date string) error {
	stmt, err := database.Prepare(`
		INSERT INTO detail(date, code, name, ins, out, total) VALUES(?, ?, ?, ?, ?, ?)
	`)

	if err != nil {
		return nil
	}

	for _, v := range n.Data {
		in, err := strconv.ParseFloat(v.In, 64)
		if err != nil {
			return err
		}

		out, err := strconv.ParseFloat(v.Out, 64)
		if err != nil {
			return err
		}

		total, err := strconv.ParseFloat(v.Total, 64)
		if err != nil {
			return err
		}

		if _, err = stmt.Exec(date, v.Code, v.Name, in, out, total); err != nil {
			return nil
		}
	}

	return nil
}

// Record -
func Record(date, data string) error {
	report := NorthDailyReport{
		&NorthDailySummary{
			Data: []NorthSummaryItem{},
		},
		&NorthDailyDetail{
			Data: []NorthDetailItem{},
		},
	}

	if err := json.Unmarshal([]byte(data), &report); err != nil {
		return err
	}

	if report[0].Size() != reportSummaryMinSize {
		return errInvalidNorthSummarySize
	}

	if report[1].Size() != reportDetailMinSize {
		return errInvalidNorthDetailSize
	}

	for _, v := range report {
		if err := v.Save(date); err != nil {
			return err
		}
	}

	return nil
}
