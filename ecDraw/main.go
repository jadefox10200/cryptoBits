//this program will draw the elliptic curve to a png file "image.png"
//the curve parameters are: 
// y^2 = x^ + 3x + 3

package main

import (	
	"image"
	"image/color"
	"image/png"
	"os"
	"log"
	"math"
	"fmt"
)

type ECC struct {
	X	[]int
	Y 	[]int
}


func main() {

	const width, height = 1000, 1000

	// Create a colored image of the given width and height.
	img := image.NewNRGBA(image.Rect(0, 0, width, height))

	var x float64
	var y float64

	//arbitrary start point:
	x = 3.0

	//draw the curve: 
	for {

		//formula inputs:
		// y^2 = x^3 - 3x + 3
		x3 := x * x * x
		threex := 2.0 * x
		three := 2.0

		y2 := x3 - threex + three 

		//when y^2 hits 0 we're crossing over the center point. we cannot find the sqrt of a negative number
		if y2 < 0 {break}

		y = math.Sqrt(y2) 
		
		//move the x value slightly to create the progression over real numbers:
		x = x - 0.0001

		//multiply by 100 to get an integer rather than a decimal and cast it.
		xi := int(x*100)
		yi := int(y*100)

		//add half the value of the image size to get the center point. 
		xi += 500
		yi += 500

		//bottom
		img.Set(xi, yi, color.NRGBA{
	        R: uint8(0),
	        G: uint8(0),
	        B: uint8(255),
	        A: 255,
	    })	 

	    

		//in order to draw the top, we subtract y from the center point. 
		yt := int(y*100)
		yi = 500 - yt

		//top
		img.Set(xi, yi, color.NRGBA{
		    R: uint8(0),
	        G: uint8(0),
	        B: uint8(255),
	        A: 255,
		})

		//print the values to the console: 
		fmt.Printf("x:%v\ty:%v\n", xi, yi)

	}

	//draw y axis:
	dx := int(width/2)
	for dy := 0; dy < height; dy++{
		img.Set(dx-1, dy, color.NRGBA{
			R: uint8(50),
	        G: uint8(50),
	        B: uint8(50),
	        A: 255,
		})

		img.Set(dx, dy, color.NRGBA{
			R: uint8(150),
	        G: uint8(150),
	        B: uint8(150),
	        A: 255,
		})

		img.Set(dx+1, dy, color.NRGBA{
			R: uint8(150),
	        G: uint8(150),
	        B: uint8(150),
	        A: 255,
		})
	}

	//draw x axis:
	dy := int(height/2)
	for dx = 0; dx < width; dx++{
		img.Set(dx, dy-1, color.NRGBA{
			R: uint8(50),
	        G: uint8(50),
	        B: uint8(50),
	        A: 255,
		})

		img.Set(dx, dy, color.NRGBA{
			R: uint8(50),
	        G: uint8(50),
	        B: uint8(50),
	        A: 255,
		})

		img.Set(dx, dy+1, color.NRGBA{
			R: uint8(50),
	        G: uint8(50),
	        B: uint8(50),
	        A: 255,
		})
	}


	f, err := os.Create("image.png")
	if err != nil {
	    log.Fatal(err)
	}

	if err := png.Encode(f, img); err != nil {
	    f.Close()
	    log.Fatal(err)
	}

	if err := f.Close(); err != nil {
	    log.Fatal(err)
	}

	return
}

	
