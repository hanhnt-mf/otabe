package pbl_otabe

import (
	"math"
	"reflect"
)

func ItemExists(arrayType interface{}, item interface{}) bool {
	arr := reflect.ValueOf(arrayType)

	if arr.Kind() != reflect.Array && arr.Kind() != reflect.Slice {
		panic("Invalid data-type")
	}

	for i := 0; i < arr.Len(); i++ {
		if arr.Index(i).Interface() == item {
			return true
		}
	}

	return false
}

// haversin(θ) function
func Hsin(theta float64) float64 {
	return math.Pow(math.Sin(theta/2), 2)
}

// distance returned is METERS!!!!!!
// http://en.wikipedia.org/wiki/Haversine_formula
func Distance(lat1, lon1, lat2, lon2 float64) float64 {
	// convert to radians
	// must cast radius as float to multiply later
	var la1, lo1, la2, lo2, r float64
	la1 = lat1 * math.Pi / 180
	lo1 = lon1 * math.Pi / 180
	la2 = lat2 * math.Pi / 180
	lo2 = lon2 * math.Pi / 180

	r = 6378100 // Earth radius in METERS

	// calculate
	h := Hsin(la2-la1) + math.Cos(la1)*math.Cos(la2)*Hsin(lo2-lo1)

	return 2 * r * math.Asin(math.Sqrt(h))
}