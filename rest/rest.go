package rest

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"pkg/blockchain"
	"pkg/p2p"
	"pkg/utils"
	"pkg/wallet"

	"github.com/gorilla/mux"
)

var port string

type url string

type urlDescription struct {
	URL         url    `json:"url"`
	Method      string `json:"method"`
	Description string `json:"description"`
	Payload     string `json:"payload,omitempty"`
}

type errorResponse struct {
	ErrorMessage string `json:"errorNotFound"`
}

type balanceResponse struct {
	Address string `json:"address"`
	Balance int    `json:"balance"`
}

type myWalletResponse struct {
	Address string `json:"address"`
}

type addTxPayload struct {
	To     string
	Amount int
}

type addPeerPayload struct {
	address, port string
}

func (u url) MarshalText() ([]byte, error) {
	url := fmt.Sprintf("http://localhost:%s%s", "4000", u)
	return []byte(url), nil
}

func documentation(w http.ResponseWriter, r *http.Request) {
	data := []urlDescription{
		{
			URL:         "/",
			Method:      "GET",
			Description: "SEE DOC",
		},
		{
			URL:         "/status",
			Method:      "GET",
			Description: "SEE The Status of Blockchain",
		},
		{
			URL:         "/blocks",
			Method:      "POST",
			Description: "SEE DOC",
			Payload:     "data:string",
		},
		{
			URL:         "/blocks/{hash}",
			Method:      "GET",
			Description: "SEE DOC",
		},
		{
			URL:         "/balance/{address}",
			Method:      "GET",
			Description: "Get Tx for an Address",
		},
		{
			URL:         "/ws",
			Method:      "GET",
			Description: "Upgrade to WebSocket",
		},
	}
	json.NewEncoder(w).Encode(data)
}

func blocks(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		utils.HandleErr(json.NewEncoder(w).Encode(blockchain.Blocks(blockchain.Blockchain())))

	case "POST":
		blockchain.Blockchain().AddBlock()
		w.WriteHeader(http.StatusCreated)
	}
}

// Middleware
func jsonContentTypeMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Content-Type", "application/json")
		next.ServeHTTP(w, r)
	})
}

func block(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r) // ID 가가져져옴옴
	hash := vars["hash"]
	encoder := json.NewEncoder(w)
	block, err := blockchain.FindBlock(hash)

	if err == blockchain.ErrNotFound {
		json.NewEncoder(w).Encode(errorResponse{fmt.Sprint(err)})
	} else {
		encoder.Encode(block)
	}
}

func status(w http.ResponseWriter, r *http.Request) {
	json.NewEncoder(w).Encode(blockchain.Blockchain())
}

func balance(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	address := vars["address"]
	total := r.URL.Query().Get("total")
	switch total {
	case "true":
		amount := blockchain.BalanceByAddress(address, blockchain.Blockchain())
		json.NewEncoder(w).Encode(balanceResponse{address, amount})
	default:
		utils.HandleErr(json.NewEncoder(w).Encode(blockchain.UTxOutsByAddress(address, blockchain.Blockchain())))
	}
}

func mempool(w http.ResponseWriter, r *http.Request) {
	utils.HandleErr(json.NewEncoder(w).Encode(blockchain.Mempool.Txs))
}

func transactions(w http.ResponseWriter, r *http.Request) {
	var payload addTxPayload
	utils.HandleErr(json.NewDecoder(r.Body).Decode(&payload))
	err := blockchain.Mempool.AddTx(payload.To, payload.Amount)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(errorResponse{err.Error()})
		return
	}
	w.WriteHeader(http.StatusCreated)
}

func myWallet(w http.ResponseWriter, r *http.Request) {
	address := wallet.Wallet().Address
	json.NewEncoder(w).Encode(myWalletResponse{Address: address})
	// json.NewEncoder(w).Encode(struct{ Address string `json:"address"`}{Address : address})
}

func loggerMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Println(r.URL)
		next.ServeHTTP(w, r)
	})
}

func peers(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "POST":
		var payload addPeerPayload
		json.NewDecoder(r.Body).Decode(&payload)
		p2p.AddPeer(payload.address, payload.port)
		w.WriteHeader(http.StatusOK)
	}
}

func Start(aPort int) {
	router := mux.NewRouter()
	port = fmt.Sprintf(":%d", aPort)
	router.Use(jsonContentTypeMiddleware, loggerMiddleware) // Accept Middleware
	router.HandleFunc("/", documentation).Methods("GET")    // ONLY HANDLE GET
	router.HandleFunc("/status", status)
	router.HandleFunc("/blocks", blocks).Methods("GET", "POST")
	router.HandleFunc("/blocks/{hash:[a-f0-9]+}", block).Methods("GET")
	router.HandleFunc("/balance/{address}", balance)
	router.HandleFunc("/mempool", mempool).Methods("GET")
	router.HandleFunc("/wallet", myWallet).Methods("GET")
	router.HandleFunc("/transactions", transactions).Methods("POST")
	router.HandleFunc("/ws", p2p.Upgrade).Methods("GET")
	router.HandleFunc("/peers", peers).Methods("POST")

	fmt.Printf("Listening on http://localhost:%s", port)
	log.Fatal(http.ListenAndServe(port, router))
}
