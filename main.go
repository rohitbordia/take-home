package main

import (
	"context"
	"database/sql"
	b64 "encoding/base64"
	"fmt"
	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
	"google.golang.org/grpc"
	"io"
	"log"
	"net"
	"os"
	"os/exec"
	db2 "take-home/db"
	pb "take-home/take-home/grpc"
	"time"
)

type server struct {
	pb.UnimplementedScriptServiceServer
	db sqlx.DB
}

type Script struct {
	Id            int64          `db:"id"`
	Name          string         `db:"name"`
	Status        string         `db:"status"`
	Content       string         `db:"content"`
	LastRunStatus sql.NullString `db:"last_run_status"`
	CreatedAt     time.Time      `db:"created_at"`
	UpdatedAt     time.Time      `db:"updated_at"`
}

func (s *server) CreateTask(taskServer pb.ScriptService_CreateTaskServer) error {
	log.Print("Creating task ")
	for {
		req, err := taskServer.Recv()
		if err == io.EOF {
			fmt.Println("server recv eof")
			return nil
		}
		if err != nil {
			fmt.Printf("failed to recv: %v\n", err)
			return err
		}

		if req.ScriptName == "" {
			return taskServer.SendAndClose(&pb.ScriptResponse{
				ScriptName:   "",
				ScriptStatus: "",
				Error: &pb.Error{
					Code: "400",
					Type: "UNKOWN",
					Desc: "script name cannot be empty",
				},
			})
		}
		if req.Content == "" {
			return taskServer.SendAndClose(&pb.ScriptResponse{
				ScriptName:   "",
				ScriptStatus: "",
				Error: &pb.Error{
					Code: "400",
					Type: "UNKOWN",
					Desc: "script content cannot be empty",
				},
			})
		}
		sDec, _ := b64.StdEncoding.DecodeString(req.Content)
		fmt.Println(string(sDec))
		fmt.Println()
		tx := s.db.MustBegin()
		tx.MustExec("INSERT INTO script ( name, status, content) VALUES ($1, $2, $3)", req.ScriptName, "CREATED", sDec)
		tx.Commit()

		return taskServer.SendAndClose(&pb.ScriptResponse{
			ScriptName:   req.ScriptName,
			ScriptStatus: "Created",
		})
	}

	return nil
}

func (s *server) ExecuteTask(ctx context.Context, in *pb.ScriptRequest) (*pb.ScriptResponse, error) {
	log.Printf("Received: %v", in.ScriptName)
	script := Script{}
	err := s.db.Get(&script, "select * from script where name=$1", in.ScriptName)
	if err != nil {
		return nil, errors.Wrapf(err, "Unable to find script status for:%s", in.ScriptName)
	}
	f, err := os.Create("/tmp/dat2")
	_, err = f.WriteString(script.Content)
	tx := s.db.MustBegin()
	if err != nil {
		s.db.Exec("UPDATE script set last_run_status=$1, updated_at=CURRENT_TIMESTAMP where id=$2", "Failed",script.Id)
		return nil, errors.WithMessage(err, "Failed to create temp file")
	}
	out, err := exec.Command("/bin/sh", "/tmp/dat2").Output()
	if err != nil {

		s.db.Exec("UPDATE script set last_run_status=$1, updated_at=CURRENT_TIMESTAMP where id=$2", "Failed", script.Id)
		return nil, errors.WithMessage(err, "Failed to execute temp file")
	}
	fmt.Printf("output is %s\n", out)
	s.db.Exec("UPDATE script set last_run_status=$1, updated_at=CURRENT_TIMESTAMP where id=$2", "Executed", script.Id)
	tx.Commit()
	return &pb.ScriptResponse{
		ScriptName:   in.ScriptName,
		ScriptStatus: script.Status,
		LastRunStatus: "Executed",
		Output: string(out),
	}, nil
}
func (s *server) GetTaskStatus(ctx context.Context, in *pb.ScriptRequest) (*pb.ScriptResponse, error) {
	log.Printf("Received: %v", in.ScriptName)
	if in.ScriptName == "" {
		return &pb.ScriptResponse{
			ScriptName:   "",
			ScriptStatus: "",
			Error: &pb.Error{
				Code: "400",
				Type: "UNKOWN",
				Desc: "script name cannot be empty",
			},
		}, nil
	}
	script := Script{}
	err := s.db.Get(&script, "select * from script where name=$1", in.ScriptName)
	if err != nil {
		return nil, errors.Wrapf(err, "Unable to find script status for:%s", in.ScriptName)
	}

	return &pb.ScriptResponse{
		ScriptName:    in.ScriptName,
		ScriptStatus:  script.Status,
		LastRunStatus: script.LastRunStatus.String,
	}, nil
}
func (s *server) GetTaskSource(ctx context.Context, in *pb.ScriptRequest) (*pb.ScriptResponse, error) {
	log.Printf("Received: %v", in.ScriptName)
	script := Script{}
	err := s.db.Get(&script, "select * from script where name=$1", in.ScriptName)
	if err != nil {
		return nil, errors.Wrapf(err, "Unable to find script status for:%s", in.ScriptName)
	}
	return &pb.ScriptResponse{
		ScriptName:   script.Name,
		ScriptStatus: script.Status,
		Content:      script.Content,
	}, nil
}

func main() {

	fmt.Println("Go gRPC take-home project!")

	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", 9000))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	grpcServer := grpc.NewServer()

	db, dbClose := db2.Connect()
	defer dbClose()
	pb.RegisterScriptServiceServer(grpcServer, &server{db: *db})

	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %s", err)
	}
}
