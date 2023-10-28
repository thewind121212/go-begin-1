package restAPI

import (
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"io"
	"log"
	"net/http"
	"os"
	"path"
	"playground/utils/database"
)

type ImageGalary struct {
	ID           string
	ThumbNail    string
	Gallery      any
	ProductColor []ProductColor
}

type ProductColor struct {
	IdColor    any
	ColorImage any
}

func InitializeGinRestAPI(port string) {
	fmt.Println("Initializing Gin Rest API")
	r := gin.Default()
	mongoClient := database.MongoClient

	db := mongoClient.Database("totodayshopClone")
	collection := db.Collection("jacketmixes")

	cursor, err := collection.Find(context.Background(), bson.M{})

	if err != nil {
		fmt.Println("Error to find data", err)
	}

	//rootFolder := "storage/images"

	if err != nil {
		fmt.Println("Error to create folder", err)
	}

	var results []bson.M
	var imageGalary []ImageGalary

	for cursor.Next(context.Background()) {
		var result bson.M
		err := cursor.Decode(&result)
		if err != nil {
			log.Fatal(err)
		}
		results = append(results, result)
	}

	for _, document := range results {

		var imageGalaryTemp ImageGalary
		var productColorTempArray []ProductColor
		imageGalaryTemp.ID = document["id"].(string)
		imageGalaryTemp.ThumbNail = document["thumbnail"].(string)
		imageGalaryTemp.Gallery = document["gallery"]

		for key, value := range document["productInfo"].(bson.M) {
			var productColorTemp ProductColor
			productColorTemp.IdColor = key
			productColorTemp.ColorImage = value

			productColorTempArray = append(productColorTempArray, productColorTemp)

		}

		imageGalaryTemp.ProductColor = productColorTempArray

		imageGalary = append(imageGalary, imageGalaryTemp)

	}

	for _, value := range imageGalary {
		color := value.ProductColor
		id := value.ID

		for _, value := range color {
			colorId := value.IdColor.(string)
			err = os.MkdirAll("storage/images/"+id+"/color/"+colorId, 0755)
			for key, value := range value.ColorImage.(bson.M) {
				if key == "colorImage" {
					fileName := path.Base(value.(string))
					res, err := http.Get(value.(string))
					if err != nil {
						fmt.Println("Error to download image", err)
					}
					preCreateOutPutFile, err := os.Create("storage/images/" + id + "/color/" + colorId + "/" + fileName)
					if err != nil {
						fmt.Println("Error to create file", err)
					}
					_, err = io.Copy(preCreateOutPutFile, res.Body)
				}
			}
		}

		//id := value.ID
		//urlThumbNail := value.ThumbNail
		//gallery := value.Gallery
		//fileThumbNailName := path.Base(urlThumbNail)
		//
		//res, err := http.Get(urlThumbNail)
		//if err != nil {
		//	fmt.Println("Error to download image", urlThumbNail)
		//	panic(err)
		//}
		//if res.StatusCode != 200 {
		//	fmt.Println("Error to download image", urlThumbNail)
		//}
		//
		//err = os.MkdirAll("storage/images/"+id+"/thumbnail", 0755)
		//if err != nil {
		//	fmt.Println("Error to create folder", err)
		//}
		//outputFile, err := os.Create("storage/images/" + id + "/thumbnail/" + fileThumbNailName)
		//
		//if err != nil {
		//	fmt.Println("Error to create file", err)
		//}
		//
		//_, err = io.Copy(outputFile, res.Body)
		//
		//if err != nil {
		//	panic(err)
		//}
		//
		//err = os.Mkdir("storage/images/"+id+"/gallery", 0755)
		//for _, value := range gallery.(bson.A) {
		//	res, err := http.Get(value.(string))
		//	fileName := path.Base(value.(string))
		//	if err != nil {
		//		fmt.Println("Error to download image", value.(string))
		//		panic(err)
		//	}
		//
		//	outputFile, err := os.Create("storage/images/" + id + "/gallery/" + fileName)
		//
		//	_, err = io.Copy(outputFile, res.Body)
		//	if err != nil {
		//		panic(err)
		//	}
		//
		//	if err != nil {
		//		fmt.Println("Error to create folder", err)
		//	}

		//}

		fmt.Println("--------------------------")

	}

	//fmt.Println(imageGalary)

	r.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "pong",
			"status":  "success",
		})
	})

	if err := r.Run(); err != nil {
		panic(err)
	}

}
