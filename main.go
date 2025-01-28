package main

import (
	_ "github.com/go-sql-driver/mysql"
	"mymodule/config"
	"mymodule/delivery/httpserver"
	"mymodule/repository/mysql"
	"mymodule/service/authservice"
	"mymodule/service/registerservice"
	"mymodule/validator/uservalidator"
	"time"
)

const (
	signingKey             string        = "tokenpass"
	accessTokenExpireTime  time.Duration = time.Hour * 1
	refreshTokenExpireTime time.Duration = time.Hour * 24 * 7
	refreshSubject                       = "rt"
	accessSubject                        = "at"
)

func main() {
	//logFile, sErr := logger.SetUpFile("errors.log")
	//if sErr != nil {
	//	log.Fatal("failed to setup logger file")
	//}
	//defer logFile.Close()
	authSvc, userSvc := setUp()
	cfg := config.Config{
		HttpConfig: config.HttpServer{Port: "8080"},
		AuthConfig: authservice.Config{
			SigningKey:             signingKey,
			AccessTokenExpireTime:  accessTokenExpireTime,
			RefreshTokenExpireTime: refreshTokenExpireTime,
			RefreshSubject:         refreshSubject,
			AccessSubject:          accessSubject,
		},
	}
	server := httpserver.New(cfg, *authSvc, *userSvc)

	server.Serve()

	//mux := http.NewServeMux()
	//mux.HandleFunc("/users/register", UserRegisterHandler)
	//mux.HandleFunc("/users/login", UserLoginHandler)
	//mux.HandleFunc("/users/profile", UserProfileHandler)
	//
	//handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	//	w.Write([]byte("you are in route that applied middleware"))
	//})
	//
	//mux.Handle("/users", profileMiddleWare(handler))
	//server := http.Server{Addr: ":8080", Handler: mux}
	//fmt.Println(textcolor.Green + "Server is running on port 8080" + textcolor.Reset)
	//log.Fatal(server.ListenAndServe().Error())
}

//

//curl -X POST -H "Content-Type: application/json" -d '{"Name":"Hosein", "PhoneNumber":"09122598501"}' http://localhost:8080/users/register

func setUp() (*authservice.Service, *registerservice.Service) {
	cfg := authservice.Config{
		SigningKey:             signingKey,
		AccessTokenExpireTime:  accessTokenExpireTime,
		RefreshTokenExpireTime: refreshTokenExpireTime,
		RefreshSubject:         refreshSubject,
		AccessSubject:          accessSubject,
	}

	dbConfig := mysql.Config{
		Username: "gameapp",
		Password: "gameappt0lk2o20",
		Host:     "localhost",
		Port:     3308,
		DbName:   "gameapp_db",
	}
	authSvc := authservice.New(cfg)
	mysqlRepo := mysql.New(dbConfig)
	validator := uservalidator.New(mysqlRepo)
	userSvc := registerservice.New(mysqlRepo, authSvc, *validator)
	return authSvc, userSvc
}
