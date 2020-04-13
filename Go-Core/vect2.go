package main

import "math"

type Vect2 struct {
	X float64 `json:"x"`
	Y float64 `json:"y"`
}

func copy(vect Vect2) Vect2 {
	return Vect2{
		X: vect.X,
		Y: vect.Y,
	}
}

func dot(vect1, vect2 Vect2) float64 {
	return vect1.X*vect2.X + vect1.Y*vect2.Y
}

func norm(vect Vect2) float64 {
	return math.Sqrt(dot(vect, vect))
}

func normalize(vect Vect2) Vect2 {
	return Vect2{
		X: vect.X / norm(vect),
		Y: vect.Y / norm(vect),
	}
}

func scalar_times(vect Vect2, number float64) Vect2 {
	return Vect2{
		X: vect.X * number,
		Y: vect.Y * number,
	}
}

func add(vect1, vect2 Vect2) Vect2 {
	return Vect2{
		X: vect1.X + vect2.X,
		Y: vect1.Y + vect2.Y,
	}
}

func substract(vect1, vect2 Vect2) Vect2 {
	return Vect2{
		X: vect1.X - vect2.X,
		Y: vect1.Y - vect2.Y,
	}
}
