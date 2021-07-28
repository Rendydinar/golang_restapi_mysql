package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"github.com/shopspring/decimal"
)

var db *gorm.DB
var err error

// Product is a representation of a product
type Product struct {
	ID    int             `json:"id"`
	Code  string          `json:"code"`
	Name  string          `json:"name"`
	Price decimal.Decimal `json:"price" sql:"type:decimal(16,2)"`
}

// Result is an array of product
type Result struct {
	Status  int         `json:"status"`
	Data    interface{} `json:"data"`
	Message string      `json:"message"`
}

func main() {
	// Please define your username and password for MySQL.
	db, err = gorm.Open("mysql", "root:@/golang_simple_restapi_crud?charset=utf8&parseTime=True")
	// NOTE: See we’re using = to assign the global var
	// instead of := which would assign it only in this function

	if err != nil {
		log.Println("Connection failed", err)
	} else {
		log.Println("Connection establine")
	}

	db.AutoMigrate(&Product{})

	handleRequest()
}

func handleRequest() {
	log.Println("Start the development server at http://127.0.0.1:8080")

	route := mux.NewRouter().StrictSlash(true)

	route.HandleFunc("/", homePage)
	route.HandleFunc("/api/products", createProduct).Methods("POST")
	route.HandleFunc("/api/products", getProducts).Methods("GET")
	route.HandleFunc("/api/product/{id}", getSingleProduct).Methods("GET")
	route.HandleFunc("/api/product/{id}", updateProduct).Methods("PUT")
	route.HandleFunc("/api/product/{id}", deleteProduct).Methods("DELETE")
	log.Fatal(http.ListenAndServe(":8080", route))
}

// Handle route home/index
func homePage(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Welcome")
}

// Handle create product
func createProduct(w http.ResponseWriter, r *http.Request) {
	payloads, _ := ioutil.ReadAll(r.Body)
	var product Product

	// extract data json body base by struct Product
	json.Unmarshal(payloads, &product)
	fmt.Println(product)

	// create into database
	db.Create(&product)

	res := Result{Status: 200, Data: product, Message: "Success create product"}
	result, err := json.Marshal(res)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(result)
}

// Handle get all product
func getProducts(w http.ResponseWriter, r *http.Request) {
	// declare variabel to store list product from database
	products := []Product{}

	// find product by id
	db.Find(&products)

	res := Result{Status: 200, Data: products, Message: "Success get all product"}
	results, err := json.Marshal(res)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(results)
}

// Handle get single product
func getSingleProduct(w http.ResponseWriter, r *http.Request) {
	var product Product

	// get id product by params url
	vars := mux.Vars(r)
	productID := vars["id"]

	// get product by id
	db.First(&product, productID)

	res := Result{Status: 200, Data: product, Message: "Success get product"}
	result, err := json.Marshal(res)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(result)
}

// Handle update product
func updateProduct(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	productID := vars["id"]
	payloads, _ := ioutil.ReadAll(r.Body)

	var productUpdate Product
	json.Unmarshal(payloads, &productUpdate)

	var product Product

	// update into database
	db.First(&product, productID)
	db.Model(&product).Update(productUpdate)

	res := Result{Status: 200, Data: product, Message: "Success update product"}
	result, err := json.Marshal(res)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(result)
}

// Handle delete single product
func deleteProduct(w http.ResponseWriter, r *http.Request) {
	// get id product by params url
	vars := mux.Vars(r)
	productID := vars["id"]

	// Delete product from database
	db.Delete(&Product{}, productID)

	res := Result{Status: 200, Data: "", Message: "Success delete product"}
	result, err := json.Marshal(res)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(result)
}
