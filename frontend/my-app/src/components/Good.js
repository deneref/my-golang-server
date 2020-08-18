import React from 'react';
import style from './orders.module.css'

const Good = ({ gid, price, status, chrt_id }) => {
    return (
        <div className={style.goods}>
            <p>gid: {gid}</p>
            <p>Цена: {price}</p>
            <p>Статус: {status}</p>
            <p>chrt_id: {chrt_id}</p>
        </div>
    )
}

export default Good;