package main

import (
	"encoding/json"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"io"
	"log"
	"mymodule/repository/mysql"
	"mymodule/service/registerservice"
	"net/http"
)

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/users/register", UserRegisterHandler)
	server := http.Server{Addr: ":8080", Handler: mux}
	//http.HandleFunc("/users/register", UserRegisterHandler)
	fmt.Println("Server is running on port 8080")
	log.Fatal(server.ListenAndServe())
}

func UserRegisterHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		w.Write([]byte(`{"message":"method not allowed","status":false}`))

		return
	}
	body, err := io.ReadAll(r.Body)
	if err != nil {
		_, wErr := w.Write([]byte(`{"message":"failed to read body","status":false}`))
		if wErr != nil {
			fmt.Println("failed to write response:", wErr)
			http.Error(w, wErr.Error(), http.StatusInternalServerError)
			return
		}

		return
	}

	defer r.Body.Close()

	bd := registerservice.RegisterRequest{}

	uErr := json.Unmarshal(body, &bd)
	if uErr != nil {
		_, wErr := w.Write([]byte(`{"message":"failed to unmarshal body","status":false}`))
		if wErr != nil {
			fmt.Println("failed to write response:", wErr)
			http.Error(w, wErr.Error(), http.StatusInternalServerError)
			return
		}

		return
	}

	registerRepo := mysql.New()
	newRegisterSvc := registerservice.New(registerRepo)
	createdUser, rErr := newRegisterSvc.RegisterUser(bd)
	if rErr != nil {
		_, wErr := w.Write([]byte(fmt.Sprintf(`{"message":%s,"status":false}\n`, rErr.Error())))
		if wErr != nil {
			fmt.Println("failed to write client response: ", wErr)

			return
		}
		return
	}

	jsonData, mErr := json.Marshal(createdUser)
	if mErr != nil {
		fmt.Println("failed to marshal response: ", mErr)

		return
	}
	w.Write(jsonData)
}

//curl -X POST -H "Content-Type: application/json" -d '{"Name":"Hosein", "PhoneNumber":"09122598501"}' http://localhost:8080/users/register
