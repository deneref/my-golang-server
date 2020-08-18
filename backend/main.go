package main

import (
	"database/sql"
	"encoding/json"
	"encoding/xml"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"

	"my-golang-server/backend/models"
	"my-golang-server/backend/orders"

	"github.com/gorilla/mux"
)

const (
	MY_URL = "YOUR_API_URL"
	/*
		url returns json like
		"{"rv":202892,"content":[
			{"goods": [{"gid": 139297, "price": 999, "status": "t", "chrt_id": 20866595}, {"gid": 139298, "price": 999, "status": "t", "chrt_id": 20866596}], "status": "u", "order_id": 361833725, "store_id": 10629, "date_created": "2019-09-06T13:48:33.34579+03:00"},
		{"goods": [{"gid": 139308, "price": 799, "status": "t", "chrt_id": 20912708}], "status": "u", "order_id": 361834045, "store_id": 10927, "date_created": "2019-09-06T14:20:13.07318+03:00"}]"
	*/
	ON_PAGE = 30
)

type application struct {
	servicePort      int
	ordersRepository orders.Orders
	s                *mux.Router
}

var conf models.Config

func init() {
	models.LoadConfig(&conf)
}

func main() {
	app := NewApplication(conf)
	app.initServer()

	go func() {
		app.UpdateOrdersList()
	}()

	log.Fatal(http.ListenAndServe(fmt.Sprintf(`:%d`, app.servicePort), app.s))
}

func NewApplication(conf models.Config) *application {
	ord := orders.NewOrdersRepository(new(sql.DB), conf)

	return &application{
		servicePort:      8147,
		ordersRepository: ord,
	}
}

func (app *application) initServer() {
	app.s = mux.NewRouter().StrictSlash(true)

	app.s.HandleFunc("/allOrders", app.AllOrdersHandler).Methods(http.MethodGet)
	app.s.HandleFunc("/allOrdersXML", app.AllOrdersXMLHandler)
	app.s.HandleFunc("/getOrder", app.GivenOrderHandler)
	app.s.HandleFunc("/getPage", app.GivenPageHandler)
	app.s.HandleFunc("/health", statusHandler).Methods(http.MethodGet)
}

// statusHandler return 200 for zabbix
func statusHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
}

func (app *application) AllOrdersHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	data, err := app.ordersRepository.List()
	if err != nil {
		log.Print(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	json, err := json.Marshal(data)
	if err != nil {
		log.Print(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write(json)
}

func (app *application) GivenOrderHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	r.ParseForm()
	id, err := strconv.Atoi(r.Form["id"][0])
	if err != nil {
		log.Print(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	data, err := app.ordersRepository.GetByOrderID(id)
	if len(data.Content) == 0 {
		w.WriteHeader(http.StatusNoContent)
		fmt.Fprintf(w, "Order with id of %d was not found", id)
		return
	}
	json, err := json.Marshal(data)
	if err != nil {
		log.Print(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write(json)
}

func (app *application) GivenPageHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	r.ParseForm()
	page, err := strconv.Atoi(r.Form["page"][0])
	if err != nil {
		log.Print(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	data, err := app.ordersRepository.GetByPageNum(page, ON_PAGE)
	if len(data.Content) == 0 {
		w.WriteHeader(http.StatusNoContent)
		fmt.Fprintf(w, "No content on this page №%d", page)
		return
	}
	json, err := json.Marshal(data)
	if err != nil {
		log.Print(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write(json)

}

func (app *application) AllOrdersXMLHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/xml")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	data, err := app.ordersRepository.List()
	if err != nil {
		log.Print(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	xml, err := xml.Marshal(data)
	if err != nil {
		log.Print(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write(xml)
}

func (app *application) UpdateOrdersList() {
	data, err := orders.MakeRequest(MY_URL)
	var currRv int
	if err != nil {
		log.Print(err)
		return
	}
	if currRv, err = app.ordersRepository.GetCurrentRv(); err != nil {
		fmt.Println(err)
		return
	}
	if data.Rv >= currRv {
		err = app.ordersRepository.Add(data)
		if err != nil {
			log.Print(err)
			return
		}
	}

	time.Sleep(60 * time.Second)
}
