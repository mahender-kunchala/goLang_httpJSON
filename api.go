package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"strconv"
	"sync"
	"time"
)

type UsersDetails struct {
	Id            string   `json:"Id"`
	SecretCode    string   `json:"SecretCode"`
	Name          string   `json:"Name"`
	Address       string   `json:"Address"`
	PhoneNumber   string   `json:"PhoneNumber"`
	UserType      string   `json:"UserType"`
	Requested     []string `json:"Requested"`
	PendingReq    []string `json:"PendingReq"`
	ConnectedUser []string `json:"ConnectedUser"`
}

var count = 2
var usersmap map[string]UsersDetails

func welcomePage(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Welcome to the HomePage!")
	fmt.Println("welcome page reached")
}

func LoginUser(w http.ResponseWriter, r *http.Request) {
	type UsersCredentials struct {
		SecretCode string `json:"SecretCode"`
	}

	var login UsersCredentials
	reqBody, _ := ioutil.ReadAll(r.Body)
	json.Unmarshal(reqBody, &login)
	var user UsersDetails
	user = usersmap[login.SecretCode]
	json.NewEncoder(w).Encode(user)
	fmt.Println("login page reached")
}
func increment(wg *sync.WaitGroup, m *sync.Mutex) {

	m.Lock()
	count = count + 1
	m.Unlock()
	wg.Done()
}

func CreateUser(w http.ResponseWriter, r *http.Request) {
	var user UsersDetails
	b := 999999
	a := 100000
	reqBody, _ := ioutil.ReadAll(r.Body)
	json.Unmarshal(reqBody, &user)
	var p sync.WaitGroup
	var m sync.Mutex
	increment(&p, &m)

	user.Id = strconv.Itoa(count)
	rand.Seed(time.Now().UnixNano())
	user.SecretCode = strconv.Itoa(a + rand.Intn(b-a+1))
	usersmap[user.SecretCode] = user
	json.NewEncoder(w).Encode(user)
	fmt.Println("create user page reached")
}

func UpdateUser(w http.ResponseWriter, r *http.Request) {
	var users UsersDetails
	reqBody, _ := ioutil.ReadAll(r.Body)
	json.Unmarshal(reqBody, &users)
	userCode := users.SecretCode
	for m, n := range usersmap {
		if n.SecretCode == userCode {
			if users.Name != " " {
				n.Name = users.Name
			}
			if users.Address != " " {
				n.Address = users.Address
			}
			if users.PhoneNumber != " " {
				n.PhoneNumber = users.PhoneNumber
			}
			if users.UserType != " " {
				n.UserType = users.UserType
			}
			usersmap[m] = n
		}
	}
	fmt.Println("update user  page reached")
}

func GetUser(w http.ResponseWriter, r *http.Request) {
	type usersID struct {
		Id string `json:"Id"`
	}
	var userid usersID
	reqBody, _ := ioutil.ReadAll(r.Body)
	json.Unmarshal(reqBody, &userid)
	var Scode string

	for k := range usersmap {
		if usersmap[k].Id == userid.Id {
			Scode = k
		}
	}
	var user UsersDetails
	user = usersmap[Scode]
	json.NewEncoder(w).Encode(user)
	fmt.Println("get user page reached")
}

func GetAllusers(w http.ResponseWriter, r *http.Request) {
	for k := range usersmap {
		json.NewEncoder(w).Encode(usersmap[k])
	}
	fmt.Println("get all users page reached")
}

func GetAllDonors(w http.ResponseWriter, r *http.Request) {
	type UsersCredentials struct {
		Id         string `json:"Id"`
		SecretCode string `json:"SecretCode"`
	}

	var patientcredentials UsersCredentials
	reqBody, _ := ioutil.ReadAll(r.Body)
	json.Unmarshal(reqBody, &patientcredentials)
	Patientid := patientcredentials.Id
	PatientUser := usersmap[patientcredentials.SecretCode]

	if PatientUser.UserType == "Patient" {
		for k := range usersmap {
			if usersmap[k].UserType == "Donor" {
				var DonorUser UsersDetails
				found := 0
				DonorUser.Id = usersmap[k].Id
				DonorUser.Name = usersmap[k].Name
				for _, queryDonor := range usersmap[k].ConnectedUser {
					if queryDonor == Patientid {
						found = 1
					}
				}
				if found == 1 {
					DonorUser.Address = usersmap[k].Address
					DonorUser.PhoneNumber = usersmap[k].PhoneNumber
					DonorUser.UserType = usersmap[k].UserType
				}
				if found == 0 {
					DonorUser.Address = " "
					DonorUser.PhoneNumber = " "
					DonorUser.UserType = " "
				}
				DonorUser.Requested = usersmap[k].Requested
				DonorUser.PendingReq = usersmap[k].PendingReq
				DonorUser.ConnectedUser = usersmap[k].ConnectedUser
				json.NewEncoder(w).Encode(DonorUser)

			}

		}

	}
	fmt.Println("get all donors page reached")

}
func GetAllPatients(w http.ResponseWriter, r *http.Request) {
	type DonorCredentials struct {
		DonorId         string `json:"Id"`
		DonorSecretCode string `json:"SecretCode"`
	}
	var donorcredentials DonorCredentials
	reqBody, _ := ioutil.ReadAll(r.Body)
	json.Unmarshal(reqBody, &donorcredentials)
	DonorId := donorcredentials.DonorId
	DonorUser := usersmap[donorcredentials.DonorSecretCode]
	if DonorUser.UserType == "Donor" {
		for k := range usersmap {
			if usersmap[k].UserType == "Patient" {
				var PatientUser UsersDetails
				found := 0
				PatientUser.Id = usersmap[k].Id
				PatientUser.Name = usersmap[k].Name
				for _, queryPatient := range usersmap[k].ConnectedUser {
					if queryPatient == DonorId {
						found = 1
					}
				}
				if found == 1 {
					PatientUser.Address = usersmap[k].Address
					PatientUser.PhoneNumber = usersmap[k].PhoneNumber
					PatientUser.UserType = usersmap[k].UserType
				}
				if found == 0 {
					PatientUser.Address = " "
					PatientUser.PhoneNumber = " "
					PatientUser.UserType = " "
				}
				PatientUser.Requested = usersmap[k].Requested
				PatientUser.PendingReq = usersmap[k].PendingReq
				PatientUser.ConnectedUser = usersmap[k].ConnectedUser
				json.NewEncoder(w).Encode(PatientUser)
			}
		}
	}
	fmt.Println("get all patients page reached")
}

func SendRequest(w http.ResponseWriter, r *http.Request) {
	type requestDetails struct {
		DonorId           string `json:"Id"`
		PatientSecretCode string `json:"SecretCode"`
	}
	var request requestDetails
	reqBody, _ := ioutil.ReadAll(r.Body)

	json.Unmarshal(reqBody, &request)
	Patientcode := request.PatientSecretCode
	DonorId := request.DonorId

	PatientId := usersmap[Patientcode].Id
	PatientUser := usersmap[Patientcode]

	PatientUser.PendingReq = append(PatientUser.PendingReq, DonorId)
	usersmap[Patientcode] = PatientUser
	json.NewEncoder(w).Encode(PatientUser)

	var DonorUser UsersDetails

	for k, v := range usersmap {
		if usersmap[k].Id == DonorId {
			DonorUser = v
		}
	}
	DonorUser.Requested = append(DonorUser.Requested, PatientId)
	usersmap[DonorUser.SecretCode] = DonorUser
	json.NewEncoder(w).Encode(DonorUser)
	fmt.Println("send request page reached")

}
func AcceptRequest(w http.ResponseWriter, r *http.Request) {
	type acceptDetails struct {
		PatientId       string `json:"Id"`
		DonarSecretCode string `json:"SecretCode"`
	}
	var accept acceptDetails
	reqBody, _ := ioutil.ReadAll(r.Body)
	json.Unmarshal(reqBody, &accept)
	PatientId := accept.PatientId
	var PatientUser UsersDetails
	Donorcode := accept.DonarSecretCode
	DonorId := usersmap[Donorcode].Id
	DonorUser := usersmap[Donorcode]
	for j := range DonorUser.Requested {
		if DonorUser.Requested[j] == PatientId {
			DonorUser.Requested = append(DonorUser.Requested[:j], DonorUser.Requested[j+1:]...)
			DonorUser.ConnectedUser = append(DonorUser.ConnectedUser, PatientId)
		}
	}
	usersmap[Donorcode] = DonorUser
	json.NewEncoder(w).Encode(DonorUser)

	for k, v := range usersmap {
		if usersmap[k].Id == PatientId {
			PatientUser = v
		}
	}
	for i := range PatientUser.PendingReq {
		if PatientUser.PendingReq[i] == DonorId {
			PatientUser.PendingReq = append(PatientUser.PendingReq[:i], PatientUser.PendingReq[i+1:]...)
			PatientUser.ConnectedUser = append(PatientUser.ConnectedUser, DonorId)
		}
	}
	usersmap[PatientUser.SecretCode] = PatientUser
	json.NewEncoder(w).Encode(PatientUser)
	fmt.Println("accept request page reached")
}

func CancelConnection(w http.ResponseWriter, r *http.Request) {
	type cancelConnDetails struct {
		DonorId           string `json:"Id"`
		PatientSecretCode string `json:"SecretCode"`
	}
	var ConCancel cancelConnDetails
	reqBody, _ := ioutil.ReadAll(r.Body)
	json.Unmarshal(reqBody, &ConCancel)
	PatientSd := ConCancel.PatientSecretCode
	DonorId := ConCancel.DonorId
	PatientId := usersmap[PatientSd].Id
	PatientUser := usersmap[PatientSd]
	var DonorUser UsersDetails
	for i := range PatientUser.ConnectedUser {
		if PatientUser.ConnectedUser[i] == DonorId {
			PatientUser.ConnectedUser = append(PatientUser.ConnectedUser[:i], PatientUser.ConnectedUser[i+1:]...)
		}
	}
	usersmap[PatientSd] = PatientUser
	json.NewEncoder(w).Encode(PatientUser)
	for k, v := range usersmap {
		if usersmap[k].Id == DonorId {
			DonorUser = v
		}
	}
	for j := range DonorUser.ConnectedUser {
		if DonorUser.ConnectedUser[j] == PatientId {

			DonorUser.ConnectedUser = append(DonorUser.ConnectedUser[:j], DonorUser.ConnectedUser[j+1:]...)
		}
	}

	usersmap[DonorUser.SecretCode] = DonorUser
	json.NewEncoder(w).Encode(DonorUser)
	fmt.Println("cancel connection page reached")
}

func CancelRequest(w http.ResponseWriter, r *http.Request) {
	type cancelReqDetails struct {
		DonorId           string `json:"Id"`
		PatientSecretCode string `json:"SecretCode"`
	}
	var reqCancel cancelReqDetails
	reqBody, _ := ioutil.ReadAll(r.Body)
	json.Unmarshal(reqBody, &reqCancel)
	PatientSd := reqCancel.PatientSecretCode
	DonorId := reqCancel.DonorId
	PatientId := usersmap[PatientSd].Id
	PatientUser := usersmap[PatientSd]
	var DonorUser UsersDetails
	for i := range PatientUser.PendingReq {
		if PatientUser.PendingReq[i] == DonorId {
			PatientUser.PendingReq = append(PatientUser.PendingReq[:i], PatientUser.PendingReq[i+1:]...)
		}
	}
	usersmap[PatientSd] = PatientUser
	json.NewEncoder(w).Encode(PatientUser)
	for k, v := range usersmap {
		if usersmap[k].Id == DonorId {
			DonorUser = v
		}
	}
	for j := range DonorUser.Requested {
		if DonorUser.Requested[j] == PatientId {
			DonorUser.Requested = append(DonorUser.Requested[:j], DonorUser.Requested[j+1:]...)
		}
	}
	usersmap[DonorUser.SecretCode] = DonorUser
	json.NewEncoder(w).Encode(DonorUser)
	fmt.Println("cancel request  page reached")

}

func DeleteUser(w http.ResponseWriter, r *http.Request) {
	type UsersCredentials struct {
		Id         string `json:"Id"`
		SecretCode string `json:"SecretCode"`
	}
	var login UsersCredentials
	reqBody, _ := ioutil.ReadAll(r.Body)
	json.Unmarshal(reqBody, &login)
	var user UsersDetails
	user = usersmap[login.SecretCode]
	delete(usersmap, login.SecretCode)
	json.NewEncoder(w).Encode(user)
	fmt.Println("delete user page reached")
}
func main() {
	usersmap = make(map[string]UsersDetails)
	usersmap["123456"] = UsersDetails{Id: "1", SecretCode: "123456", Name: "maahi", Address: "hyderabad", PhoneNumber: "8282771313", UserType: "Donor", Requested: []string{""}, PendingReq: []string{""}, ConnectedUser: []string{"2"}}
	usersmap["234456"] = UsersDetails{Id: "2", SecretCode: "234456", Name: "hari", Address: "chennai", PhoneNumber: "7293823323", UserType: "Patient", Requested: []string{""}, PendingReq: []string{""}, ConnectedUser: []string{"1"}}
	http.HandleFunc("/", welcomePage)
	http.HandleFunc("/login", LoginUser)
	http.HandleFunc("/createUser", CreateUser)
	http.HandleFunc("/updateUser", UpdateUser)
	http.HandleFunc("/getAllDonors", GetAllDonors)
	http.HandleFunc("/getAllPatients", GetAllPatients)
	http.HandleFunc("/getAllusers", GetAllusers)
	http.HandleFunc("/getUser", GetUser)
	http.HandleFunc("/sendRequest", SendRequest)
	http.HandleFunc("/acceptRequest", AcceptRequest)
	http.HandleFunc("/cancelConnection", CancelConnection)
	http.HandleFunc("/cancelRequest", CancelRequest)
	http.HandleFunc("/deleteUser", DeleteUser)
	log.Fatal(http.ListenAndServe(":5000", nil))
}
