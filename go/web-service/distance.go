package main

import (
	"math"
)

/*
Haversine: Compute the shortest distance over the earth’s surface between two positions

formula:	a = sin²(Δφ/2) + cos φ1 ⋅ cos φ2 ⋅ sin²(Δλ/2)

	c = 2 ⋅ atan2( √a, √(1−a) )
	d = R ⋅ c

where:

	φ is latitude, λ is longitude, R is earth’s radius (mean radius = 6,371km);

See. https://www.movable-type.co.uk/scripts/latlong.html
*/
func haversine(srcLat float64, srcLng float64, dstLat float64, dstLng float64) float64 {
	const earthRadius = 6371e3
	const pi = math.Pi

	// convert to radians
	φ1 := (srcLat * pi) / 180
	φ2 := (dstLat * pi) / 180
	Δφ := ((dstLat - srcLat) * pi) / 180
	Δλ := ((dstLng - srcLng) * pi) / 180

	a := math.Sin(Δφ/2)*math.Sin(Δφ/2) + math.Cos(φ1)*math.Cos(φ2)*math.Sin(Δλ/2)*math.Sin(Δλ/2)
	c := 2 * math.Atan2(math.Sqrt(a), math.Sqrt(1-a))
	d := earthRadius * c
	return d
}

/*
		Computes distance in kilometers from a source position to a destination position.
	    Position is an object like:
	        { lat: float, lng: float }

	    @param {*} start_pos Position of the source geolocation
	    @param {*} dest_pos Position of the destination geolocation
	    @returns distance in kilometers between source position and destination position
	    """
*/
func distance(startPos position, destPos latlong) int {
	distanceInMeters := haversine(startPos.Lat, startPos.Lng, destPos.Lat, destPos.Lng)
	return int(math.Round(distanceInMeters / 1000))
}

func (t *tractor) checkDistance(userPosition position, radius int) *tractorWithDistance {

	productDistanceFromUser := distance(userPosition, t.Position)

	if productDistanceFromUser < radius {
		return &tractorWithDistance{
			Title:            t.Title,
			Modele:           t.Modele,
			Position:         t.Position,
			NomCommune:       t.NomCommune,
			Categorie:        t.Categorie,
			ImageUrl:         t.ImageUrl,
			AnneeFabrication: t.AnneeFabrication,
			Distance:         productDistanceFromUser,
		}
	}
	return nil
}
