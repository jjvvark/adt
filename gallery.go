package main

import (
	"fmt"
	"log"
	"os"
	"io"
	"image"
	"image/draw"
	"image/jpeg"
	"io/ioutil"
	"encoding/json"
	"github.com/nfnt/resize"
	"mime/multipart"
	"strconv"
)

type gallery struct {
	man string
	source string
	dest string
	path string
}

const (
	FOLDER_ONE string = "one"
	FOLDER_TWO string = "two"
	FOLDER_THREE string = "three"
)

type images struct {
	Id string			`json:"id"`
	Nm int				`json:"nm"`
	Text string			`json:"text"`
	X int64				`json:"x"`
	Y int64				`json:"y"`
	W int64				`json:"w"`
	H int64				`json:"h"`

}

type galleryData struct {
	One		[]images		`json:"one"`
	Two		[]images		`json:"two"`
	Three	[]images		`json:"three"`
}

func newGallery( man, source, dest, path string ) ( *gallery, error ) {

	g := &gallery{ fmt.Sprintf( "%s/pics", man ), fmt.Sprintf( "%s/pics", source ), fmt.Sprintf( "%s/pics", dest ), path }

	// check source dirs
	err := g.checkFiles( fmt.Sprintf( "%s/pics/%s", source, FOLDER_ONE ),
		fmt.Sprintf( "%s/pics/%s", source, FOLDER_TWO ),
		fmt.Sprintf( "%s/pics/%s", source, FOLDER_THREE ),
		fmt.Sprintf( "%s/pics/%s", man, FOLDER_ONE ),
		fmt.Sprintf( "%s/pics/%s", man, FOLDER_TWO ),
		fmt.Sprintf( "%s/pics/%s", man, FOLDER_THREE ) )
	if err != nil {

		return nil, err

	}

	// check if data file exists
	_, err = os.Stat( path )
	if err != nil {

		if os.IsNotExist( err ) {

			// data file does not exist
			g.resetData()

		} else {

			return nil, err

		}

	}

	return g, nil

}

func ( g *gallery ) RemoveImage( name, which string ) ( error ) {

	d, err := g.readData()
	if err != nil {

		return err

	}

	if which == "one" {

		var newData []images = make( []images, 0 )

		for _, v := range d.One {

			if v.Id != name {

				newData = append( newData, v )

			}

		}

		d.One = newData

	} else if which == "two" {

		var newData []images = make( []images, 0 )

		for _, v := range d.Two {

			if v.Id != name {

				newData = append( newData, v )

			}

		}

		d.Two = newData

	} else if which == "three" {

		var newData []images = make( []images, 0 )

		for _, v := range d.Three {

			if v.Id != name {

				newData = append( newData, v )

			}

		}

		d.Three = newData

	}

	err = g.writeData( d )
	if err != nil {

		return err

	}

	// remove all files
	err = os.RemoveAll( fmt.Sprintf( "%s/%s/%s", g.man, which, name ) )
	if err != nil {

		return err

	}

	err = os.RemoveAll( fmt.Sprintf( "%s/%s/%s", g.source, which, name ) )
	if err != nil {

		return err

	}
	
	err = os.RemoveAll( fmt.Sprintf( "%s/%s/%s", g.dest, which, name ) )
	if err != nil {

		return err

	}
	

	return nil

}

func ( g *gallery ) GetData(result interface{}) ( error ) {

	r, ok := result.( *galleryData )
	if !ok {

		log.Fatal("gallery :: GetData :: result must be of *galleryData type.")

	}

	data, err := g.readData()
	if err != nil {

		return err

	}

	*r = data

	return nil

}

func ( g *gallery ) UpdateImageInfo( which, name, text string, x, y, w, h int64 ) ( error ) {

	data, err := g.readData()
	if err != nil {
		return err
	}

	if which == FOLDER_ONE {
		
		for index, i := range data.One {

			if i.Id == name {

				data.One[ index ] = images{ data.One[ index ].Id, data.One[ index ].Nm, text, x, y, w, h }

			}

		}

	} else if which == FOLDER_TWO {

		for index, i := range data.Two {

			if i.Id == name {

				data.Two[ index ] = images{ data.Two[ index ].Id, data.Two[ index ].Nm, text, x, y, w, h }

			}

		}

	} else if which == FOLDER_THREE {

		for index, i := range data.Three {

			if i.Id == name {

				data.Three[ index ] = images{ data.Three[ index ].Id, data.Three[ index ].Nm, text, x, y, w, h }

			}

		}

	}

	err = g.writeData( data )
	if err != nil {

		return err

	}

	// recreate thumb from original...
	file, err := os.Open( fmt.Sprintf( "%s/%s/%s/orig.jpg", g.man, which, name ) )
	if err != nil {

		return err

	}

	im, err := jpeg.Decode( file )
	if err != nil {

		return err

	}

	thumb := g.createThumb( im, x, y, w, h )

	t1, err := os.Create( fmt.Sprintf( "%s/%s/%s/%s", g.man, which, name, "thumb.jpg" ) )
	if err != nil {
		return err
	}

	defer t1.Close()

	t2, err := os.Create( fmt.Sprintf( "%s/%s/%s/%s", g.source, which, name, "thumb.jpg" ) )
	if err != nil {
		return err
	}

	defer t2.Close()

	t3, err := os.Create( fmt.Sprintf( "%s/%s/%s/%s", g.dest, which, name, "thumb.jpg" ) )
	if err != nil {
		return err
	}

	defer t3.Close()

	writer := io.MultiWriter( t1, t2, t3 )
	err = jpeg.Encode( writer, thumb, nil )
	if err != nil {

		return err

	}

	return nil

}

func ( g *gallery ) UpdateImage( header *multipart.FileHeader, which, name, text string, x, y, w, h int64 ) ( error ) {

	ni, err := strconv.ParseInt( name, 10, 0 )
	if err != nil {

		return err

	}

	var nameInt int = int( ni )

	file, err := header.Open()
	if err != nil {

		return err

	}

	im, err := jpeg.Decode( file )
	if err != nil {

		return err

	}

	thumb := g.createThumb( im, x, y, w, h )

	ww := im.Bounds().Max.X - im.Bounds().Min.X
	hh := im.Bounds().Max.Y - im.Bounds().Min.Y

	var img image.Image

	if ww > hh {

		// landscape
		img = resize.Resize( 1200, 0, im, resize.Lanczos3 )


	} else {

		// portait or rectangle
		img = resize.Resize( 0, 800, im, resize.Lanczos3 )

	}

	data, err := g.readData()
	if err != nil {
		return err
	}

	var fol string

	if which == FOLDER_ONE {

		fol = FOLDER_ONE
		
		for index, i := range data.One {

			if i.Id == name {

				data.One[ index ] = images{ name, nameInt, text, x, y, w, h }

			}

		}

	} else if which == FOLDER_TWO {

		fol = FOLDER_TWO
		
		for index, i := range data.Two {

			if i.Id == name {

				data.Two[ index ] = images{ name, nameInt, text, x, y, w, h }

			}

		}

	} else if which == FOLDER_THREE {

		fol = FOLDER_THREE
		
		for index, i := range data.Three {

			if i.Id == name {

				data.Three[ index ] = images{ name, nameInt, text, x, y, w, h }

			}

		}

	}

	err = g.writeData( data )
	if err != nil {

		return err

	}

	// write files
	err = g.createFolders( fmt.Sprintf( "%s/%s/%s", g.man, fol, name ), fmt.Sprintf( "%s/%s/%s", g.source, fol, name ), fmt.Sprintf( "%s/%s/%s", g.dest, fol, name ) )
	if err != nil {

		return err

	}

	t1, err := os.Create( fmt.Sprintf( "%s/%s/%s/%s", g.man, fol, name, "thumb.jpg" ) )
	if err != nil {
		return err
	}

	defer t1.Close()

	t2, err := os.Create( fmt.Sprintf( "%s/%s/%s/%s", g.source, fol, name, "thumb.jpg" ) )
	if err != nil {
		return err
	}

	defer t2.Close()

	t3, err := os.Create( fmt.Sprintf( "%s/%s/%s/%s", g.dest, fol, name, "thumb.jpg" ) )
	if err != nil {
		return err
	}

	defer t3.Close()

	i1, err := os.Create( fmt.Sprintf( "%s/%s/%s/%s", g.man, fol, name, "image.jpg" ) )
	if err != nil {
		return err
	}

	defer i1.Close()

	i2, err := os.Create( fmt.Sprintf( "%s/%s/%s/%s", g.source, fol, name, "image.jpg" ) )
	if err != nil {
		return err
	}

	defer i2.Close()

	i3, err := os.Create( fmt.Sprintf( "%s/%s/%s/%s", g.dest, fol, name, "image.jpg" ) )
	if err != nil {
		return err
	}

	defer i3.Close()

	writer := io.MultiWriter( t1, t2, t3 )
	err = jpeg.Encode( writer, thumb, nil )
	if err != nil {

		return err

	}

	writer = io.MultiWriter( i1, i2, i3 )
	err = jpeg.Encode( writer, img, nil )
	if err != nil {

		return err

	}

	file, err = header.Open()
	if err != nil {

		return err

	}

	buf, err := ioutil.ReadAll(file)
	if err != nil {

		return err

	}

	err = ioutil.WriteFile(fmt.Sprintf( "%s/%s/%s/%s", g.man, fol, name, "orig.jpg" ), buf, 0777)
	if err != nil {

		return err

	}


	return nil

	return nil

}

func ( g *gallery ) AddImage( header *multipart.FileHeader, which, text string, x, y, w, h int64 ) ( error ) {

	file, err := header.Open()
	if err != nil {

		return err

	}

	im, err := jpeg.Decode( file )
	if err != nil {

		return err

	}

	thumb := g.createThumb( im, x, y, w, h )

	ww := im.Bounds().Max.X - im.Bounds().Min.X
	hh := im.Bounds().Max.Y - im.Bounds().Min.Y

	var img image.Image

	if ww > hh {

		// landscape
		img = resize.Resize( 1200, 0, im, resize.Lanczos3 )


	} else {

		// portait or rectangle
		img = resize.Resize( 0, 800, im, resize.Lanczos3 )

	}

	data, err := g.readData()
	if err != nil {
		return err
	}

	var name string
	var fol string
	var nameInt int

	if which == FOLDER_ONE {

		fol = FOLDER_ONE
		name, nameInt = g.getName( data.One )
		data.One = append( data.One, images{ name, nameInt, text, x, y, w, h } )

	} else if which == FOLDER_TWO {

		fol = FOLDER_TWO
		name, nameInt = g.getName( data.Two )
		data.Two = append( data.Two, images{ name, nameInt, text, x, y, w, h } )

	} else if which == FOLDER_THREE {

		fol = FOLDER_THREE
		name, nameInt = g.getName( data.Three )
		data.Three = append( data.Three, images{ name, nameInt, text, x, y, w, h } )

	}

	err = g.writeData( data )
	if err != nil {

		return err

	}

	// write files
	err = g.createFolders( fmt.Sprintf( "%s/%s/%s", g.man, fol, name ), fmt.Sprintf( "%s/%s/%s", g.source, fol, name ), fmt.Sprintf( "%s/%s/%s", g.dest, fol, name ) )
	if err != nil {

		return err

	}

	t1, err := os.Create( fmt.Sprintf( "%s/%s/%s/%s", g.man, fol, name, "thumb.jpg" ) )
	if err != nil {
		return err
	}

	defer t1.Close()

	t2, err := os.Create( fmt.Sprintf( "%s/%s/%s/%s", g.source, fol, name, "thumb.jpg" ) )
	if err != nil {
		return err
	}

	defer t2.Close()

	t3, err := os.Create( fmt.Sprintf( "%s/%s/%s/%s", g.dest, fol, name, "thumb.jpg" ) )
	if err != nil {
		return err
	}

	defer t3.Close()

	i1, err := os.Create( fmt.Sprintf( "%s/%s/%s/%s", g.man, fol, name, "image.jpg" ) )
	if err != nil {
		return err
	}

	defer i1.Close()

	i2, err := os.Create( fmt.Sprintf( "%s/%s/%s/%s", g.source, fol, name, "image.jpg" ) )
	if err != nil {
		return err
	}

	defer i2.Close()

	i3, err := os.Create( fmt.Sprintf( "%s/%s/%s/%s", g.dest, fol, name, "image.jpg" ) )
	if err != nil {
		return err
	}

	defer i3.Close()

	writer := io.MultiWriter( t1, t2, t3 )
	err = jpeg.Encode( writer, thumb, nil )
	if err != nil {

		return err

	}

	writer = io.MultiWriter( i1, i2, i3 )
	err = jpeg.Encode( writer, img, nil )
	if err != nil {

		return err

	}

	file, err = header.Open()
	if err != nil {

		return err

	}

	buf, err := ioutil.ReadAll(file)
	if err != nil {

		return err

	}

	err = ioutil.WriteFile(fmt.Sprintf( "%s/%s/%s/%s", g.man, fol, name, "orig.jpg" ), buf, 0777)
	if err != nil {

		return err

	}


	return nil

}

func ( g *gallery ) getName( data []images ) ( string, int ) {

	if len( data ) == 0 {

		return "000", 0

	}

	var dataMap map[int]string = make( map[int]string )
	for _, img := range data {

		dataMap[img.Nm] = img.Id

	}

	for i := 0; i < 1000; i++ {

		if _, ok := dataMap[i]; !ok {

			return fmt.Sprintf( "%03d", i ), i

		}

	}

	log.Fatal( "gallery :: getName :: Fatal Error -> data fucker..." )

	return "", 0

}

func ( g *gallery ) createThumb( im image.Image, x, y, w, h int64 ) ( image.Image ) {

	imageBounds := im.Bounds().Max
	thumbWidth := int( ( 200.0 / float64( w ) ) * float64( imageBounds.X ) )
	thumbHeight := int( ( 170.0 / float64( h ) ) * float64( imageBounds.Y ) )
	thumbX := int( ( float64( x ) - 188.0 ) * ( float64( imageBounds.X ) /  float64( w ) ) * -1.0 )
	thumbY := int( ( float64( y ) - 151.0 ) * ( float64( imageBounds.Y ) /  float64( h ) ) * -1.0 )

	log.Printf( "Hello from thumbs creator. %d, %d, %d, %d \n", thumbX, thumbY, thumbWidth, thumbHeight )

	if ( thumbWidth + thumbX ) > imageBounds.X {
		thumbX = imageBounds.X - thumbWidth
	}

	if ( thumbHeight + thumbY ) > imageBounds.Y {
		thumbY = imageBounds.Y - thumbHeight
	}

	croppedImage := image.NewRGBA( image.Rect( 0, 0, thumbWidth, thumbHeight ) )
	draw.Draw( croppedImage, croppedImage.Bounds(), im, image.Point{ thumbX, thumbY }, draw.Src )

	newImg := resize.Resize( uint( 200 ), 0, croppedImage, resize.Lanczos3)
	return newImg

}

func ( g *gallery ) resetData() ( error ) {

	// remove all folders
	var err error

	err = os.RemoveAll( g.source )
	if err != nil {
		return err
	}

	err = os.RemoveAll( g.man )
	if err != nil {
		return err
	}

	// recreate them :: empty
	err = g.createFolders(
			fmt.Sprintf( "%s/%s", g.man, FOLDER_ONE ),
			fmt.Sprintf( "%s/%s", g.man, FOLDER_TWO ),
			fmt.Sprintf( "%s/%s", g.man, FOLDER_THREE ),
			fmt.Sprintf( "%s/%s", g.source, FOLDER_ONE ),
			fmt.Sprintf( "%s/%s", g.source, FOLDER_TWO ),
			fmt.Sprintf( "%s/%s", g.source, FOLDER_THREE ) )
	if err != nil {
		return err
	}

	// write new and emtpy data file
	var d galleryData = galleryData{ One:make( []images, 0 ), Two:make( []images, 0 ), Three:make( []images, 0 ) }
	err = g.writeData( d )
	if err != nil{
		return err
	}
	
	return nil

}

func (g *gallery) writeData(d galleryData) error {

	value, err := json.Marshal(d)
	if err != nil {

		return err

	}

	err = ioutil.WriteFile(g.path, value, 0777)
	if err != nil {

		return err

	}

	return nil

}

func (g *gallery) readData() (galleryData, error) {

	var d galleryData

	data, err := g.GetDataJson()
	if err != nil {

		return d, err

	}

	err = json.Unmarshal(data, &d)
	if err != nil {

		return d, err

	}

	return d, nil

}

func (g *gallery) GetDataJson() ([]byte, error) {

	var data []byte
	var err error

	data, err = ioutil.ReadFile(g.path)
	if err != nil {

		return data, err

	}

	return data, nil

}

func ( g *gallery ) createFolders( folders ...string ) ( error ) {

	for _, folder := range folders {

		err := os.MkdirAll( folder, 0777 )
		if err != nil {

			return err

		}

	}

	return nil

}

func ( g *gallery ) checkFiles( folders ...string ) ( error ) {

	for _, folder := range folders {

		_, err := os.Stat( folder )
		if err != nil {

			if os.IsNotExist( err ) {

				log.Fatal( fmt.Sprintf( "newGallery :: check dirs :: %s must exist...", folder ) )

			} else {

				return err

			}

		}

	}

	return nil

}