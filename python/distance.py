from math import atan2, cos, sin, sqrt, pi


def haversine(src_lat, src_lng, dst_lat, dst_lng):
    """
    Haversine: Compute the shortest distance over the earth’s surface between two positions
       formula:	a = sin²(Δφ/2) + cos φ1 ⋅ cos φ2 ⋅ sin²(Δλ/2)
                   c = 2 ⋅ atan2( √a, √(1−a) )
                   d = R ⋅ c
       where:
                   φ is latitude, λ is longitude, R is earth’s radius (mean radius = 6,371km);
       See. https://www.movable-type.co.uk/scripts/latlong.html
    """
    R = 6371e3  # Earth radius in meters
    # convert to radians
    φ1 = (src_lat * pi) / 180
    φ2 = (dst_lat * pi) / 180
    Δφ = ((dst_lat - src_lat) * pi) / 180
    Δλ = ((dst_lng - src_lng) * pi) / 180

    a = sin(Δφ / 2) * sin(Δφ / 2) + cos(φ1) * cos(φ2) * sin(Δλ / 2) * sin(Δλ / 2)
    c = 2 * atan2(sqrt(a), sqrt(1 - a))
    d = R * c  # distance d in meters
    return d


def distance(start_pos, dest_pos):
    """
    Computes distance in kilometers from a source position to a destination position.
    Position is an object like:
        { lat: float, lng: float }

    @param {*} start_pos Position of the source geolocation
    @param {*} dest_pos Position of the destination geolocation
    @returns distance in kilometers between source position and destination position
    """
    distance_in_meters = haversine(
        start_pos["lat"], start_pos["lng"], dest_pos["x"], dest_pos["y"]
    )
    return round(distance_in_meters / 1000)
