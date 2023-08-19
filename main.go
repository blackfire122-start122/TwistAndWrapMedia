package main

import (
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"io"
	"net/http"
	"os"
	"time"
)

type User struct {
	gorm.Model
	Id    uint64 `gorm:"primaryKey"`
	Image string
}

type Product struct {
	gorm.Model
	Id    uint64 `gorm:"primaryKey"`
	Image string
}

var DB *gorm.DB

func main() {
	db, err := gorm.Open(postgres.New(postgres.Config{DSN: "host=localhost user=u0_a146 dbname=postgres port=5432 sslmode=disable TimeZone=Asia/Shanghai"}), &gorm.Config{})

	if err != nil {
		panic("failed to connect database " + err.Error())
	}

	DB = db

	http.HandleFunc("/saveImageUser", saveImage)
	http.HandleFunc("/saveImageProduct", saveImage)
	if err := http.ListenAndServe(":8000", nil); err != nil {
		panic("failed run serve " + err.Error())
	}
}

func saveImage(w http.ResponseWriter, r *http.Request) {
	queryValues := r.URL.Query()
	id := queryValues.Get("id")
	file, handler, err := r.FormFile("file")

	if err != nil {
		http.Error(w, "Error reading file", http.StatusBadRequest)
		return
	}
	defer file.Close()

	var path string
	var filename = generateUniqueFileName(handler.Filename)

	if r.URL.Path == "/saveImageUser" {
		path = "media/userImages/"
		var user User
		if err := DB.First(&user, id).Error; err != nil {
			http.Error(w, "Error find user", http.StatusBadRequest)
			return
		}

		user.Image = filename
		if err := DB.Save(&user).Error; err != nil {
			http.Error(w, "Error save user", http.StatusBadRequest)
			return
		}
	} else if r.URL.Path == "/saveImageProduct" {
		path = "media/productImages/"
		var product Product
		if err := DB.First(&product, id).Error; err != nil {
			http.Error(w, "Error find product", http.StatusBadRequest)
			return
		}

		product.Image = filename
		if err := DB.Save(&product).Error; err != nil {
			http.Error(w, "Error save product", http.StatusBadRequest)
			return
		}
	}

	f, err := os.OpenFile(path+filename, os.O_WRONLY|os.O_CREATE, 0666)
	if err != nil {
		http.Error(w, "Error saving file", http.StatusInternalServerError)
		return
	}

	defer f.Close()
	io.Copy(f, file)

	w.WriteHeader(http.StatusOK)
}

func generateUniqueFileName(originalName string) string {
	currentTime := time.Now().Format("20060102150405")
	return currentTime + "_" + originalName
}
