package main
import(
	"fmt"
	"net/http"
	"encoding/json"
	"strings"
	"math/rand"
	"time"
)

var db = make(map[string]string)

type Short1 struct { // long url
	Long string `json:"long_url"`
}

type Short2 struct { // short url
	Short string `json:"short_url"`
}

func sub(){
	// handle the route
	http.HandleFunc("/html", HandleFront)
	http.HandleFunc("/short", handleShort)
	http.HandleFunc("/", handleRedirect)

	fmt.Println("http://localhost:8080/html")
	err := http.ListenAndServe(":8080", nil)

	if err != nil {
		fmt.Println(err)
	}
}

// generate the code
func gen() string {
	const set = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

	random := rand.New(rand.NewSource(time.Now().UnixNano())) // for random set generation

	var sb strings.Builder
	for i := 0 ; i < 6; i++ {
		sb.WriteByte(set[random.Intn(len(set))])
	}

	return sb.String()
}

// to handle post for shorting url
func handleShort(w http.ResponseWriter , r *http.Request){
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid", http.StatusBadRequest)
		return
	}
	// decoding that json
	var req Short1
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil || req.Long == "" {
		http.Error(w, "Error having url", http.StatusBadRequest)
		return
	}
	//1. generate the code
	code := gen()

	// 2. save it on above map
	db[code] = req.Long

	// 3. response maker
	res := Short2 {
		Short : fmt.Sprintf("http://localhost:8080/%s", code),
	}

	// 4. finally coinnect with http
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(res)
	
}

func HandleFront(w http.ResponseWriter , r *http.Request){
	if r.Method != http.MethodGet {
		http.Error(w, "Invalid", http.StatusBadRequest)
		return
	}

		http.ServeFile(w , r , "./index.html")
}

// to get "get" response from it 
func handleRedirect(w http.ResponseWriter , r *http.Request){
	if r.Method != http.MethodGet {
		http.Error(w, "Invalid", http.StatusBadRequest)
		return
	}

	// extracting the code " /this-code "
	code := strings.TrimPrefix(r.URL.Path , "/")

	// if typed : localhost:8080/

	if code == "" {
		fmt.Fprint(w, "Welcome to the api ... !!")
		return
	}

	// look up in db
	Long , exist := db[code]

	if !exist {
		http.Error(w , "No url exists" , http.StatusNotFound)
		return
	}


	// redirect to the browser
	http.Redirect(w , r , Long , http.StatusFound)
}