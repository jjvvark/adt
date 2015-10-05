package main

import (
	"flag"
	"github.com/gorilla/mux"
	"log"
	"mauscode/configurationfile"
	"net/http"
	"image"
	"mime/multipart"
)

var (
	conf                string
	configurationValues = map[string]string{
		"manager":     "/Users/joostvanvark/www/adt-www/manager",
		"source":      "/Users/joostvanvark/www/adt-www/source",
		"dest":        "/Users/joostvanvark/www/adt-www/dest",
		"portManager": ":8000",
		"portDest":":8001",
		"data": "data.json",
		"user": "user.json",
		"images": "images.json",
	}
)

type dataManager interface {
	GetDataJson() ([]byte, error)
	GetData(result interface{}) error
	UpdateTitle(value, page string) error
	UpdateParagraph(value []string, page string) error
	UpdateLogoBrightness(value int) error
}

type imageManager interface {
	SetJumboImage( im image.Image ) ( error )
	SetHeaderOneImage( im image.Image ) ( error )
	SetHeaderTwoImage( im image.Image ) ( error )
}

type galleryManager interface {
	GetDataJson() ([]byte, error)
	GetData(result interface{}) ( error )
	AddImage( header *multipart.FileHeader, which, text string, x, y, w, h int64 ) ( error )
	UpdateImageInfo( which, name, text string, x, y, w, h int64 ) ( error )
	UpdateImage( header *multipart.FileHeader, which, name, text string, x, y, w, h int64 ) ( error )
	RemoveImage( name, which string ) ( error )
}

type userManager interface {
	Login( name, pw string ) ( bool, error )
	UpdateUser( origName , newName, newPw string ) ( error )
	GetUsername() ( string, error )
}

func init() {
	flag.StringVar(&conf, "conf", "settings.conf", "Setting configuration file location.")
	flag.Parse()
}

func main() {

	// init configuation file, exit if non existant.
	configurationValues, err := configurationfile.Parse(configurationValues, conf)
	if err != nil {

		log.Fatal(err)

	}

	// init usermanager
	var um userManager
	um, err = newUserManager( configurationValues[ "user" ] )
	if err != nil {

		log.Fatal( err )

	}

	//init datamanager
	var dm dataManager
	dm, err = newJsonDatabase(configurationValues["data"])
	if err != nil {

		log.Fatal(err)

	}

	// init gallery -> before renderer....
	var gm galleryManager
	gm, err = newGallery( configurationValues["manager"], configurationValues["source"], configurationValues["dest"], configurationValues["images"] )
	if err != nil {

		log.Fatal(err)

	}

	// init renderer
	re, err := newRender(dm, gm, configurationValues["source"], configurationValues["dest"])
	if err != nil {

		log.Fatal(err)

	}

	var im imageManager
	im = newImager( configurationValues["manager"], configurationValues["source"], configurationValues["dest"] )

	// init manager
	router := mux.NewRouter()
	setManager(router, configurationValues["manager"], dm, re, im, gm, um)
	if err != nil {

		log.Fatal(err)

	}

	// start manager
	go func() {
		log.Fatal(http.ListenAndServe(configurationValues["portManager"], router))
	}()

	site := mux.NewRouter()
	site.PathPrefix( "/" ).Handler( http.FileServer( http.Dir( configurationValues[ "dest" ] ) ) )

	// start site
	go func() {
		log.Fatal(http.ListenAndServe(configurationValues["portDest"], site))
	}()

	// loop forever
	select {}

}
