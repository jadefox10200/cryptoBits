package ec

import (
	"fmt"	
	"math/big"
)

var big2 = big.NewInt(2)
var big3 = big.NewInt(3)

//y^2 = x^3 + a*x + b mod p
var a = big.NewInt(2)
var b = big.NewInt(2)
var p = big.NewInt(17)

type Point struct {
	X	*big.Int
	Y	*big.Int
}

func printPointOrders(p *big.Int) {

	pts := getPoints(p)

	for _, v := range pts {
		i := ord(v.X, v.Y, p)
		fmt.Printf("ord(%02d,%02d) = %02d\n", v.X, v.Y, i)
	}

	return 
}

func scalarMult(x, y, d, p *big.Int) {

	x1, y1 := big.NewInt(x.Int64()), big.NewInt(y.Int64())
	t := d.BitLen()		
	// fmt.Println(x1, y1)
	for i := t-2; i >= 0; i-- {
		x1, y1 = doubleP(x1, y1, p)
		// fmt.Printf("Double i:%v \tx:%v\ty:%v\n", i, x1, y1)
		if (d.Uint64() >> uint(i) & 1) == 1 {
			x1, y1 = addP(x, y, x1, y1, p)
			// fmt.Printf("Add i:%v \tx:%v\ty:%v\n", i, x1, y1)
		}
				
	}

	fmt.Printf("%vP\tx:%v\ty:%v \n", d, x1, y1)

	return
}

func ord(x, y, p *big.Int) int {
		
	//start at 1 to count the generator of the group
	var counter int = 1

	//add again for the first double
	counter++
	x1, y1 := doubleP(x, y, p)

	for {
		//add for each iteration
		counter++
		x1, y1 = addP(x, y, x1, y1, p)	
		//check if we got the inverse of P. if so, we've hit the end of the cyclic group:
		//add 1 to counter before break to include the O element. 
		if x1.Cmp(x) == 0 && new(big.Int).Sub(p, y1).Cmp(y) == 0 {counter++; break}
	}
	
	return counter

}

func printGroup(x, y, p *big.Int) {

	fmt.Println(x, y)
	//start at 1 to count the generator of the group
	var counter int = 1

	//add again for the first double
	counter++
	x1, y1 := doubleP(x, y, p)
	fmt.Printf("%02dP:(%02d,%02d) : %v\n", counter,  x1, y1, isOnCurve(x1, y1, p))

	for {
		//add for each iteration
		counter++
		x1, y1 = addP(x, y, x1, y1, p)
		fmt.Printf("%02dP:(%02d,%02d) : %v\n", counter, x1, y1, isOnCurve(x1, y1, p))
		//check if we got the inverse of P. if so, we've hit the end of the cyclic group:
		//add 1 to counter before break to include the O element. 
		if x1.Cmp(x) == 0 && new(big.Int).Sub(p, y1).Cmp(y) == 0 {counter++; break}
	}

	fmt.Println("Order:", counter)

	return
}

func getPoints(p *big.Int) []Point {

	points := make([]Point, 0, 0)

	counter := int(p.Int64())
	for xi := 0; xi < counter; xi++ {
		for yi := 0; yi < counter; yi++ {
			x := big.NewInt(int64(xi))
			x3 := new(big.Int).Mul(x, x)
	  		x3.Mul(x3, x)
	  		twoX := new(big.Int).Lsh(x, 1)
	  		// threeX.Add(threeX, x)
	  	
		  	x3.Add(x3, twoX)
		  	x3.Add(x3, b)
		  	x3.Mod(x3, p)

		  	y := big.NewInt(int64(yi))
		  	y2 := new(big.Int).Mul(y, y)  	
	  		y2.Mod(y2, p)

		  	if x3.Cmp(y2) == 0 {
		  		// fmt.Println(xi, yi, isOnCurve(x,y,p)) 
		  		pt := Point{
		  			X: x,
		  			Y: y,
		  		}
		  		points = append(points, pt)

		  	}
		} 
	}

	return points
}

func isOnCurve(x, y, p *big.Int) bool {

	// y² = x³ + a*x + b
  	y2 := new(big.Int).Mul(y, y)  	
  	y2.Mod(y2, p)

  	x3 := new(big.Int).Mul(x, x)
  	x3.Mul(x3, x)
  	
  	twoX := new(big.Int).Lsh(x, 1)
  	// threeX.Add(threeX, x)
  	
  	x3.Add(x3, twoX)
  	x3.Add(x3, b)
  	x3.Mod(x3, p)

  	return x3.Cmp(y2) == 0
}

func addP(x1, y1, x2, y2, p *big.Int) (*big.Int, *big.Int) {

	s := big.NewInt(0)

	//if P = P then we are actuallying doing a doubling and need to use the correct S value for this:
	if x1.Cmp(x2) == 0 && y1.Cmp(y2) == 0 {
		sy := new(big.Int).ModInverse(new(big.Int).Mul(big2, y1) , p) 
		sx := new(big.Int).Add(new(big.Int).Mul(big3, new(big.Int).Exp(x1, big2, nil) ), big2)

		s = new(big.Int).Mod(new(big.Int).Mul(sy, sx), p)
	} else {
		sy := new(big.Int).Sub(y2, y1)
		sx := new(big.Int).ModInverse(new(big.Int).Sub(x2, x1), p)
		s = new(big.Int).Mod(new(big.Int).Mul(sy, sx), p)	
	}

	x3 := new(big.Int).Mod(new(big.Int).Sub(new(big.Int).Sub(new(big.Int).Exp(s, big2, nil), x1), x2), p)

	y3 := new(big.Int).Mod(new(big.Int).Sub(new(big.Int).Mul(s, new(big.Int).Sub(x1, x3)), y1), p)

	return x3, y3
}

func doubleP(x1, y1, p *big.Int) (*big.Int, *big.Int) {

	sy := new(big.Int).ModInverse(new(big.Int).Mul(big2, y1) , p) 
	sx := new(big.Int).Add(new(big.Int).Mul(big3, new(big.Int).Exp(x1, big2, nil) ), big2)

	s := new(big.Int).Mod(new(big.Int).Mul(sy, sx), p)

	//the second x1 should be x2 but since we are doing point doubling, x1 and x2 are the same: 
	x3 := new(big.Int).Mod(new(big.Int).Sub(new(big.Int).Sub(new(big.Int).Exp(s, big2, nil), x1), x1), p)

	y3 := new(big.Int).Mod(new(big.Int).Sub(new(big.Int).Mul(s, new(big.Int).Sub(x1, x3)), y1), p)

	return x3, y3

}
