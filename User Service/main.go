package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"os"

	pb "grpccrud/user_server/genproto/user_service"

	"github.com/jackc/pgx/v4"
	"google.golang.org/grpc"
)

const (
	port = ":50051"
)

func NewUserManagementServer() *UserManagementServer {
	return &UserManagementServer{}
}

type UserManagementServer struct {
	conn *pgx.Conn
	pb.UnimplementedUserManagementServer
}

func (server *UserManagementServer) Run() error {
	lis, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	s := grpc.NewServer()
	pb.RegisterUserManagementServer(s, server)
	log.Printf("server listening at %v", lis.Addr())
	return s.Serve(lis)
}

func (server *UserManagementServer) CreateNewUser(ctx context.Context, in *pb.NewUser) (*pb.User, error) {

	log.Printf("Received: %v", in.GetName())
	createSql := `
	create table if not exists users(
		id SERIAL PRIMARY KEY, 
		name text, 
		age int
	);
	`
	_, err := server.conn.Exec(context.Background(), createSql)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Table creation failed: %v\n", err)
		os.Exit(1)
	}

	created_user := &pb.User{
		Name: in.GetName(),
		Age:  in.GetAge(),
	}
	tx, err := server.conn.Begin(context.Background())
	if err != nil {
		log.Fatalf("conn.Begin failed. %v", err)
	}
	_, err = tx.Exec(context.Background(), "insert into users(name, age) values ($1, $2)", in.Name, in.Age)
	if err != nil {
		log.Fatalf("tx.Exec failed: %v", err)
	}
	tx.Commit(context.Background())

	return created_user, nil
}

func (server *UserManagementServer) GetUsers(ctx context.Context, in *pb.GetUsersParams) (*pb.UserList, error) {

	var users_list *pb.UserList = &pb.UserList{}
	rows, err := server.conn.Query(context.Background(), "select * from users")
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		user := pb.User{}
		err = rows.Scan(&user.Id, &user.Name, &user.Age)
		if err != nil {
			return nil, err
		}
		users_list.Users = append(users_list.Users, &user)
	}

	return users_list, nil
}

func (server *UserManagementServer) GetById(ctx context.Context, in *pb.UserRequest) (*pb.UserRequestResponse, error) {

	row, err := server.conn.Query(context.Background(), "select * from users where id=$1", in.Id)
	if err != nil {
		return nil, err
	}
	defer row.Close()

	if !row.Next() {
		if err := row.Err(); err != nil {
			log.Fatalf("failed to retrieve data from ToDo-> %v" + err.Error())
		}
		log.Fatalf("ToDo with ID='%d' is not found", in.Id)
	}

	user := pb.User{}
	if err := row.Scan(&user.Id, &user.Name, &user.Age); err != nil {
		log.Fatalf("failed to retrieve data from ToDo-> %v" + err.Error())
	}

	if err != nil {
		log.Fatalf("failed to retrieve data from ToDo-> %v" + err.Error())
	}

	if row.Next() {
		log.Fatalf("ToDo with ID='%d' is not found", in.Id)
	}

	return &pb.UserRequestResponse{User: &user}, nil
}

func (server *UserManagementServer) UpdateUser(ctx context.Context, in *pb.UpdateUserRequest) (*pb.UpdateUserResponse, error) {

	res, err := server.conn.Exec(context.Background(), "update users set name=$2, age=$3 where id=$1", in.Id, in.User.Name, in.User.Age)
	if err != nil {
		log.Fatalf("Unable to execute the query. %v", err)
	}

	row := res.RowsAffected()

	fmt.Printf("Total rows/record affected %v", row)

	return &pb.UpdateUserResponse{Success: true}, nil
}

func (server *UserManagementServer) DeleteUser(ctx context.Context, in *pb.UserRequest) (*pb.UserResponse, error) {

	res, err := server.conn.Exec(context.Background(), "delete from users where id=$1", in.Id)
	if err != nil {
		log.Fatalf("Unable to execute the query. %v", err)
	}
	rows := res.RowsAffected()

	fmt.Printf("Total rows/record affected %v", rows)

	return &pb.UserResponse{Success: true}, nil
}

func main() {

	database_url := "postgres://postgres:RIVOJmz777@localhost:5432/postgres"
	conn, err := pgx.Connect(context.Background(), database_url)
	if err != nil {
		log.Fatalf("Unable to establish connection: %v", err)
	}
	defer conn.Close(context.Background())
	var user_mgmt_server *UserManagementServer = NewUserManagementServer()

	user_mgmt_server.conn = conn
	if err := user_mgmt_server.Run(); err != nil {
		log.Fatalf("failed to server: %v", err)
	}
}
