package main

import (
	"errors"
	"io"
	"io/ioutil"
	"os"
	"log"
	"bufio"
	"fmt"
	"html/template"
)

type render struct {
	dm		dataManager
	gm		galleryManager
	source	string
	dest	string
}

func newRender(dm dataManager, gm galleryManager, source, dest string) (*render, error) {

	// check is source dir exists
	file, err := os.Stat(source)
	if err != nil {

		return nil, err

	}

	if !file.IsDir() {

		return nil, errors.New("newRender :: source is not a directory.")

	}

	// remove dest
	err = os.RemoveAll(dest)
	if err != nil {

		return nil, err

	}

	// re-create empty dir
	err = os.MkdirAll(dest, 0777)
	if err != nil {

		return nil, err

	}

	r := &render{dm, gm, source, dest}

	err = r.copyNormalFiles(source)
	if err != nil {

		return nil, err

	}

	err = r.render()
	if err != nil {

		return r, err

	}

	return r, nil

}

func (r *render) render() error {

	var p page
	err := r.dm.GetData(&p)
	if err != nil {

		return err

	}

	var g galleryData
	err = r.gm.GetData( &g )
	if err != nil {

		return err

	}

	var allData map[string]interface{} = make( map[string]interface{} )

	var data map[int]map[int]images = make( map[int]map[int]images )
	allData["one"] = data
	var row map[int]images

	for index, im := range g.One {

		if ( index % 4 ) == 0 {

			row = make( map[int]images )
			data[ index ] = row

		}

		row[index] = im

	}

	data = make( map[int]map[int]images )
	allData["two"] = data

	for index, im := range g.Two {

		if ( index % 4 ) == 0 {

			row = make( map[int]images )
			data[ index ] = row

		}

		row[index] = im

	}

	data = make( map[int]map[int]images )
	allData["three"] = data

	for index, im := range g.Three {

		if ( index % 4 ) == 0 {

			row = make( map[int]images )
			data[ index ] = row

		}

		row[index] = im

	}

	log.Println( allData )

	update := func(s, d string) error {

		t := template.New("template")
		t, err := t.ParseFiles(s)
		if err != nil {

			return err

		}

		out, err := os.Create(d)
		if err != nil {

			return err

		}

		defer out.Close()

		writer := bufio.NewWriter(out)
		defer writer.Flush()

		t.Execute(writer, map[string]interface{}{ "p":p, "g":allData })

		return nil

	}

	err = r.getJsts(update)
	if err != nil {

		return err

	}

	return nil

}

func (r *render) getJsts(fn func(s, d string) error, path ...string) error {

	var p string

	if len(path) == 0 {

		p = r.source

	} else {

		p = path[0]

	}

	fis, err := ioutil.ReadDir(p)
	if err != nil {

		return err

	}

	for _, fi := range fis {

		s := fmt.Sprintf("%s/%s", p, fi.Name())
		d := fmt.Sprintf("%s%s", r.dest, s[len(r.source):])

		if fi.IsDir() {

			// make dir for safety
			err = os.MkdirAll(d, 0777)
			if err != nil {

				return err

			}

			r.getJsts(fn, s)

		} else {

			l := len(s)

			if l >= 4 && s[l-4:] == ".jst" {

				err = fn(s, d[:len(d)-4])
				if err != nil {

					return err

				}

			}

		}

	}

	return nil

}

func (r *render) copyNormalFiles(path string) error {

	fis, err := ioutil.ReadDir(path)
	if err != nil {

		return err

	}

	for _, fi := range fis {

		s := fmt.Sprintf("%s/%s", path, fi.Name())
		d := fmt.Sprintf("%s%s", r.dest, s[len(r.source):])

		if fi.IsDir() {

			// create the dest dir
			err = os.MkdirAll(d, 0777)
			if err != nil {

				return err

			}

			// start 'copyNormalFiles' on this dir
			err = r.copyNormalFiles(s)
			if err != nil {

				return err

			}

		} else {

			n := fi.Name()
			l := len(n)

			if !(l >= 4 && n[l-4:] == ".jst") {

				err = copyFile(s, d)
				if err != nil {

					return err

				}

			}

		}

	}

	return nil

}

func copyFile(s, d string) error {

	// check source file
	in, err := os.Open(s)
	if err != nil {

		return err

	}

	// eventualy close source file
	defer in.Close()

	// create new file
	out, err := os.Create(d)
	if err != nil {

		return err

	}

	// eventually close new file
	defer out.Close()

	_, err = io.Copy(out, in)
	if err != nil {

		return err

	}

	return nil

}
