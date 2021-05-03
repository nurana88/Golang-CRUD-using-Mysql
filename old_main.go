package main

import (
	"database/sql"
	"fmt"
	"net/http"
	"strconv"
	"text/template"

	_ "github.com/go-sql-driver/mysql"
)

var db *sql.DB
var tpl *template.Template

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
	db, err = sql.Open("mysql", "root:mypassword@tcp(127.0.0.1:3306)/testdb")
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
	tpl = template.Must(template.ParseGlob("templates/*.html"))
}

func main() {
	http.HandleFunc("/items", showAll)
	http.HandleFunc("/item", showOne)
	http.HandleFunc("/added", add)
	http.HandleFunc("/add", addForm)
	http.HandleFunc("/items/delete", delete)
	http.HandleFunc("/items/update", updateForm)
	http.HandleFunc("/items/update/process", update)

	http.ListenAndServe(":8080", nil)

}

func showAll(w http.ResponseWriter, req *http.Request) {

	if req.Method != "GET" {
		http.Error(w, http.StatusText(405), 405)
		fmt.Println("Error in showing", req.Method)
		return
	}

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
	tpl.ExecuteTemplate(w, "items.html", its)
}

func showOne(w http.ResponseWriter, req *http.Request) {
	if req.Method != "GET" {
		http.Error(w, http.StatusText(405), http.StatusMethodNotAllowed)
		fmt.Println("Error in showOne", req.Method)
		return
	}

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

	tpl.ExecuteTemplate(w, "show.html", its)

}

func addForm(w http.ResponseWriter, req *http.Request) {
	tpl.ExecuteTemplate(w, "addForm.html", nil)
}

func add(w http.ResponseWriter, req *http.Request) {
	if req.Method != "POST" {
		http.Error(w, http.StatusText(405), http.StatusMethodNotAllowed)
		fmt.Println("Error in posting", req.Method)
		return
	}

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

	tpl.ExecuteTemplate(w, "added.html", its)

}

func delete(w http.ResponseWriter, req *http.Request) {
	if req.Method != "GET" {
		http.Error(w, http.StatusText(405), http.StatusMethodNotAllowed)
		fmt.Println("Error in deleting", req.Method)
		return
	}

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

func updateForm(w http.ResponseWriter, req *http.Request) {
	if req.Method != "GET" {
		http.Error(w, http.StatusText(405), http.StatusMethodNotAllowed)
		fmt.Println("Error in update form", req.Method)
		return
	}
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
	tpl.ExecuteTemplate(w, "updateForm.html", it)

}

func update(w http.ResponseWriter, req *http.Request) {
	if req.Method != "POST" {
		http.Error(w, http.StatusText(405), http.StatusMethodNotAllowed)
		fmt.Println("Error in update", req.Method)
		return
	}

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

	tpl.ExecuteTemplate(w, "updated.html", it)

}
