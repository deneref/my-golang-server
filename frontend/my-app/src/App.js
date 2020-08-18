import React, { useEffect, useState } from 'react';
import Order from "./components/Order";
import Head from "./components/Head";
import ErrorMessage from "./components/ErrorMessage";
import './App.css';

function App() {
  const [orders, setOrders] = useState([]);
  const [search, setSearch] = useState("");
  const [query, setQuery] = useState("getPage?page=1");
  const [contentError, setContentError] = useState(false);

  useEffect(() => {
    getOrders();
  }, [query]);

  const getOrders = async () => {
    try {
      const response = await fetch(
        `http://localhost:8147/${query}`
      );
      if (!response.ok) {
        throw new Error('Такого заказа нет');
      }
      let data = await response.json();
      setOrders(data.content)
      setContentError(false)
    } catch (error) {
      console.log('Не нашли такой заказ');
      setContentError(true)
    }
  }

  const saveXML = async () => {
    try {
      const response = await fetch(
        `http://localhost:8147/allOrdersXML`
      );
      if (!response.ok) {
        throw new Error('Server error');
      }
      let data = response.text()
        .then(data => {
          var FileSaver = require('file-saver');
          var blob = new Blob([data], { type: "text/xml;charset=utf-8" });
          FileSaver.saveAs(blob, "myOrdersXML.xml");
        });
      console.log(data)
    } catch (error) {
      console.log('Server did not response');
      console.log(error)
    }
  }

  const updateSearch = e => {
    setSearch(e.target.value)
  }

  const getSearch = e => {
    e.preventDefault();
    setOrders([])
    setQuery("getOrder?id=" + search)
  }

  return (
    <div className="App">
      <Head />
      <form onSubmit={getSearch} className="search-form">
        <input className="search-bar" type="text"
          value={search} onChange={updateSearch} />
        <button className="search-button" type="submit">
          Поиск
        </button>
        <button className="search-button" onClick={saveXML} type="button">
          Загрузить данные в формате XML
        </button>
      </form>
      <div>
        {
          (!contentError) ?
            <div className="orders">
              {orders.map(order => (
                <Order
                  key={order.order_id}
                  order_id={order.order_id}
                  status={order.status}
                  store_id={order.store_id}
                  goods={order.goods}
                />
              ))}
            </div> :
            <div>
              <ErrorMessage />
            </div>
        }
      </div>
    </div>
  );
}

export default App;
