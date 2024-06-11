/**
 *
 *  Haversine: Compute the shortest distance over the earth’s surface between two positions
 *
 *  formula:	a = sin²(Δφ/2) + cos φ1 ⋅ cos φ2 ⋅ sin²(Δλ/2)
 *              c = 2 ⋅ atan2( √a, √(1−a) )
 *              d = R ⋅ c
 *  where:
 *              φ is latitude, λ is longitude, R is earth’s radius (mean radius = 6,371km);
 *
 * See. https://www.movable-type.co.uk/scripts/latlong.html
 *
 * @returns distance in meters
 */
function haversine(
  sourceLatitude,
  sourceLongitude,
  destLatitude,
  destLongitude
) {
  const R = 6371e3; // Earth radius in meters

  // convert to radians
  const φ1 = (sourceLatitude * Math.PI) / 180;
  const φ2 = (destLatitude * Math.PI) / 180;
  const Δφ = ((destLatitude - sourceLatitude) * Math.PI) / 180;
  const Δλ = ((destLongitude - sourceLongitude) * Math.PI) / 180;

  const a =
    Math.sin(Δφ / 2) * Math.sin(Δφ / 2) +
    Math.cos(φ1) * Math.cos(φ2) * Math.sin(Δλ / 2) * Math.sin(Δλ / 2);
  const c = 2 * Math.atan2(Math.sqrt(a), Math.sqrt(1 - a));
  const d = R * c; // in meters
  return d;
}

/**
 * Computes distance in kilometers from a source position
 * to a destination position.
 *
 * Position is an object like:
 * { latitude: number, longitude: number }
 *
 * @param {*} startPos Position of the source geolocation
 * @param {*} destPos Position of the destination geolocation
 * @returns distance in kilometers between source position and destination position
 */
export function distance(startPos, destPos) {
  const distanceInMeters = haversine(
    startPos.lat,
    startPos.lng,
    destPos.x,
    destPos.y
  );
  return Math.round(distanceInMeters / 1000);
}
