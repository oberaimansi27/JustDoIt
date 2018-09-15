# JustDoIt
It contains the backend files of the application 'to-do-list'.
there are three files:
1. 'linking.go' which contains the linking of the do database with the postgresql database.
2. 'structure.go' which contains the structure of the go database using struct module.
3. 'todolist.go' contains the various queries for inserting and accessing the values of the database.

Schema for postgresql:
table 1: lists
id: int (primary key)
type_of_list: string 

table 2: list_items
list_no: int (foreign key to id of lists table)
item_name: string
completed: boolean
