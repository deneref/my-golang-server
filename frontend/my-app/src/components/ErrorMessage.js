import React from 'react';

import style from './orders.module.css'
import err_img from './err_img.jpg'

const ErrorMessage = () => {
    return (
        <div className={style.errorMessage}>
            <div> Упс, ошибо4ка вышла - такого заказа нет</div>
            <img className={style.image} src={err_img} alt="error" />
        </div>
    )
}

export default ErrorMessage;