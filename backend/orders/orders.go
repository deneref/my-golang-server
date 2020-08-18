package orders

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"strconv"
	"time"

	"my-golang-server/backend/models"
)

type Orders interface {
	GetCurrentRv() (int, error)
	Add(order *models.Order) error
	List() (models.Order, error)
	GetAmountOfOrders() (int, error) //возвращает количество заказов в базе
	CloseDB() error
	GetByOrderID(int) (models.Order, error)      // возвращает конкретный заказ по айдишнику
	GetByPageNum(int, int) (models.Order, error) // возвращает страницу с заказами по возрастанию айдишника
}

type orders struct {
	db *sql.DB
}

func NewOrdersRepository(db *sql.DB, conf models.Config) Orders {
	dbinfo := fmt.Sprintf("host=%s user=%s password=%s dbname=%s sslmode=disable",
		conf.SQLDataBase.Server, conf.SQLDataBase.UserID, conf.SQLDataBase.Password, conf.SQLDataBase.Database)

	db, err := sql.Open("postgres", dbinfo)
	if err != nil {
		return nil
	}
	log.Println("connected succesfully")
	return &orders{db: db}
}

func (r *orders) CloseDB() error {
	return r.db.Close()
}

func (r *orders) Add(order *models.Order) error {

	tx, err := r.db.Begin()
	if err != nil {
		return err
	}
	defer func() {
		if err != nil {
			tx.Rollback()
			return
		}
		err = tx.Commit()
	}()

	//rewrites rv regardless
	var stmt *sql.Stmt
	if stmt, err = r.db.Prepare("SELECT insert_into_rv($1)"); err != nil {
		return err
	}
	if _, err := stmt.Exec(order.Rv); err != nil {
		//когда в бд rv интовый здесь ошибка 'rv тип инт, а выражение varchar'
		return err
	}
	//functions insert or update on conflict (order_id)
	stmt_c, err := r.db.Prepare("SELECT insert_into_orders($1,$2,$3,$4)")
	if err != nil {
		return err
	}
	//upd on conflict on (chrt_id)
	stmt_g, err := r.db.Prepare("SELECT insert_into_goods($1,$2,$3,$4,$5)")
	if err != nil {
		return err
	}
	for _, i := range order.Content {
		_, err = stmt_c.Exec(i.Status, i.OrderID, i.StoreID, i.DateCreated)
		if err != nil {
			return err
		}
		for _, j := range i.Goods {
			_, err = stmt_g.Exec(j.Gid, j.Price, j.Status, j.ChrtID, i.OrderID)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func (r *orders) List() (models.Order, error) {
	var result models.Order
	var err error
	if result.Rv, err = r.GetCurrentRv(); err != nil {
		return result, err
	}

	row_ord, err := r.db.Query("SELECT status, order_id, store_id, date_created FROM orders ORDER BY order_id")
	if err != nil {
		return result, err
	}
	for row_ord.Next() {
		curr_ord := new(models.Content)
		if err = row_ord.Scan(&curr_ord.Status, &curr_ord.OrderID, &curr_ord.StoreID, &curr_ord.DateCreated); err != nil {
			return result, err
		}
		row_goods, err := r.db.Query("SELECT gid, price, status, chrt_id FROM goods WHERE order_id = $1", curr_ord.OrderID)
		if err != nil {
			return result, err
		}
		for row_goods.Next() {
			curr_good := new(models.Goods)
			if err = row_goods.Scan(&curr_good.Gid, &curr_good.Price, &curr_good.Status, &curr_good.ChrtID); err != nil {
				return result, err
			}
			curr_ord.Goods = append(curr_ord.Goods, *curr_good)
		}
		result.Content = append(result.Content, *curr_ord)
	}
	return result, nil
}

func (r *orders) GetByOrderID(id int) (models.Order, error) {
	var result models.Order
	var err error
	if result.Rv, err = r.GetCurrentRv(); err != nil {
		return result, err
	}

	row_ord, err := r.db.Query("SELECT status, order_id, store_id, date_created FROM orders WHERE order_id=$1 ORDER BY order_id", id)
	if err != nil {
		return result, err
	}
	for row_ord.Next() {
		curr_ord := new(models.Content)
		if err = row_ord.Scan(&curr_ord.Status, &curr_ord.OrderID, &curr_ord.StoreID, &curr_ord.DateCreated); err != nil {
			return result, err
		}
		row_goods, err := r.db.Query("SELECT gid, price, status, chrt_id FROM goods WHERE order_id = $1", id)
		if err != nil {
			return result, err
		}
		for row_goods.Next() {
			curr_good := new(models.Goods)
			if err = row_goods.Scan(&curr_good.Gid, &curr_good.Price, &curr_good.Status, &curr_good.ChrtID); err != nil {
				return result, err
			}
			curr_ord.Goods = append(curr_ord.Goods, *curr_good)
		}
		result.Content = append(result.Content, *curr_ord)
	}
	return result, nil
}

func (r *orders) GetByPageNum(page, on_page int) (models.Order, error) {
	var result models.Order
	var err error
	if result.Rv, err = r.GetCurrentRv(); err != nil {
		return result, err
	}

	row_ord, err := r.db.Query("SELECT status, order_id, store_id, date_created FROM get_by_page($1, $2)", page, on_page)
	if err != nil {
		return result, err
	}
	for row_ord.Next() {
		curr_ord := new(models.Content)
		if err = row_ord.Scan(&curr_ord.Status, &curr_ord.OrderID, &curr_ord.StoreID, &curr_ord.DateCreated); err != nil {
			return result, err
		}
		row_goods, err := r.db.Query("SELECT gid, price, status, chrt_id FROM goods WHERE order_id = $1", curr_ord.OrderID)
		if err != nil {
			return result, err
		}
		for row_goods.Next() {
			curr_good := new(models.Goods)
			if err = row_goods.Scan(&curr_good.Gid, &curr_good.Price, &curr_good.Status, &curr_good.ChrtID); err != nil {
				return result, err
			}
			curr_ord.Goods = append(curr_ord.Goods, *curr_good)
		}
		result.Content = append(result.Content, *curr_ord)
	}
	return result, nil
}

func (r *orders) GetCurrentRv() (int, error) {
	var rv string
	err := r.db.QueryRow("select rv from public.rv_const").Scan(&rv)
	if err != nil {
		fmt.Println(err)
		return -1, err
	}
	res, err := strconv.Atoi(rv)
	return res, err
}

func (r *orders) GetAmountOfOrders() (int, error) {
	row, err := r.db.Query("SELECT count(order_id) FROM orders")
	if err != nil {
		return -1, err
	}
	var res int
	for row.Next() {
		if err = row.Scan(&res); err != nil {
			return -1, err
		}
	}
	return res, nil
}

func MakeRequest(url string) (*models.Order, error) {
	client := &http.Client{
		Transport: &http.Transport{
			Dial: (&net.Dialer{
				Timeout:   30 * time.Second,
				KeepAlive: 30 * time.Second,
			}).Dial,
		},
	}
	resp, err := client.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	defer client.CloseIdleConnections()
	var result models.Order
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal(body, &result)
	if err != nil {
		return nil, err
	}

	return &result, nil
}
