package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"
)

func getGlovoFoodItems(url string) error {

}

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal(err)
	}
	glovoUrl := os.Getenv("glovo_url")
	err = getGlovoFoodItems(glovoUrl)
	fmt.Println("Running")
}
