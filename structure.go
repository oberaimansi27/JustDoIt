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