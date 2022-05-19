package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"os"

	pb "grpccrud/book_server/genproto/book_service"

	"github.com/jackc/pgx/v4"
	"google.golang.org/grpc"
)

const (
	port = ":40041"
)

func NewBookManagementServer() *BookManagementServer {
	return &BookManagementServer{}
}

type BookManagementServer struct {
	conn *pgx.Conn
	pb.UnimplementedBookManagementServer
}

func (server *BookManagementServer) Run() error {
	lis, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("Failed to listen. %v", err)
	}

	s := grpc.NewServer()
	pb.RegisterBookManagementServer(s, server)
	log.Printf("Server listening at %v", lis.Addr())
	return s.Serve(lis)
}

func (server *BookManagementServer) Create(ctx context.Context, in *pb.NewBook) (*pb.Book, error) {

	log.Printf("Received: %v", in.GetBookName())
	createSql := `
	create table if not exists books(
		book_id SERIAL PRIMARY KEY, 
		book_name text,
		author text,
		user_id int,
		CONSTRAINT fk_user FOREIGN KEY(user_id) REFERENCES users(id)
	);
	`
	_, err := server.conn.Exec(context.Background(), createSql)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Table creation failed: %v\n", err)
		os.Exit(1)
	}

	created_book := &pb.Book{
		BookName: in.GetBookName(),
		Author:   in.GetAuthor(),
		UserId:   in.GetUserId(),
	}
	tx, err := server.conn.Begin(context.Background())
	if err != nil {
		log.Fatalf("conn.Begin failed. %v", err)
	}
	_, err = tx.Exec(context.Background(), "insert into books(book_name, author, user_id) values ($1, $2, $3)", in.BookName, in.Author, in.UserId)
	if err != nil {
		log.Fatalf("tx.Exec failed: %v", err)
	}
	tx.Commit(context.Background())

	return created_book, nil
}

func (server *BookManagementServer) GetAll(ctx context.Context, in *pb.BookParams) (*pb.BookList, error) {

	var books_list *pb.BookList = &pb.BookList{}
	rows, err := server.conn.Query(context.Background(), "select * from books")
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		book := pb.Book{}
		err = rows.Scan(&book.BookId, &book.BookName, &book.Author, &book.UserId)
		if err != nil {
			return nil, err
		}
		books_list.Books = append(books_list.Books, &book)
	}

	return books_list, nil
}

func (server *BookManagementServer) GetUserBooks(ctx context.Context, in *pb.UserBookRequest) (*pb.UserBook, error) {

	var books_array *pb.UserBook = &pb.UserBook{}
	row, err := server.conn.Query(context.Background(), "select book_id, book_name, author from books where user_id=$1", in.UserId)
	if err != nil {
		return nil, err
	}
	defer row.Close()

	for row.Next() {
		book := pb.Book{}
		err = row.Scan(&book.BookId, &book.BookName, &book.Author)
		if err != nil {
			return nil, err
		}
		books_array.Book = append(books_array.Book, &book)
	}

	return books_array, nil
}

func (server *BookManagementServer) Get(ctx context.Context, in *pb.BookRequest) (*pb.BookResponse, error) {

	row, err := server.conn.Query(context.Background(), "select * from books where book_id=$1", in.BookId)
	if err != nil {
		return nil, err
	}
	defer row.Close()

	if !row.Next() {
		if err := row.Err(); err != nil {
			log.Fatalf("failed to retrieve data from ToDo-> %v" + err.Error())
		}
		log.Fatalf("ToDo with ID='%d' is not found", in.BookId)
	}

	book := pb.Book{}
	if err := row.Scan(&book.BookId, &book.BookName, &book.Author, &book.UserId); err != nil {
		log.Fatalf("failed to retrieve data from ToDo-> %v" + err.Error())
	}

	if err != nil {
		log.Fatalf("failed to retrieve data from ToDo-> %v" + err.Error())
	}

	if row.Next() {
		log.Fatalf("ToDo with ID='%d' is not found", in.BookId)
	}

	return &pb.BookResponse{Book: &book}, nil
}

func (server *BookManagementServer) Update(ctx context.Context, in *pb.UpdateBookRequest) (*pb.UpdateBookResponse, error) {

	res, err := server.conn.Exec(context.Background(), "update books set book_name=$2, author=$3, user_id=$4 where book_id=$1", in.BookId, in.Book.BookName, in.Book.Author, in.Book.UserId)
	if err != nil {
		log.Fatalf("Unable to execute the query. %v", err)
	}

	row := res.RowsAffected()

	fmt.Printf("Total rows/record affected %v", row)

	return &pb.UpdateBookResponse{Success: true}, nil
}

func (server *BookManagementServer) Delete(ctx context.Context, in *pb.BookRequest) (*pb.UpdateBookResponse, error) {

	res, err := server.conn.Exec(context.Background(), "delete from books where book_id=$1", in.BookId)
	if err != nil {
		log.Fatalf("Unable to execute the query. %v", err)
	}
	rows := res.RowsAffected()

	fmt.Printf("Total rows/record affected %v", rows)

	return &pb.UpdateBookResponse{Success: true}, nil
}

func main() {

	database_url := "postgres://postgres:RIVOJmz777@localhost:5432/postgres"
	conn, err := pgx.Connect(context.Background(), database_url)
	if err != nil {
		log.Fatalf("Unable to establish connection: %v", err)
	}
	defer conn.Close(context.Background())
	var book_server *BookManagementServer = NewBookManagementServer()

	book_server.conn = conn
	if err := book_server.Run(); err != nil {
		log.Fatalf("failed to server: %v", err)
	}
}
