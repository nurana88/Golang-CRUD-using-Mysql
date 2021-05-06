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
	ID           int     `json:"id,omitempty"             sql:"id"`
	Items        string  `json:"items,omitempty"          sql:"items"`
	UnitPrice    float32 `json:"unit_price,omitempty"     sql:"unit_price"`
	ItemCategory string  `json:"item_category,omitempty"  sql:"item_category"`
	Quantity     int     `json:"quantity,omitempty"       sql:"quantity"`
	SoldQuantity int     `json:"sold_quantity,omitempty"  sql:"sold_quantity"`
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
	router.GET("/items", showAll)
	router.GET("/item", showOne)
	router.POST("/insert", insert)
	// router.GET("/items/delete", delete)
	router.PATCH("/items/update", update)

	log.Fatal(http.ListenAndServe(":8080", router))

}

func update(w http.ResponseWriter, req *http.Request, _ httprouter.Params) {
	idQuery := req.URL.Query()["id"] // map[string][]string
	id, err := strconv.Atoi(idQuery[0])
	if err != nil {
		fmt.Println("Error in converting int", err)
		return
	}

	selectRow := db.QueryRow("SELECT * FROM warehouse WHERE id = ?", id)

	it := Warehouse{}
	err = selectRow.Scan(&it.ID, &it.Items, &it.UnitPrice, &it.ItemCategory, &it.Quantity, &it.SoldQuantity)
	switch {
	case err == sql.ErrNoRows:
		http.NotFound(w, req)
		fmt.Println("error in ErrNoRows in select", err)
		return
	case err != nil:
		http.Error(w, http.StatusText(500), http.StatusInternalServerError)
		fmt.Println("Error in update select", err)
		return
	}

	result, err := db.Exec("UPDATE warehouse SET items = ?, unit_price=?, item_category=?, quantity=?, sold_quantity=? WHERE id=?;", it.Items, it.UnitPrice, it.ItemCategory, it.Quantity, it.SoldQuantity, it.ID)
	if err != nil {
		http.Error(w, http.StatusText(500), http.StatusInternalServerError)
		fmt.Println("Error in update exec", err)
		return
	}

	fmt.Println(result)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(it)
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

	idQuery := req.URL.Query()["id"] // map[string][]string
	id, err := strconv.Atoi(idQuery[0])
	if err != nil {
		fmt.Println("Error in converting int", err)
		return
	}

	row := db.QueryRow("SELECT * FROM warehouse WHERE id=?;", id)

	err = row.Scan(&its.ID, &its.Items, &its.UnitPrice, &its.ItemCategory, &its.Quantity, &its.SoldQuantity)

	if err != nil || err == sql.ErrNoRows {
		fmt.Println("Error in showOne scan", err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(its)
}

func insert(w http.ResponseWriter, req *http.Request, _ httprouter.Params) {

	var err error
	var price int
	it := Warehouse{}

	idQuery := req.URL.Query()["id"] // map[string][]string
	it.ID, err = strconv.Atoi(idQuery[0])
	if err != nil {
		fmt.Println("Error in converting int", err)
		return
	}
	priceQuery := req.URL.Query()["unit_price"] // map[string][]string
	price, err = strconv.Atoi(priceQuery[2])
	if err != nil {
		fmt.Println("Error in converting int", err)
		return
	}

	it.UnitPrice = float32(price)

	qQuery := req.URL.Query()["quantity"] // map[string][]string
	it.Quantity, err = strconv.Atoi(qQuery[4])
	if err != nil {
		fmt.Println("Error in converting int", err)
		return
	}
	sqQuery := req.URL.Query()["sold_quantity"] // map[string][]string
	it.SoldQuantity, err = strconv.Atoi(sqQuery[5])
	if err != nil {
		fmt.Println("Error in converting int", err)
		return
	}

	_, err = db.Exec("INSERT INTO warehouse(items,unit_price,item_category,quantity,sold_quantity) VALUES(?,?,?,?,?) WHERE id=?;", it.Items, it.UnitPrice, it.ItemCategory, it.Quantity, it.SoldQuantity, it.ID)
	if err != nil {
		http.Error(w, http.StatusText(500), http.StatusInternalServerError)
		fmt.Println("Error in adding item", err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(it)
}

// func delete(w http.ResponseWriter, req *http.Request, _ httprouter.Params) {

// 	id, err := strconv.Atoi(req.FormValue("id"))
// 	if err != nil {
// 		fmt.Println("Can not be converted into integer", err)
// 		return
// 	}

// 	// var row *sql.Rows
// 	_, err = db.Exec("DELETE FROM warehouse WHERE id=?;", id)
// 	if err != nil {
// 		http.Error(w, http.StatusText(500), http.StatusInternalServerError)
// 		fmt.Println("Error in delete query", err)
// 		return
// 	}
// 	// err = row.Scan(&id)
// 	// if err != nil {
// 	// 	http.Error(w, http.StatusText(500), http.StatusInternalServerError)
// 	// 	fmt.Println("Error in delete scan", err)
// 	// 	return
// 	// }
// 	http.Redirect(w, req, "/items", http.StatusSeeOther)
// }
