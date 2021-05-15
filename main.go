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
	ID           int     `json:"id"`
	Items        string  `json:"items"`
	UnitPrice    float32 `json:"unit_price"`
	ItemCategory string  `json:"item_category"`
	Quantity     int     `json:"quantity"`
	SoldQuantity int     `json:"sold_quantity"`
}

func init() {
	var err error
	db, err = sql.Open("mysql", configurator.Configurator())
	if err != nil {
		fmt.Println("Can not be connected to DB", err)
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
	router.GET("/items", allItems)
	router.GET("/item", oneItem)
	router.POST("/items/insert", insertItem)
	router.GET("/items/delete", deleteItem)
	router.PATCH("/items/update/:id", updateItem)

	log.Fatal(http.ListenAndServe(":8080", router))

}

func updateItem(w http.ResponseWriter, req *http.Request, ps httprouter.Params) {

	var it Warehouse
	var err error

	idQuery := ps.ByName("id")
	it.ID, err = strconv.Atoi(idQuery)
	if err != nil {
		fmt.Println("Error in converting int", err)
		return
	}

	selectRow := db.QueryRow("SELECT * FROM warehouse WHERE id = ?", it.ID)

	err = selectRow.Scan(&it.ID, &it.Items, &it.UnitPrice, &it.ItemCategory, &it.Quantity, &it.SoldQuantity)
	switch {
	case err == sql.ErrNoRows:
		http.Error(w, http.StatusText(400), http.StatusBadRequest)
		fmt.Println("Row doesn't exist in Database", err)
		return
	case err != nil:
		http.Error(w, http.StatusText(500), http.StatusInternalServerError)
		fmt.Println("Error in update select", err)
		return
	}

	json.NewDecoder(req.Body).Decode(&it)

	result, dbErr := db.Prepare("UPDATE warehouse SET items = ?, unit_price=?, item_category=?, quantity=?, sold_quantity=? WHERE id=?;")
	if dbErr != nil {
		http.Error(w, http.StatusText(500), http.StatusInternalServerError)
		fmt.Println("Error in update prepare", err)
		return
	}
	defer result.Close()

	res, err := result.Exec(it.Items, it.UnitPrice, it.ItemCategory, it.Quantity, it.SoldQuantity, it.ID)
	if err != nil {
		fmt.Println("Problem in update exec")
		return
	}
	_, err = res.RowsAffected()
	if err != nil {
		fmt.Println("Problem in rowsAffected")
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(it)
}

func allItems(w http.ResponseWriter, req *http.Request, _ httprouter.Params) {

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

func oneItem(w http.ResponseWriter, req *http.Request, _ httprouter.Params) {

	its := Warehouse{}

	idQuery := req.URL.Query().Get("id") // map[string][]string
	id, err := strconv.Atoi(idQuery)
	if err != nil {
		fmt.Println("Error in converting int", err)
		return
	}

	row := db.QueryRow("SELECT * FROM warehouse WHERE id=?;", id)

	err = row.Scan(&its.ID, &its.Items, &its.UnitPrice, &its.ItemCategory, &its.Quantity, &its.SoldQuantity)

	if err != nil || err == sql.ErrNoRows {
		fmt.Println("Doesnt exist in DB", err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(its)
}

func insertItem(w http.ResponseWriter, req *http.Request, _ httprouter.Params) {

	it := Warehouse{}

	json.NewDecoder(req.Body).Decode(&it)

	_, err := db.Exec("INSERT INTO warehouse(items,unit_price,item_category,quantity,sold_quantity) VALUES(?,?,?,?,?);", it.Items, it.UnitPrice, it.ItemCategory, it.Quantity, it.SoldQuantity)
	if err != nil {
		http.Error(w, http.StatusText(500), http.StatusInternalServerError)
		fmt.Println("Error in adding item", err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(it)

}

func deleteItem(w http.ResponseWriter, req *http.Request, _ httprouter.Params) {

	idQuery := req.URL.Query().Get("id") // map[string][]string
	id, err := strconv.Atoi(idQuery)
	if err != nil {
		fmt.Println("Error in converting int", err)
		return
	}

	// var row *sql.Rows
	_, err = db.Exec("DELETE FROM warehouse WHERE id=?;", id)
	if err != nil {
		http.Error(w, http.StatusText(500), http.StatusInternalServerError)
		fmt.Println("Error in delete query", err)
		return
	}

	http.Redirect(w, req, "/items", http.StatusSeeOther)
}
