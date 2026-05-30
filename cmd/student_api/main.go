package main

import (
	"context"
	"fmt"
	"log"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/siddharth7actowiz/student_api/internal/config"
	"github.com/siddharth7actowiz/student_api/internal/http/handlers/student"
	"github.com/siddharth7actowiz/student_api/internal/storage/mysql"
)

func main() {
	//load config
	cfg := config.MustLoad("config/local.yaml")

	//db setup
	//it is plugin you can add storage interface as you want add postgres,sql ,sqlite as you want
	//dependency injection
	storage, err := mysql.New(cfg)
	if err != nil {
		log.Fatal("Error creating mysql storage:", err)
	}
	slog.Info("MySQL Storage Created Successfully", slog.String("env:", cfg.Env))

	//setup router
	router := http.NewServeMux()
	
	
	// Route for create student
	router.HandleFunc("POST /api/students", student.New(storage))

	//Get Student by Id
	//endpint/{}-->dynamic query
	router.HandleFunc("GET /api/student/{id}",student.GetByID(storage))	

	//get all students
	router.HandleFunc("GET /api/students",student.GetAll(storage))	


	//update 
	router.HandleFunc(
	"PUT /api/student/{id}",
	student.Update(storage),
)	

	//delete
	router.HandleFunc("DELETE /api/student/{id}",student.Delete(storage))	

	//Channel for Synchronization
	done := make(chan os.Signal, 1)
	//signal
	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM) //for handling any incoming interrupt signals then notify
	//setup http server
	server := http.Server{
		Addr:    cfg.HTTPServer.Addr,
		Handler: router,
	}
	fmt.Println(cfg.HTTPServer.Addr)
	//Using lgs like production Quality is good practice because it provides structured logging and better performance than the standard log package. It also allows you to easily add context to your logs, which can be very helpful for debugging and monitoring your application.
	slog.Info("Server Started on: ", slog.String("Addr", cfg.HTTPServer.Addr))

	go func() {
		err := server.ListenAndServe()
		if err != nil {
			log.Fatal("Error starting server:", err)
		}
	}()

	<-done //wait for the signal to be received
	slog.Info("Shutting Down the Server..")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err = server.Shutdown(ctx) //close the server gracefully
	if err != nil {
		slog.Error("Error closing server:", slog.String("error", err.Error()))
	}

	slog.Info("Server Stopped Succesfully")
}
