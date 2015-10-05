package main

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"log"
	"os"
)

type item struct {
	Title     string   `json:"title"`
	Paragraph []string `json:"paragraph"`
}

type page struct {
	One            item `json:"one"`
	Two            item `json:"two"`
	Three          item `json:"three"`
	Four           item `json:"four"`
	Five           item `json:"five"`
	LogoBrightness int  `json:"logobrightness"` // 0 : low -  1 : mid - 2 : high
}

type jsonDataBase struct {
	path string
}

func (db *jsonDataBase) GetDataJson() ([]byte, error) {

	var data []byte
	var err error

	data, err = ioutil.ReadFile(db.path)
	if err != nil {

		return data, err

	}

	return data, nil

}

func (db *jsonDataBase) GetData(result interface{}) error {

	r, ok := result.(*page)
	if !ok {

		log.Fatal("jsonDataBase :: GetData :: result must be of *page type.")

	}

	data, err := db.readData()
	if err != nil {

		return err

	}

	*r = data

	return nil

}

func (db *jsonDataBase) UpdateLogoBrightness(value int) error {

	if !(value >= 0 && value < 3) {

		return errors.New("jsonDataBase :: UpdateLogoBrightness :: value must be 0, 1 or 2.")

	}

	data, err := db.readData()
	if err != nil {

		return err

	}

	if data.LogoBrightness == value {

		return nil

	}

	data.LogoBrightness = value

	err = db.writeData(data)
	if err != nil {

		return err

	}

	return nil

}

func (db *jsonDataBase) UpdateTitle(value, page string) error {

	update := func(i *item) {
		i.Title = value
	}

	err := db.updateItem(page, update)
	if err != nil {

		return err

	}

	return nil

}

func (db *jsonDataBase) UpdateParagraph(value []string, page string) error {

	update := func(i *item) {
		i.Paragraph = value
	}

	err := db.updateItem(page, update)
	if err != nil {

		return err

	}

	return nil

}

func (db *jsonDataBase) updateItem(which string, fn func(i *item)) error {

	p, err := db.readData()
	if err != nil {

		return err

	}

	if which == "one" {

		fn(&p.One)

	} else if which == "two" {

		fn(&p.Two)

	} else if which == "three" {

		fn(&p.Three)

	} else if which == "four" {

		fn(&p.Four)

	} else if which == "five" {

		fn(&p.Five)

	} else {

		return errors.New("jsonDataBase :: insertItem :: which has not a valid value")

	}

	err = db.writeData(p)
	if err != nil {

		return err

	}

	return nil

}

func newJsonDatabase(path string) (*jsonDataBase, error) {

	db := &jsonDataBase{path}

	// check if file exists
	_, err := os.Stat(path)
	if err != nil {

		if os.IsNotExist(err) {

			// File doesn't exits, create new one with dummy data.
			errr := db.createInitDataAndFile()
			if errr != nil {

				return db, errr

			}

		} else {

			return nil, err

		}

	}

	return db, nil

}

func (db *jsonDataBase) createInitDataAndFile() error {

	p := page{
		One: item{
			Title:     "title one",
			Paragraph: []string{"This is paragraph one.", "Hello everybody."},
		},
		Two: item{
			Title:     "title two",
			Paragraph: []string{"This is paragraph two.", "Hello everybody."},
		},
		Three: item{
			Title:     "title three",
			Paragraph: []string{"This is paragraph three.", "Hello everybody."},
		},
		Four: item{
			Title:     "title four",
			Paragraph: []string{"This is paragraph four.", "Hello everybody."},
		},
		Five: item{
			Title:     "title five",
			Paragraph: []string{"This is paragraph five.", "Hello everybody."},
		},
		LogoBrightness: 1,
	}

	err := db.writeData(p)
	if err != nil {

		return err

	}

	return nil

}

func (db *jsonDataBase) writeData(p page) error {

	value, err := json.Marshal(p)
	if err != nil {

		return err

	}

	err = ioutil.WriteFile(db.path, value, 0777)
	if err != nil {

		return err

	}

	return nil

}

func (db *jsonDataBase) readData() (page, error) {

	var p page

	data, err := db.GetDataJson()
	if err != nil {

		return p, err

	}

	err = json.Unmarshal(data, &p)
	if err != nil {

		return p, err

	}

	return p, nil

}

// func( db *jsonDataBase ) readFile()
