package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strconv"

	"github.com/gin-gonic/gin"
)

var tractors []tractor

type latlong struct {
	Lat float64 `json:"x"`
	Lng float64 `json:"y"`
}

type position struct {
	Lat float64 `json:"lat"`
	Lng float64 `json:"lng"`
}

type userInput struct {
	Position position `json:"position"`
}

/**
Tractor JSON structure looks like:
{
    "raison_social": "EARL FERME DU MONT",
    "nom_commune": "Saint-Denis-lès-Bourg",
    "lat_long": {
      "x": 46.2111998453,
      "y": 5.18462806771
    },
    "image_url": "https://touslestracteurs.com/images/iseki/TG233.jpg",
    "modele": "TG233",
    "categorie": "Tracteur utilitaire compact",
    "annee_fabrication": [
      1994,
      1995,
      1996,
      1997,
      1998,
      1999,
      2000,
      2001,
      2002,
      2003,
      2004
    ]
  }
*/

type tractor struct {
	Title            string  `json:"raison_social"`
	Modele           string  `json:"modele"`
	Position         latlong `json:"lat_long"`
	NomCommune       string  `json:"nom_commune"`
	Categorie        string  `json:"categorie"`
	ImageUrl         string  `json:"image_url"`
	AnneeFabrication []int16 `json:"annee_fabrication"`
}

type tractorWithDistance struct {
	Title            string  `json:"raison_social"`
	Modele           string  `json:"modele"`
	Position         latlong `json:"lat_long"`
	NomCommune       string  `json:"nom_commune"`
	Categorie        string  `json:"categorie"`
	ImageUrl         string  `json:"image_url"`
	AnneeFabrication []int16 `json:"annee_fabrication"`
	Distance         int     `json:"distance"`
}

func loadDataset() []tractor {
	tractorJsonFilePath := os.Getenv("DATA_JSON_FILE_PATH")
	content, err := os.ReadFile(tractorJsonFilePath)
	if err != nil {
		fmt.Println("Error during Unmarshal(): ", err)
	}

	json.Unmarshal(content, &tractors)

	return tractors
}

func ping(ctx *gin.Context) {
	ctx.String(http.StatusOK, "{\"Hello\": \"Sotracteur\"}")
}

func noRoute(ctx *gin.Context) {
	ctx.Status(http.StatusNotImplemented)
}

/**
 * Find all products that are 10km for the position of the user
 */
func allTractorsAround(ctx *gin.Context) {
	userInput := userInput{}
	radiusParam, exist := ctx.GetQuery("radius")
	if !exist {
		fmt.Println("params missing")
		ctx.Status(http.StatusBadRequest)
		return
	}

	radius, err := strconv.Atoi(radiusParam)
	if err != nil {
		fmt.Println("convertion params failed")
		ctx.Status(http.StatusBadRequest)
		return
	}

	err = ctx.ShouldBindJSON(&userInput)
	if err != nil {
		fmt.Println("failed bind json:", err)
		ctx.Status(http.StatusBadRequest)
		return
	}

	var tractorsMatched []tractorWithDistance

	for _, tractor := range tractors {
		matchedTractors := tractor.checkDistance(userInput.Position, radius)
		if matchedTractors != nil {
			tractorsMatched = append(tractorsMatched, *matchedTractors)
		}
	}
	fmt.Println("GO: tractors returned for userPosition: ", userInput.Position)
	ctx.JSON(http.StatusOK, tractorsMatched)
}

func main() {

	// load json
	fmt.Println("Loading json")
	loadDataset()
	fmt.Println("finish Loading json, key: ", len(tractors))

	r := gin.New()
	r.NoRoute(noRoute)
	r.GET("/ping", ping)
	v1 := r.Group("/v1")
	{
		v1.POST("/tractors", allTractorsAround)
	}

	// listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
	r.Run()
}
