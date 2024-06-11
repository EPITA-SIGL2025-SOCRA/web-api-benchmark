from typing import Tuple, Union

from fastapi import FastAPI

from pydantic import BaseModel

from distance import distance

from json import load
from os import environ


class Position(BaseModel):
    lat: float
    lng: float


class PositionInput(BaseModel):
    position: Position


app = FastAPI()


def getRequiredEnvVar(env_key: str):
    env_value = environ.get(env_key)
    if env_value is None:
        raise Exception(f"missing required env var {env_key}")
    return env_value


def load_tractors():
    DATA_JSON_FILE_PATH = getRequiredEnvVar("DATA_JSON_FILE_PATH")
    tractors = []
    with open(DATA_JSON_FILE_PATH, "r") as tractors_json_file:
        tractors = load(tractors_json_file)
    return tractors


# loads all tractors **in memory** when web-service is starting
TRACTORS = list(load_tractors())


@app.get("/")
def read_root():
    return {"Hello": "Sotracteur"}


@app.post("/v1/tractors")
async def read_tractors(position_input: PositionInput, radius: Union[str, None] = None):
    user_position = dict(position_input.position)
    radius_int = int(radius)
    tractors_near_user = []

    def compute_distance(tractor):
        tractor_distance_from_user = distance(
            start_pos=user_position, dest_pos=tractor["lat_long"]
        )
        return {**tractor, "distance": tractor_distance_from_user}

    def is_inside_radius(tractor_with_distance):
        user_distance_from_tractor = tractor_with_distance["distance"]
        return user_distance_from_tractor < radius_int

    tractors_near_user = filter(is_inside_radius, map(compute_distance, TRACTORS))
    print(f"PYTHON: products sent to user with position {user_position}")
    return list(tractors_near_user)
