package main

import (
	"github.com/gorilla/mux"
	"github.com/gorilla/securecookie"
	"log"
	// "fmt"
	"net/http"
	"strconv"
	"encoding/json"
	"image/jpeg"
	"errors"
	"time"
)

type manager struct {
	sc         *securecookie.SecureCookie
	cookieName string
	dm         dataManager
	re         *render
	im 			imageManager
	gm			galleryManager
	um			userManager
}

func setManager(router *mux.Router, path string, dm dataManager, re *render, im imageManager, gm galleryManager, us userManager) {

	m := manager{securecookie.New(securecookie.GenerateRandomKey(64), securecookie.GenerateRandomKey(32)), "atelierdt", dm, re, im, gm, us}

	router.HandleFunc("/data", m.sessionHandler( m.dataHandler ) )
	router.HandleFunc("/img", m.sessionHandler( m.imgHandler) )
	router.HandleFunc("/update", m.sessionHandler( m.updateHandler) )
	router.HandleFunc("/upload", m.sessionHandler( m.uploadHandler) )
	router.HandleFunc("/thumb", m.sessionHandler( m.thumbHandler) )
	router.HandleFunc("/remove", m.sessionHandler( m.removeHandler) )
	router.HandleFunc("/updateimageinfo", m.sessionHandler( m.updateImageInfoHandler) )
	router.HandleFunc("/thumbUpdate", m.sessionHandler( m.thumbUpdateHandler) )
	router.HandleFunc("/login", m.loginHandler)
	router.HandleFunc("/logout", m.logoutHandler)
	router.HandleFunc("/username", m.sessionHandler( m.usernameHandler) )
	router.HandleFunc("/updatepw", m.sessionHandler( m.pwHandler) )

	router.PathPrefix("/").Handler(http.FileServer(http.Dir(path)))

}

func ( m manager ) pwHandler( rw http.ResponseWriter, req *http.Request ) {

	old := req.FormValue( "old" )
	user := req.FormValue( "user" )
	pass := req.FormValue( "pass" )

	err := m.um.UpdateUser( old, user, pass )
	if err != nil {

		log.Println( err )
		http.Error(rw, "Internal error.", http.StatusInternalServerError )
		return
		
	}

}

func ( m manager ) usernameHandler( rw http.ResponseWriter, req *http.Request ) {

	value, err := m.um.GetUsername()
	if err != nil {

		log.Println( err )
		http.Error(rw, "Internal error.", http.StatusInternalServerError )
		return

	}

	rw.Write( []byte( value ) )

}

func (m manager) updateImageInfoHandler( rw http.ResponseWriter, req *http.Request ) {

	which := req.FormValue( "which" )

	if !( which == "one" || which == "two" || which == "three" ) {

		log.Println( "manager :: removeHandler :: which type not valid." )
		http.Error(rw, "Type not found", http.StatusInternalServerError )
		return

	}

	name := req.FormValue( "name" )

	if len( name ) != 3 {

		log.Println( "manager :: updateImageInfo :: name type not valid." )
		http.Error(rw, "Type not found", http.StatusInternalServerError )
		return

	}

	text := req.FormValue( "text" )

	x, err := strconv.ParseInt( req.FormValue( "x" ), 10, 0 )
	if err != nil {

		log.Println( err )
		http.Error(rw, "Internal error.", http.StatusInternalServerError )
		return

	}

	y, err := strconv.ParseInt( req.FormValue( "y" ), 10, 0 )
	if err != nil {

		log.Println( err )
		http.Error(rw, "Internal error.", http.StatusInternalServerError )
		return

	}

	w, err := strconv.ParseInt( req.FormValue( "w" ), 10, 0 )
	if err != nil {

		log.Println( err )
		http.Error(rw, "Internal error.", http.StatusInternalServerError )
		return

	}

	h, err := strconv.ParseInt( req.FormValue( "h" ), 10, 0 )
	if err != nil {

		log.Println( err )
		http.Error(rw, "Internal error.", http.StatusInternalServerError )
		return

	}

	err = m.gm.UpdateImageInfo( which, name, text, x, y, w, h )
	if err != nil {

		log.Println( err )
		http.Error(rw, "Internal error.", http.StatusInternalServerError )
		return

	}

	err = m.re.render()
	if err != nil {

		log.Println(err)
		http.Error(rw, "Internal error", http.StatusInternalServerError)
		return

	}

}

func (m manager) removeHandler( rw http.ResponseWriter, req *http.Request ) {

	which := req.FormValue( "which" )

	if !( which == "one" || which == "two" || which == "three" ) {

		log.Println( "manager :: removeHandler :: which type not valid." )
		http.Error(rw, "Type not found", http.StatusInternalServerError )
		return

	}

	name := req.FormValue( "name" )

	if len( name ) != 3 {

		log.Println( "manager :: removeHandler :: name type not valid." )
		http.Error(rw, "Type not found", http.StatusInternalServerError )
		return

	}

	err := m.gm.RemoveImage( name, which )
	if err != nil {

		log.Println( err )
		http.Error(rw, "Internal error.", http.StatusInternalServerError )
		return

	}

	err = m.re.render()
	if err != nil {

		log.Println(err)
		http.Error(rw, "Internal error", http.StatusInternalServerError)
		return

	}

}

func (m manager) thumbUpdateHandler( rw http.ResponseWriter, req *http.Request ) {

	id := req.FormValue("which")
	log.Println(id)
	if !( id == "one" || id == "two" || id == "three" ) {

		log.Println( errors.New( "manager :: thumbHandler :: Illegal value for id/which" ) )
		http.Error( rw, "Error", http.StatusInternalServerError )
		return

	}

	name := req.FormValue( "name" )

	if len( name ) != 3 {

		log.Println( "manager :: removeHandler :: name type not valid." )
		http.Error(rw, "Type not found", http.StatusInternalServerError )
		return

	}

	text := req.FormValue("text")

	ix, err := strconv.ParseInt( req.FormValue("x"), 10, 0 )
	if err != nil {
		log.Println( err )
		http.Error( rw, "Error", http.StatusInternalServerError )
		return
	}

	iy, err := strconv.ParseInt( req.FormValue("y"), 10, 0 )
	if err != nil {
		log.Println( err )
		http.Error( rw, "Error", http.StatusInternalServerError )
		return
	}

	iw, err := strconv.ParseInt( req.FormValue("w"), 10, 0 )
	if err != nil {
		log.Println( err )
		http.Error( rw, "Error", http.StatusInternalServerError )
		return
	}

	ih, err := strconv.ParseInt( req.FormValue("h"), 10, 0 )
	if err != nil {
		log.Println( err )
		http.Error( rw, "Error", http.StatusInternalServerError )
		return
	}

	_, h, err := req.FormFile( "file" )
	if err != nil {
		log.Println( err )
		http.Error( rw, "Error", http.StatusInternalServerError )
		return
	}

	m.gm.UpdateImage( h, id, name, text, ix, iy, iw, ih )

	err = m.re.render()
	if err != nil {

		log.Println(err)
		http.Error(rw, "Internal error", http.StatusInternalServerError)
		return

	}

}

func (m manager) thumbHandler( rw http.ResponseWriter, req *http.Request ) {

	id := req.FormValue("which")
	log.Println(id)
	if !( id == "one" || id == "two" || id == "three" ) {

		log.Println( errors.New( "manager :: thumbHandler :: Illegal value for id/which" ) )
		http.Error( rw, "Error", http.StatusInternalServerError )
		return

	}

	text := req.FormValue("text")

	ix, err := strconv.ParseInt( req.FormValue("x"), 10, 0 )
	if err != nil {
		log.Println( err )
		http.Error( rw, "Error", http.StatusInternalServerError )
		return
	}

	iy, err := strconv.ParseInt( req.FormValue("y"), 10, 0 )
	if err != nil {
		log.Println( err )
		http.Error( rw, "Error", http.StatusInternalServerError )
		return
	}

	iw, err := strconv.ParseInt( req.FormValue("w"), 10, 0 )
	if err != nil {
		log.Println( err )
		http.Error( rw, "Error", http.StatusInternalServerError )
		return
	}

	ih, err := strconv.ParseInt( req.FormValue("h"), 10, 0 )
	if err != nil {
		log.Println( err )
		http.Error( rw, "Error", http.StatusInternalServerError )
		return
	}

	_, h, err := req.FormFile( "file" )
	if err != nil {
		log.Println( err )
		http.Error( rw, "Error", http.StatusInternalServerError )
		return
	}

	m.gm.AddImage( h, id, text, ix, iy, iw, ih )

	err = m.re.render()
	if err != nil {

		log.Println(err)
		http.Error(rw, "Internal error", http.StatusInternalServerError)
		return

	}

}

func (m manager) uploadHandler( rw http.ResponseWriter, req *http.Request ) {

	t := req.FormValue( "type" )

	ff, _, err := req.FormFile( "file" )
	if err != nil {

		log.Println( err )
		http.Error( rw, "Inter Error.", http.StatusInternalServerError )
		return
	}

	i, err := jpeg.Decode( ff )
	if err != nil {

		log.Println(err)
		http.Error( rw, "Not a jpeg.", http.StatusNotAcceptable )
		return

	}

	switch t {
	case "header":
		m.im.SetJumboImage( i )
		break
	case "headerone":
		m.im.SetHeaderOneImage( i )
		break
	case "headertwo":
		m.im.SetHeaderTwoImage( i )
		break
	default:
		log.Println( "header type not found." )
		http.Error(rw, "Type not found", http.StatusInternalServerError)
		return
		break
	}

}

func (m manager) updateHandler(rw http.ResponseWriter, req *http.Request) {

	t := req.FormValue("type")

	switch t {
	case "logobrightness":
		m.updateLogoBrightness(rw, req.FormValue("value"))
		break

	case "title":
		m.updateTitle(rw, req.FormValue("which"), req.FormValue("value"))
		break

	case "paragraph":
		m.updateParagraph( rw, req.FormValue("which"), req.FormValue("value"))
		break

	default:
		http.Error(rw, "Type not found", http.StatusInternalServerError)
		return
		break

	}

}

func (m manager) updateParagraph(rw http.ResponseWriter, which, value string) {

	if !m.checkWhich(rw, which) {

		return

	}

	var result []string
	err := json.Unmarshal( []byte( value ), &result )
	if err != nil {

		log.Println( err )
		http.Error( rw, "Internal error.", http.StatusInternalServerError )
		return

	}

	err = m.dm.UpdateParagraph( result, which )
	if err != nil {

		log.Println( err )
		http.Error( rw, "Internal error.", http.StatusInternalServerError )
		return

	}

	err = m.re.render()
	if err != nil {

		log.Println(err)
		http.Error(rw, "Internal error", http.StatusInternalServerError)
		return

	}

	valueByte := []byte(value)
	rw.Header().Set("Content-Length", strconv.Itoa(len(valueByte)))
	rw.Header().Set("Content-Type", "application/text")
	rw.Write(valueByte)

}

func (m manager) updateTitle(rw http.ResponseWriter, which, value string) {

	if !m.checkWhich(rw, which) {

		return

	}

	err := m.dm.UpdateTitle(value, which)
	if err != nil {

		log.Println(err)
		http.Error(rw, "Internal error", http.StatusInternalServerError)
		return

	}

	err = m.re.render()
	if err != nil {

		log.Println(err)
		http.Error(rw, "Internal error", http.StatusInternalServerError)
		return

	}

	valueByte := []byte(value)
	rw.Header().Set("Content-Length", strconv.Itoa(len(valueByte)))
	rw.Header().Set("Content-Type", "application/text")
	rw.Write(valueByte)

}

func (m manager) checkWhich(rw http.ResponseWriter, which string) bool {

	if !(which == "one" || which == "two" || which == "one" || which == "three" || which == "four" || which == "five") {

		log.Println("manager :: checkwich :: which must be one, two, three, four or five.")
		http.Error(rw, "Internal error", http.StatusInternalServerError)
		return false

	}

	return true

}

func (m manager) updateLogoBrightness(rw http.ResponseWriter, value string) {

	if !(value == "0" || value == "1" || value == "2") {

		log.Println("manager :: updateLogoBrightness :: value must be 0, 1 or 2.")
		http.Error(rw, "Internal error", http.StatusInternalServerError)
		return

	}

	valueInt, err := strconv.ParseInt(value, 10, 0)
	if err != nil {

		log.Println("manager :: updateLogoBrightness :: strconv issue")
		log.Println(err)
		http.Error(rw, "Internal error", http.StatusInternalServerError)
		return

	}

	err = m.dm.UpdateLogoBrightness(int(valueInt))
	if err != nil {

		log.Println(err)
		http.Error(rw, "Internal error", http.StatusInternalServerError)
		return

	}

	err = m.re.render()
	if err != nil {

		log.Println(err)
		http.Error(rw, "Internal error", http.StatusInternalServerError)
		return

	}

	valueByte := []byte(value)
	rw.Header().Set("Content-Length", strconv.Itoa(len(valueByte)))
	rw.Header().Set("Content-Type", "application/text")
	rw.Write(valueByte)

}

func (m manager) dataHandler(rw http.ResponseWriter, req *http.Request) {

	data, err := m.dm.GetDataJson()
	if err != nil {

		log.Println("manager :: dataHandler :: GetDataJson")
		log.Println(err)

		http.Error(rw, "Something went wrong.", http.StatusInternalServerError)
		return

	}

	m.writeJson(rw, data)

}

func (m manager) imgHandler(rw http.ResponseWriter, req *http.Request) {

	data, err := m.gm.GetDataJson()
	if err != nil {

		log.Println("manager :: dataHandler :: GetDataJson")
		log.Println(err)

		http.Error(rw, "Something went wrong.", http.StatusInternalServerError)
		return

	}

	m.writeJson(rw, data)

}

func (m manager) writeJson(rw http.ResponseWriter, d []byte) {

	rw.Header().Set("Content-Length", strconv.Itoa(len(d)))
	rw.Header().Set("Content-Type", "application/json")
	rw.Write(d)

}



// Login session
func (m manager) sessionHandler( fn func( rw http.ResponseWriter, req *http.Request ) ) ( http.HandlerFunc ) {

	return func( rw http.ResponseWriter, req *http.Request ) {

		loggedIn, err := m.isLoggedIn( req )
		if err != nil {

			log.Println(err)
			m.clearSession( rw )
			http.Error(rw, "There was an error.", http.StatusInternalServerError)
			return

		}

		if !loggedIn {

			http.Error(rw, "You must be logged in.", http.StatusUnauthorized)
			return

		}

		fn( rw, req )

	}

}

func ( m manager ) loginHandler( rw http.ResponseWriter, req *http.Request ) {

	name := req.FormValue( "name" )
	pw := req.FormValue( "pw" )

	if name == "" || pw == "" {
		log.Println( "loginHandler : name and/or pw empty." )
		http.Error( rw, "Did not receive enough values.", http.StatusBadRequest )
		return
	}

	ok, err := m.um.Login( name, pw )
	if err != nil {

		log.Println( err )
		http.Error( rw, "username or password not correct.", http.StatusUnauthorized )
		return

	}

	if ok {
		
		m.setSession( name, rw )
		return

	}

	http.Error( rw, "username or password not correct.", http.StatusUnauthorized )

}

func( m manager ) logoutHandler( rw http.ResponseWriter, req *http.Request ) {

	m.clearSession( rw )

}

func (m manager) setSession( name string, rw http.ResponseWriter ) {

	value := map[string]string{
		"name":name,
	}

	encodedValue, err := m.sc.Encode( m.cookieName, value )
	if err != nil {

		log.Println( err )
		http.Error( rw, "An error occured.", http.StatusInternalServerError )
		return

	}

	http.SetCookie( rw, &http.Cookie{
			Name : m.cookieName,
			Value : encodedValue,
			Path : "/",
			Expires : time.Now().Add( time.Hour*48 ),
		} )

}

func ( m manager ) clearSession( rw http.ResponseWriter ) {

	http.SetCookie( rw, &http.Cookie{
			Name : m.cookieName,
			Value : "",
			Path : "/",
			MaxAge : -1,
		} )

}

func ( m manager ) isLoggedIn( req *http.Request ) (bool, error) {

	cookie, err := req.Cookie( m.cookieName )
	if err != nil {

		if err == http.ErrNoCookie {

			return false, nil

		} else {

			return false, err

		}

	}

	cookieValue := make( map[string]string )
	err = m.sc.Decode( m.cookieName, cookie.Value, &cookieValue )
	if err != nil {

		return false, err

	}

	if cookieValue["name"] != "" {

		return true, nil

	}

	return false, nil

}