import Pool from "pg-pool";
import { writeFile } from "fs/promises";

const pgConfig = {
  user: "socra",
  password: "sigl2025",
  host: "localhost",
  port: 5432,
  database: "sotracteur",
};

async function query(sql) {
  var pool = new Pool(pgConfig);
  var client = await pool.connect();
  let result;
  try {
    result = await client.query(sql);
  } catch (e) {
    console.error(e.message, e.stack);
  } finally {
    client.release();
  }
  return result;
}

async function readDataFromPostgres() {
  console.log("Reading data from Postgres");
  const results = await query(
    "SELECT * FROM public.tracteurs_disponibles LIMIT 2000"
  );
  try {
    const content = JSON.stringify(results.rows, null, 2);
    await writeFile("./tractors.json", content);
  } catch (err) {
    console.log(err);
  }
}

readDataFromPostgres();
