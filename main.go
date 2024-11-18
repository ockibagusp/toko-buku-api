package main

import (
	"net/http"
	"toko-buku-api/app"
	"toko-buku-api/controller"
	"toko-buku-api/helper"
	"toko-buku-api/middleware"
	"toko-buku-api/repository"
	"toko-buku-api/service"

	"github.com/go-playground/validator/v10"
	_ "github.com/go-sql-driver/mysql"
)

func main() {
	db := app.NewDB()
	validate := validator.New()
	categoryRepository := repository.NewCategoryRepository()
	categoryService := service.NewCategoryService(categoryRepository, db, validate)
	categoryController := controller.NewCategoryController(categoryService)
	router := app.NewRouter(categoryController)

	server := http.Server{
		Addr:    "localhost:3000",
		Handler: middleware.NewAuthMiddleware(router),
	}

	err := server.ListenAndServe()
	helper.PanicIfError(err)

	// if errors.Is(err, http.ErrServerClosed) {
	// 	fmt.Printf("server closed\n")
	// } else if err != nil {
	// 	fmt.Printf("error starting server: %s\n", err)
	// 	os.Exit(1)
	// }

	// ////
	// // start the Echo server
	// go func() {
	// 	if err := e.Start(":8000"); err != nil && err != http.ErrServerClosed {
	// 		e.Logger.Fatal("shutting down the Echo server")
	// 	}
	// }()

	// quit := make(chan os.Signal, 1)
	// signal.Notify(quit, os.Interrupt)
	// <-quit

	// // a timeout context after 10 seconds
	// ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	// defer cancel()

	// // shutdown the Echo server
	// if err := e.Shutdown(ctx); err != nil {
	// 	e.Logger.Fatal(fmt.Sprintf("failed the Echo server: %v", err))
	// } else {
	// 	e.Logger.Info("successfully the Echo server")
	// }
}
