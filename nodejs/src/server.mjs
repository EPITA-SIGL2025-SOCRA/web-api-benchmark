import express from "express";
import { distance } from "./distance.mjs";
import { env, exit } from "process";

const DATA_JSON_FILE_PATH = env["DATA_JSON_FILE_PATH"] || exit(1);
import { createRequire } from "module";
const require = createRequire(import.meta.url);
const TRACTORS = require(DATA_JSON_FILE_PATH);

const app = express();
const port = 3000;

app.use(express.json());

app.get("/", function (request, response) {
  response.send({ Hello: "Sotracteur!" });
});

/**
 * Find all tractors that are <radius>-km for the position of the user
 */
app.post("/v1/tractors", function (request, response) {
  try {
    const radius = +request.query.radius;
    const userPosition = request.body.position;
    const tractorsAroundUsers = TRACTORS.map(function (tractor) {
      const distanceToProduct = distance(userPosition, tractor.lat_long);
      return {
        ...tractor,
        distance: distanceToProduct,
      };
    }).filter(function ({ distance }) {
      return distance < radius;
    });
    console.log("NODEJS: tractors sent to user with position", userPosition);
    response.send(tractorsAroundUsers);
  } catch (e) {
    console.log(`ERROR with /v1/tractors: ${e}`);
    response.status = 500;
    response.send({ details: e });
  }
});

app.listen(port, function () {
  console.log(`Server started on localhost:${port}`);
});
