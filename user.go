package main

import (
	"os"
	"io/ioutil"
	"encoding/json"
	"errors"
)

type user struct {
	Name string		`json:"name"`
	Pw string		`json:"pw"`
}

type users struct {
	path string
}

func ( u *users ) GetUsername() ( string, error ) {

	data, err := u.readData()
	if err != nil {

		return "", err

	}

	if len( data ) != 1 {

		return "", errors.New( "There must only be one user." )

	}

	return data[0].Name, nil

}

func ( u *users ) Login( name, pw string ) ( bool, error ) {

	data, err := u.readData()
	if err != nil {

		return false, err

	}

	for _, us := range data {

		if us.Name == name && us.Pw == pw {

			return true, nil

		}

	}

	return false, nil

}

func ( u *users ) UpdateUser( origName, newName, newPw string ) ( error ) {

	data, err := u.readData()
	if err != nil {

		return err

	}

	var item int = -1

	for index, us := range data {

		if us.Name == origName {

			item = index

		}

	}

	if item == -1 {

		return errors.New( "user to update not found" )

	}

	data[item] = user{ newName, newPw }
	err = u.writeData( data )
	if err != nil {

		return err

	}

	return nil

}


func newUserManager( path string ) ( *users, error ) {

	u := &users{ path }

	_, err := os.Stat( path )
	if err != nil {

		if os.IsNotExist( err ) {

			// File does not exists, make it.
			err = u.createInitDataFile()
			if err != nil {

				return nil, err

			}

		} else {

			return nil, err

		}

	}

	return u, nil

}

func ( u *users ) createInitDataFile() ( error ) {

	var data []user = []user{ user{ "dasja", "dasja" } }

	err := u.writeData( data )
	if err != nil {

		return err

	}

	return nil

}

func ( u *users ) writeData( data []user ) ( error ) {

	value, err := json.Marshal( data )
	if err != nil {

		return err

	}

	err = ioutil.WriteFile( u.path, value, 0777 )
	if err != nil {

		return err

	}

	return nil

}

func ( u *users ) readData() ( []user, error ) {

	var value []user

	data, err := ioutil.ReadFile(u.path)
	if err != nil {

		return nil, err

	}

	err = json.Unmarshal( data, &value )
	if err != nil {

		return nil, err

	}

	return value, nil
	
}


























