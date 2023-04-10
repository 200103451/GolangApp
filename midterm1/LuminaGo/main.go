package main

import (
	"bytes"
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
	"html/template"
	"net/http"
)

var host = ":8080"

type Product struct {
	Id          uint16
	Brand       string
	Type        string
	Description string
	Price       float32
	Photo       string
	Rating      uint16
	RatingCount uint16
}

var products = []Product{}
var showProduct = Product{}

func getProductData(db *sql.DB) ([]Product, error) {
	rows, err := db.Query("SELECT * FROM `products`")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var item []Product
	for rows.Next() {
		var product Product
		err = rows.Scan(&product.Id, &product.Brand, &product.Type, &product.Description, &product.Price, &product.Photo, &product.Rating, &product.RatingCount)
		if err != nil {
			return nil, err
		}
		item = append(item, product)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}

	return item, nil
}

func index(w http.ResponseWriter, r *http.Request) {
	t, err := template.ParseFiles("templates/index.html", "templates/header.html", "templates/footer.html", "templates/smallcard.html")

	if err != nil {
		fmt.Fprintf(w, err.Error())
	}
	//connecting db
	db, err := sql.Open("mysql", "root:root@tcp(127.0.0.1:8889)/golang")
	if err != nil {
		panic(err)
	}
	defer db.Close()

	// Get the product data from the database
	products, err := getProductData(db)
	if err != nil {
		panic(err.Error())
	}

	t.ExecuteTemplate(w, "index", products)
}
func createUser(w http.ResponseWriter, r *http.Request) {
	t, err := template.ParseFiles("templates/registration.html", "templates/header.html", "templates/footer.html")

	if err != nil {
		fmt.Fprintf(w, err.Error())
	}

	t.ExecuteTemplate(w, "registration", nil)
}
func loginUser(w http.ResponseWriter, r *http.Request) {
	t, err := template.ParseFiles("templates/login.html", "templates/header.html", "templates/footer.html")
	if err != nil {
		fmt.Fprintf(w, err.Error())
	}

	t.ExecuteTemplate(w, "login", nil)
}
func saveUser(w http.ResponseWriter, r *http.Request) {
	nickname := r.FormValue("nickname")
	firstname := r.FormValue("firstname")
	lastname := r.FormValue("lastname")
	age := r.FormValue("age")
	password := r.FormValue("password")

	if nickname == "" || firstname == "" || lastname == "" || age == "" || password == "" {
		fmt.Fprintf(w, "You didnt write something")
	} else {
		db, err := sql.Open("mysql", "root:root@tcp(127.0.0.1:8889)/golang")
		if err != nil {
			panic(err)
		}
		defer db.Close()

		insert, err := db.Query(fmt.Sprintf("INSERT INTO `users` (`nickname`,`firstname`,`lastname`,`age`,`password`) VALUES('%s','%s','%s','%s','%s')", nickname, firstname, lastname, age, password))
		if err != nil {
			panic(err)
		}
		defer insert.Close()

		http.Redirect(w, r, "/", http.StatusSeeOther)
	}
}
func checkUser(w http.ResponseWriter, r *http.Request) {
	nickname := r.FormValue("nickname")
	password := r.FormValue("password")

	if nickname == "" || password == "" {
		fmt.Fprintf(w, "You didnt write something")
	} else {
		db, err := sql.Open("mysql", "root:root@tcp(127.0.0.1:8889)/golang")
		if err != nil {
			panic(err)
		}
		defer db.Close()

		var exists bool
		err = db.QueryRow(fmt.Sprintf("SELECT EXISTS(SELECT * from `users` WHERE `nickname`='%s' AND `password`='%s')", nickname, password)).Scan(&exists)
		if err != nil {
			panic(err)
		}
		if exists {
			http.Redirect(w, r, "/", http.StatusSeeOther)
		} else {
			fmt.Fprintf(w, "User don't exist")
		}
	}
}
func catalog(w http.ResponseWriter, r *http.Request) {
	typeFilter := r.URL.Query().Get("type")
	brandFilter := r.URL.Query().Get("brand")
	sortValue := r.URL.Query().Get("sorting")

	t, err := template.ParseFiles("templates/header.html", "templates/footer.html", "templates/smallcard.html", "templates/filters.html")
	if err != nil {
		fmt.Fprintf(w, err.Error())
	}

	db, err := sql.Open("mysql", "root:root@tcp(127.0.0.1:8889)/golang")
	if err != nil {
		panic(err)
	}
	defer db.Close()

	query := "SELECT * FROM `products`"
	if typeFilter != "" {
		query += " WHERE `type`='" + typeFilter + "'"
	}
	if brandFilter != "" {
		if typeFilter == "" {
			query += " WHERE "
		} else {
			query += " AND "
		}
		query += "`brand`='" + brandFilter + "'"
	}
	if sortValue != "" {
		query += " ORDER BY `" + sortValue + "` DESC"
	}

	rows, err := db.Query(query)
	if err != nil {
		panic(err)
	}
	defer rows.Close()

	//iterate over results and write HTML for each row
	var count = 0
	buf := new(bytes.Buffer)
	err = t.ExecuteTemplate(buf, "header", nil)
	fmt.Fprintf(w, "%s", buf.String()) //header
	if typeFilter != "" {
		fmt.Fprintf(w, "<div class = \"badge\">\n    <a>Catalog</a>\n <a>Searched type of product: %s </a>\n </div>", typeFilter)
	} else if brandFilter != "" {
		fmt.Fprintf(w, "<div class = \"badge\">\n    <a>Catalog</a>\n <a>Searched brand of product: %s </a>\n </div>", brandFilter)
	} else {
		fmt.Fprintf(w, "<div class = \"badge\">\n    <a>Catalog</a>\n </div>")
	}

	fmt.Fprintf(w, "<div class =\"catalog-container\">")
	buf = new(bytes.Buffer)
	err = t.ExecuteTemplate(buf, "filters", nil)
	fmt.Fprintf(w, "%s", buf.String()) //filters
	fmt.Fprintf(w, "<div class=\"catalog-products\">")
	fmt.Fprintf(w, "<div class =\"row\">")
	for rows.Next() {
		var id uint16
		var brand string
		var productType string
		var description string
		var price float32
		var photoAddress string
		var rating uint16
		var ratingCount uint16

		err = rows.Scan(&id, &brand, &productType, &description, &price, &photoAddress, &rating, &ratingCount)
		if err != nil {
			panic(err.Error())
		}

		if count == 3 {
			fmt.Fprintf(w, "</div>")
			fmt.Fprintf(w, "<div class=\"row\">")
			count = 0
		}
		// write HTML for product row
		fmt.Fprintf(w, "<div class=\"card\">\n  <div class=\"card-header\">\n    <h2>%s</h2>\n  </div>\n  <div class=\"card-image\">\n    <img src=\"/static/images/products/%s\" alt=\"some product\">\n  </div>\n  <p>Price: $%0.2f</p>\n  <div class=\"card-description\">\n    <p>%s</p>\n  </div>\n  <form>\n    <button formaction=\"/product/%d\">Product Overview</button>\n <button>Add to cart</button>\n  </form>\n</div>", brand, photoAddress, price, description, id)
		count += 1
	}
	fmt.Fprintf(w, "</div>")
	fmt.Fprintf(w, "</div>")
	buf = new(bytes.Buffer)
	err = t.ExecuteTemplate(buf, "footer", nil)
	fmt.Fprintf(w, "%s", buf.String()) //footer
}

func productFullInfo(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	t, err := template.ParseFiles("templates/product.html", "templates/header.html", "templates/footer.html")
	if err != nil {
		fmt.Fprintf(w, err.Error())
	}

	db, err := sql.Open("mysql", "root:root@tcp(127.0.0.1:8889)/golang")
	if err != nil {
		panic(err)
	}
	defer db.Close()

	//vyborka dannyh
	rows, err := db.Query(fmt.Sprintf("SELECT * FROM `products` WHERE `id` = '%s'", vars["id"]))
	if err != nil {
		panic(err)
	}
	defer rows.Close()

	var item Product
	for rows.Next() {
		var product Product
		err = rows.Scan(&product.Id, &product.Brand, &product.Type, &product.Description, &product.Price, &product.Photo, &product.Rating, &product.RatingCount)
		if err != nil {
			panic(err)
		}
		item = product
	}
	if err = rows.Err(); err != nil {
		panic(err)
	}

	t.ExecuteTemplate(w, "product", item)
}

func handleFunc() {
	r := mux.NewRouter()
	r.HandleFunc("/", index).Methods("GET")
	r.HandleFunc("/registration", createUser).Methods("GET")
	r.HandleFunc("/save_user", saveUser).Methods("POST")
	r.HandleFunc("/login", loginUser).Methods("GET")
	r.HandleFunc("/check_user", checkUser).Methods("POST")
	r.HandleFunc("/catalog", catalog).Methods("GET")
	r.HandleFunc("/product/{id:[0-9]+}", productFullInfo).Methods("GET")

	http.Handle("/", r)
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("./static/"))))
	fmt.Printf("server is listening on host %s \n", host)
	http.ListenAndServe(host, nil)
}

func main() {
	handleFunc()
}
