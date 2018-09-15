package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"log"
	"net/http"
	"regexp"
)

type route struct {
	pattern *regexp.Regexp
	verb    string
	handler http.Handler
}

type RegexpHandler struct {
	routes []*route
}

func (h *RegexpHandler) Handler(pattern *regexp.Regexp, verb string, handler http.Handler) {
	h.routes = append(h.routes, &route{pattern, verb, handler})
}

func (h *RegexpHandler) HandleFunc(r string, v string, handler func(http.ResponseWriter, *http.Request)) {
	re := regexp.MustCompile(r)
	h.routes = append(h.routes, &route{re, v, http.HandlerFunc(handler)})
}

func (h *RegexpHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	for _, route := range h.routes {
		if route.pattern.MatchString(r.URL.Path) && route.verb == r.Method {
			route.handler.ServeHTTP(w, r)
			return
		}
	}
	http.NotFound(w, r)
}

type lists struct {
	id           int    `json:"id,sting"`
	type_of_list string `json:"type_of_list"`
}

type list_items struct {
	list_no          int    `json:"list_no,sting"`
	item_name string `json:"item_name"`
	completed bool `json:"completed"` 
}

type Server struct {
	db *sql.DB
}

func main() {
	db, err := sql.Open("mysql", "root:@tcp(127.0.0.1:3306)/golang_todo_dev")
	if err != nil {
		log.Fatal(err)
	}
	db.SetMaxIdleConns(100)
	defer db.Close()

	server := &Server{db: db}

	reHandler := new(RegexpHandler)

	reHandler.HandleFunc("/lists/$", "GET", server.todoIndex)
	reHandler.HandleFunc("/lists/$", "POST", server.todoCreate)
	reHandler.HandleFunc("/lists/[0-9]+$", "GET", server.todoShow)
	reHandler.HandleFunc("/lists/[0-9]+$", "PUT", server.todoUpdate)
	reHandler.HandleFunc("/lists/[0-9]+$", "DELETE", server.todoDelete)

	reHandler.HandleFunc(".*.[js|css|png|eof|svg|ttf|woff]", "GET", server.assets)
	reHandler.HandleFunc("/", "GET", server.homepage)

	fmt.Println("Starting server on port 3000")
	http.ListenAndServe(":3000", reHandler)
}


func (s *Server) homepage(res http.ResponseWriter, req *http.Request) {
	http.ServeFile(res, req, "./index.html")
}

func (s *Server) assets(res http.ResponseWriter, req *http.Request) {
	http.ServeFile(res, req, req.URL.Path[1:])
}



func (s *Server) listsIndex(res http.ResponseWriter, req *http.Request) {
	var lists2 []*lists

	rows, err := s.db.Query("SELECT id, type_of_list FROM lists")
	error_check(res, err)
	for rows.Next() {
		lists := &lists{}
		rows.Scan(&lists.id, &lists.type_of_list)
		lists2 = append(lists2, lists)
	}
	rows.Close()

	jsonResponse(res, lists2)
}

func (s *Server) list_itemsIndex(res http.ResponseWriter, req *http.Request) {
	var lists3 []*list_items

	rows, err := s.db.Query("SELECT item_name FROM list_items")
	error_check(res, err)
	for rows.Next() {
		list_items := &list_items{}
		rows.Scan( &list_items.item_name)
		lists3 = append(lists3, list_items)
	}
	rows.Close()

	jsonResponse(res, lists3)
}


func (s *Server) listsCreate(res http.ResponseWriter, req *http.Request) {
	lists := &lists{}

	decoder := json.NewDecoder(req.Body)
	err := decoder.Decode(&lists)
	if err != nil {
		fmt.Println("ERROR decoding JSON - ", err)
		return
	}
	defer req.Body.Close()
	
	result, err := s.db.Exec("INSERT INTO lists(id, type_of_list) VALUES(?, ?)", lists.id, lists.type_of_list)
	if err != nil {
		fmt.Println("ERROR saving to db - ", err)
	}

	result, err := s.db.Exec("INSERT INTO list_items(list_no, item_name, completed) VALUES(?, ?, ?)", list_items.list_no, list_items.item_name, list_items.completed)
	if err != nil {
		fmt.Println("ERROR saving to db - ", err)
	}


	Id64, err := result.LastInsertId()
	id := int(id64)
	lists = &lists{id: id}

	s.db.QueryRow("SELECT completed, item_name FROM list_items WHERE Id=?", lists.id).Scan(&lists.id, &lists.type_of_list, &list_items.item_name, &list_items.completed)

	jsonResponse(res, list_items)
}

func (s *Server) todoShow(res http.ResponseWriter, req *http.Request) {
	fmt.Println("Render list_items json")
}



func (s *Server) todoDelete(res http.ResponseWriter, req *http.Request) {
	r, _ := regexp.Compile(`\d+$`)
	Id := r.FindString(req.URL.Path)
	s.db.Exec("DELETE FROM Todo WHERE Id=?", Id)
	res.WriteHeader(200)
}

func jsonResponse(res http.ResponseWriter, data interface{}) {
	res.Header().Set("Content-Type", "application/json; charset=utf-8")

	payload, err := json.Marshal(data)
	if error_check(res, err) {
		return
	}

	fmt.Fprintf(res, string(payload))
}

func error_check(res http.ResponseWriter, err error) bool {
	if err != nil {
		http.Error(res, err.Error(), http.StatusInternalServerError)
		return true
	}
	return false
}