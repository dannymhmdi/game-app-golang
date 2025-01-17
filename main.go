package main

import (
	"encoding/json"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"io"
	"log"
	"mymodule/pkg/textcolor"
	"mymodule/repository/mysql"
	"mymodule/service/registerservice"
	"net/http"
)

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/users/register", UserRegisterHandler)
	mux.HandleFunc("/users/login", UserLoginHandler)
	mux.HandleFunc("/users/profile", UserProfileHandler)
	server := http.Server{Addr: ":8080", Handler: mux}
	fmt.Println(textcolor.Green + "Server is running on port 8080" + textcolor.Reset)
	log.Fatal(server.ListenAndServe().Error())
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
			fmt.Println(textcolor.Red + fmt.Sprintf("failed to write response:%v\n", wErr) + textcolor.Reset)
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
			fmt.Println(textcolor.Red + fmt.Sprintf("failed to write response:%v\n", wErr) + textcolor.Reset)
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

func UserLoginHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		_, wErr := w.Write([]byte(`{"message":"method not allowed","status":false}`))
		if wErr != nil {
			fmt.Println(textcolor.Red + fmt.Sprintf("failed to write response:%v\n", wErr) + textcolor.Reset)

			return
		}
	}

	body, rErr := io.ReadAll(r.Body)
	if rErr != nil {
		fmt.Println("failed to read body:", rErr)

		return
	}

	bd := registerservice.LoginRequest{}

	jErr := json.Unmarshal(body, &bd)
	if jErr != nil {
		fmt.Println("failed to unmarshal body:", jErr)

		return
	}

	loginRepo := mysql.New()
	loginsvc := registerservice.New(loginRepo)

	Respone, lErr := loginsvc.Login(bd)
	if lErr != nil {
		_, wErr := w.Write([]byte(fmt.Sprintf(`{"message":%v,"status":false}`, lErr)))
		if wErr != nil {
			fmt.Println(textcolor.Red + fmt.Sprintf("failed to write response:%v\n", wErr) + textcolor.Reset)

			return
		}

		return
	}

	jsonResponse, jErr := json.Marshal(Respone)
	if jErr != nil {
		fmt.Println("failed to marshal response:", jErr)

		return
	}

	_, wErr := w.Write(jsonResponse)
	if wErr != nil {
		fmt.Println(textcolor.Red + fmt.Sprintf("failed to write response:%v\n", wErr) + textcolor.Reset)

		return
	}

}

func UserProfileHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		_, wErr := w.Write([]byte(`{"message":"method not allowed","status":false}`))
		if wErr != nil {
			fmt.Println(textcolor.Red+"failed to write response:\n", wErr)

			return
		}
	}
	body, rErr := io.ReadAll(r.Body)
	if rErr != nil {
		fmt.Println(textcolor.Red + fmt.Sprintf("failed to read body:%v\n", rErr) + textcolor.Reset)

		return
	}

	bd := registerservice.ProfileRequest{}

	jErr := json.Unmarshal(body, &bd)
	if jErr != nil {
		fmt.Println(textcolor.Red + fmt.Sprintf("failed to unmarshall json object:%v\n", rErr) + textcolor.Reset)

		return
	}
	profileRepo := mysql.New()
	profileSvc := registerservice.New(profileRepo)
	profInfo, gErr := profileSvc.GetUserProfile(bd)
	if gErr != nil {
		w.Write([]byte(gErr.Error()))

		return
	}

	response, mErr := json.Marshal(profInfo)
	if mErr != nil {
		fmt.Println(textcolor.Red + fmt.Sprintf("failed to marshall json object:%v\n", rErr) + textcolor.Reset)

		return
	}

	w.Write(response)

}

//curl -X POST -H "Content-Type: application/json" -d '{"Name":"Hosein", "PhoneNumber":"09122598501"}' http://localhost:8080/users/register
