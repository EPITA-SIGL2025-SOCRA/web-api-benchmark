package loadtest


import io.gatling.core.Predef._
import io.gatling.http.Predef._

import scala.concurrent.duration._
import scala.math._
import scala.util.Random

case class EnvironmentVariableException(private val message: String = "",
                                        private val cause: Throwable = None.orNull) extends Exception(message, cause)

class GetTractorsSimulation extends Simulation {

  @throws(classOf[EnvironmentVariableException])
  def getEnvOrThrow(name: String): String = {
    val env = System.getenv(name)
    if (env != null) env else throw EnvironmentVariableException(s"missing env var ${name}")
  }
  private val backendUrl: String = getEnvOrThrow("LOADTEST_BACKEND_URL")
  private val constantConcurrentUsersNumber: String = getEnvOrThrow("LOADTEST_CONSTANT_CONCURRENT_USERS")
  private val loadtestDuration: String = getEnvOrThrow("LOADTEST_DURATION")

  private val httpProtocol = http
    .baseUrl(backendUrl)
    .inferHtmlResources()
    .acceptHeader("*/*")
    .acceptEncodingHeader("gzip, deflate")
    .acceptLanguageHeader("en-US,en;q=0.5")
    .userAgentHeader("Mozilla/5.0 (X11; Linux x86_64; rv:104.0) Gecko/20100101 Firefox/104.0")

  /**
   * pick a random geo position (lat, lon) within a RANDOM_POSITION_RADIUS
   * from the REFERENCE_POSITION.
   *
   *  Returns a tuple (lat, lon)
   *  See. https://gis.stackexchange.com/questions/25877/generating-random-locations-nearby
   *  @param ref_position
   *  @return
   */
  def getRandomLocation(ref_position: Map[String, Double]) : (Double, Double) = {
    val rand = new Random()

    val RANDOM_POSITION_RADIUS = 100e3
    val radius_in_degrees = RANDOM_POSITION_RADIUS / 111000
    val x0 = ref_position("lng")
    val y0 = ref_position("lat")

    val u = rand.nextDouble()
    val v = rand.nextDouble()

    val w = radius_in_degrees * sqrt (u)
    val t = 2 * Pi * v
    val x = w * cos (t)
    val y = w * sin (t)

    val new_x = x / cos (toRadians(y0))

    val lng = BigDecimal(new_x + x0).setScale(6, BigDecimal.RoundingMode.HALF_UP).toDouble
    val lat = BigDecimal(y + y0).setScale(6, BigDecimal.RoundingMode.HALF_UP).toDouble

    (lat, lng)
  }

  def nextRandDouble(base: Double): Double = {
    val rand = new Random()
    val pair = rand.nextInt(2)
    val bias = rand.nextDouble() + 2
    if (pair % 2 == 0) base + bias else base - bias
  }

  def nextRandomPause(): Int = {
    val rand = new Random()
    rand.nextInt(4) + 1 // rand pause between 1 to 3 sec
  }


  private val post_header = Map(
    "Cache-Control" -> "no-cache",
    "Content-Type" -> "application/json",
    "Pragma" -> "no-cache"
  )

  val REFERENCE_POSITION: Map[String, Double] = Map("lat" -> 48.8115336, "lng" -> 2.3681119)

  val feeder: Iterator[Map[String, Double]] = Iterator.continually {
    val randomLocation = getRandomLocation(REFERENCE_POSITION)
    Map(
      "lat" -> randomLocation._1,
      "lng" -> randomLocation._2,
    )
  }

  private val scn = scenario("GetTractorsSimulation")
    .feed(feeder)
    .pause(1)
    .exec(
      http("get-tractors")
        .post("/v1/tractors?radius=340")
        .headers(post_header)
        .body(StringBody(
          """
            |{"position":{"lat":${lat},"lng":${lng}}}
            |""".stripMargin)
        )
    )

        setUp(scn.inject(constantConcurrentUsers(constantConcurrentUsersNumber.toInt) during (loadtestDuration.toInt minutes)))
          .protocols(httpProtocol)
}
