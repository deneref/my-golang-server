import React from 'react';
import Good from "./Good";
import style from './orders.module.css'

const Order = ({ order_id, status, store_id, goods }) => {
    return (
        <div className={style.orders}>
            <h2>Заказ №</h2>
            <h1>{order_id}</h1>
            <p>Текущий сатус: {status}</p>
            <p>Заказ в магазине №{store_id}</p >
            <h2>В заказе: </h2>
            <ol>
                {goods.map(good => (
                    <Good
                        key={good.chrt_id}
                        gid={good.gid}
                        price={good.price}
                        status={good.status}
                        chrt_id={good.chrt_id} />
                ))}
            </ol>
        </div>
    )
}

export default Order;