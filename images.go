package main

import (
	"image"
	"io"
	"os"
	// "log"
	"fmt"
	"image/jpeg"
	"image/draw"
	"github.com/nfnt/resize"
)

type imager struct {
	man string
	source string
	dest string
}

func ( i *imager ) SetJumboImage( im image.Image ) ( error ) {

	// write images
	err := i.writeImage( "/images/jumbo-bg.jpg", im, 16.0, 9.0, 1400 )
	if err != nil {

		return err

	}

	return nil

}


func ( i *imager ) SetHeaderOneImage( im image.Image ) ( error ) {

	// write images
	err := i.writeImage( "/images/title-header-one.jpg", im, 3.0, 1.0, 1200 )
	if err != nil {

		return err

	}

	return nil

}

func ( i *imager ) SetHeaderTwoImage( im image.Image ) ( error ) {

	// write images
	err := i.writeImage( "/images/title-header-two.jpg", im, 3.0, 1.0, 1200 )
	if err != nil {

		return err

	}

	return nil

}

func ( i *imager ) writeImage( path string, im image.Image, ratioX, ratioY float64, width uint ) ( error ) {

	var sb image.Rectangle = im.Bounds()
	var sw, sh, cw, ch, cx, cy float64

	sw = float64( sb.Max.X - sb.Min.X )
	sh = float64( sb.Max.Y - sb.Min.Y )

	if ( sw / ratioX ) >= ( sh / ratioY ) {

		ch = sh
		cw = ( ch / ratioY ) * ratioX
		cy = 0.0
		cx = ( sw - cw ) / 2.0

	} else {


		// hoger
		cw = sw
		ch = ( sw / ratioX ) * ratioY
		cx = 0.0
		cy = ( sh - ch ) / 2.0

	}

	// log.Println(  fmt.Sprintf( "cx : %f, cy : %f, cw : %f, ch : %f ", cx, cy, cw, ch ) )

	dst := image.NewRGBA( image.Rect( 0, 0, int( cw ), int( ch ) ) )
	draw.Draw( dst, dst.Bounds(), im, image.Point{ int( cx ), int( cy ) } , draw.Src )

	rslt := resize.Resize( width, 0, dst, resize.Lanczos3 )

	f1, err := os.Create( fmt.Sprintf( "%s%s", i.man, path ) )
	if err != nil {
		return err
	}
	defer f1.Close()

	f2, err := os.Create( fmt.Sprintf( "%s%s", i.source, path ) )
	if err != nil {
		return err
	}
	defer f2.Close()

	f3, err := os.Create( fmt.Sprintf( "%s%s", i.dest, path ) )
	if err != nil {
		return err
	}
	defer f3.Close()

	err = jpeg.Encode( io.MultiWriter( f1, f2, f3 ), rslt, nil )
	if err != nil {

		return err

	}

	return nil

}

func newImager( man, source, dest string ) ( *imager ) {

	i := &imager{ man, source, dest }
	return i

}