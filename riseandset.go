package riseandset

import (
	"math"
)

const (
	DegToRad     = (2.0 * math.Pi) / 360.0
	Jan1200012pm = 2451545.0
)

// This will generate a sunrise and sunset time for a given date,
// longitude, latitude, and altitude.
//
// Implementation based on https://en.wikipedia.org/wiki/Sunrise_equation
//
// Date is a julian day.
// Longitude and latitude should be specified in degrees north and west.
// The altitude should be supplied in meters.
// The sunrise and sunset times are returned as julian day values.
func Times(date int, longitude, latitude, altitude float64) (float64, float64) {
	jd := (float64(date) - Jan1200012pm) + 0.0008
	meanSolNoon := meanSolarNoon(jd, longitude)
	meanSolAnomaly := meanSolarAnomaly(meanSolNoon)
	eclipticLongitude := eclipticLongitude(meanSolAnomaly)
	hourAngle := hourAngle(altitude, latitude, eclipticLongitude)
	return julRiseAndSet(
		meanSolNoon,
		meanSolAnomaly,
		eclipticLongitude,
		hourAngle)
}

func meanSolarNoon(julDay, longitude float64) float64 {
	return (longitude / 360.0) + julDay
}

func meanSolarAnomaly(meanSolarNoon float64) float64 {
	return math.Mod(
		357.5291+0.98560028*meanSolarNoon,
		360.0)
}

func center(meanSolAnomaly float64) float64 {
	meanSolAnomalyRadians := meanSolAnomaly * DegToRad

	return 1.9148*math.Sin(meanSolAnomalyRadians) +
		0.02*math.Sin(2.0*meanSolAnomalyRadians) +
		0.0003*math.Sin(3.0*meanSolAnomalyRadians)
}

func eclipticLongitude(meanSolAnomaly float64) float64 {
	return math.Mod(
		meanSolAnomaly+
			center(meanSolAnomaly)+
			180.0+
			102.9372,
		360.0)
}

func hourAngle(altitude, latitude, eclipticLongitude float64) float64 {
	altitudeModifier := -2.076 * (math.Sqrt(altitude) / 60.0)
	latitudeRadians := latitude * DegToRad

	solDeclination :=
		math.Asin(
			math.Sin(eclipticLongitude*DegToRad)*
				math.Sin(23.43713*DegToRad)) / DegToRad

	solDeclinationRadians := solDeclination * DegToRad

	return math.Acos(
		(math.Sin((-0.83+altitudeModifier)*DegToRad)-
			math.Sin(latitudeRadians)*
				math.Sin(solDeclinationRadians))/
			(math.Cos(latitudeRadians)*
				math.Cos(solDeclinationRadians))) / DegToRad
}

func julRiseAndSet(
	meanSolNoon,
	meanSolAnomaly,
	eclipticLongitude,
	hourAngle float64) (float64, float64) {
	solTransit :=
		meanSolNoon +
			(0.0053 * math.Sin(meanSolAnomaly*DegToRad)) -
			(0.0069 * math.Sin(2.0*eclipticLongitude*DegToRad))

	hourAngUnit := hourAngle / 360.0

	riseJulian := Jan1200012pm + (solTransit - hourAngUnit)
	setJulian := Jan1200012pm + (solTransit + hourAngUnit)
	return riseJulian, setJulian
}
