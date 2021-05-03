package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"mysql/configurator"
	"net/http"
	"strconv"

	_ "github.com/go-sql-driver/mysql"
	"github.com/julienschmidt/httprouter"
)

var db *sql.DB

type Warehouse struct {
	ID           int
	Items        string
	UnitPrice    float32
	ItemCategory string
	Quantity     int
	SoldQuantity int
}

func init() {
	var err error
	db, err = sql.Open("mysql", configurator.Configurator())
	if err != nil {
		fmt.Println("Can not be connected", err)
		return
	}

	err = db.Ping()
	if err != nil {
		fmt.Println("Can not be connected in Ping", err)
		return
	}

	fmt.Println("Successfully connected to DB...")
}

func main() {
	router := httprouter.New()
	router.GET("/items", showAll)
	router.GET("/item", showOne)
	router.POST("/added", add)
	router.GET("/items/delete", delete)
	router.GET("/items/update", updateForm)
	router.POST("/items/update/process", update)

	log.Fatal(http.ListenAndServe(":8080", router))

}

func showAll(w http.ResponseWriter, req *http.Request, _ httprouter.Params) {

	rows, err := db.Query("SELECT * FROM warehouse")
	if err != nil {
		fmt.Println("Can't select the id", err)
		return
	}
	defer rows.Close()

	var its []Warehouse

	for rows.Next() {
		it := Warehouse{}
		err := rows.Scan(&it.ID, &it.Items, &it.UnitPrice, &it.ItemCategory, &it.Quantity, &it.SoldQuantity)
		if err != nil {
			fmt.Println("Can't scan the items", err)
			return
		}
		its = append(its, it)
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(its)
}

func showOne(w http.ResponseWriter, req *http.Request, _ httprouter.Params) {

	its := Warehouse{}

	var err error
	its.ID, err = strconv.Atoi(req.FormValue("id"))
	if err != nil {
		fmt.Println("Can not be converted in showOne", err)
	}
	its.Items = req.FormValue("item")
	p, _ := strconv.ParseFloat(req.FormValue("price"), 32)
	its.UnitPrice = float32(p)
	its.ItemCategory = req.FormValue("category")
	its.Quantity, _ = strconv.Atoi(req.FormValue("quantity"))
	its.SoldQuantity, _ = strconv.Atoi(req.FormValue("soldQuantity"))

	row := db.QueryRow("SELECT * FROM warehouse WHERE id=?;", its.ID)

	err = row.Scan(&its.ID, &its.Items, &its.UnitPrice, &its.ItemCategory, &its.Quantity, &its.SoldQuantity)

	if err != nil || err == sql.ErrNoRows {
		fmt.Println("Error in showOne scan", err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(its)
}

func add(w http.ResponseWriter, req *http.Request, _ httprouter.Params) {

	its := Warehouse{}

	its.Items = req.FormValue("item")
	p, _ := strconv.ParseFloat(req.FormValue("price"), 32)
	its.UnitPrice = float32(p)
	its.ItemCategory = req.FormValue("category")
	its.Quantity, _ = strconv.Atoi(req.FormValue("quantity"))
	its.SoldQuantity, _ = strconv.Atoi(req.FormValue("soldQuantity"))

	_, err := db.Exec("INSERT INTO warehouse(items,unit_price,item_category,quantity,sold_quantity) VALUES(?,?,?,?,?);", its.Items, its.UnitPrice, its.ItemCategory, its.Quantity, its.SoldQuantity)
	if err != nil {
		http.Error(w, http.StatusText(500), http.StatusInternalServerError)
		fmt.Println("Error in adding item", err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(its)
}

func delete(w http.ResponseWriter, req *http.Request, _ httprouter.Params) {

	id, err := strconv.Atoi(req.FormValue("id"))
	if err != nil {
		fmt.Println("Can not be converted into integer", err)
		return
	}

	// var row *sql.Rows
	_, err = db.Exec("DELETE FROM warehouse WHERE id=?;", id)
	if err != nil {
		http.Error(w, http.StatusText(500), http.StatusInternalServerError)
		fmt.Println("Error in delete query", err)
		return
	}
	// err = row.Scan(&id)
	// if err != nil {
	// 	http.Error(w, http.StatusText(500), http.StatusInternalServerError)
	// 	fmt.Println("Error in delete scan", err)
	// 	return
	// }
	http.Redirect(w, req, "/items", http.StatusSeeOther)
}

func updateForm(w http.ResponseWriter, req *http.Request, _ httprouter.Params) {

	id, err := strconv.Atoi(req.FormValue("id"))
	if err != nil {
		fmt.Println("Can not convert into integer in updateForm", err)
		return
	}

	row := db.QueryRow("SELECT * FROM warehouse WHERE id = ?", id)

	it := Warehouse{}
	err = row.Scan(&it.ID, &it.Items, &it.UnitPrice, &it.ItemCategory, &it.Quantity, &it.SoldQuantity)
	switch {
	case err == sql.ErrNoRows:
		http.NotFound(w, req)
		fmt.Println("error in ErrNoRows", err)
		return
	case err != nil:
		http.Error(w, http.StatusText(500), http.StatusInternalServerError)
		fmt.Println("Error in updateForm", err)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(it)
}

func update(w http.ResponseWriter, req *http.Request, _ httprouter.Params) {

	it := Warehouse{}

	it.ID, _ = strconv.Atoi(req.FormValue("id"))
	it.Items = req.FormValue("item")
	p, _ := strconv.ParseFloat(req.FormValue("price"), 32)
	it.UnitPrice = float32(p)
	it.ItemCategory = req.FormValue("category")
	it.Quantity, _ = strconv.Atoi(req.FormValue("quantity"))
	it.SoldQuantity, _ = strconv.Atoi(req.FormValue("SoldQuantity"))

	result, err := db.Exec("UPDATE warehouse SET items = ?, unit_price=?, item_category=?, quantity=?, sold_quantity=? WHERE id=?;", it.Items, it.UnitPrice, it.ItemCategory, it.Quantity, it.SoldQuantity, it.ID)
	if err != nil {
		http.Error(w, http.StatusText(500), http.StatusInternalServerError)
		fmt.Println("Error in update exec", err)
		return
	}

	fmt.Println("Result of Update is: ", result)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(it)
}
