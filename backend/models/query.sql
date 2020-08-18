create table if not exists rv_const
(
	rv varchar(8) not null,
	primary key(rv)
);

Create table if not exists orders 
(
	status varchar(4),
	order_id int not null unique,
	store_id int not null,
	date_created timestamp,
	primary key(order_id) 
);
create table if not exists goods
(
	gid int not null unique,
	price int not null,
	status varchar(4),
	chrt_id int not null,
	order_id int not null,
	primary key(gid),
	constraint goods_fk foreign key(order_id)
	references orders (order_id)
);

create or replace function 
insert_into_rv(new_rv varchar(8)) returns void as
$$
begin
	delete from public.rv_const;
	insert into public.rv_const(rv) values (new_rv);
end;
$$ language plpgsql;

--insert into orders
create or replace function
insert_into_orders(new_status varchar(4),
				   new_order_id int,
				   new_store_id int, new_date_created timestamp)
				   returns void as
$$
begin
	insert into orders(status, order_id, store_id, date_created)
	values(new_status, new_order_id, new_store_id, new_date_created)
	on conflict (order_id) do update set 
	status = new_status, store_id = new_store_id, date_created = new_date_created;
end;
$$ language plpgsql;
drop function insert_into_orders(varchar(4), int, int, timestamp)
--insert into goods
create or replace function
insert_into_goods(new_gid int,
				   new_price int,
				   new_status varchar(4),
				   new_chrt_id int,
				   new_order_id int) returns void as
$$
begin
	insert into goods(gid, price, status, chrt_id, order_id)
	values(new_gid, new_price, new_status, new_chrt_id, new_order_id)
	on conflict (gid) do update set 
	gid = new_gid, price = new_price,
	status = new_status, order_id = new_order_id;
end;
$$ language plpgsql;
drop function insert_into_goods(int, int, varchar(4), int, int)

--get orders in page by it's number
create or replace function
get_by_page(page int, n_on_page int) returns table
(
	status varchar(4),
	order_id int,
	store_id int,
	date_created timestamp
)as
$$
begin
	 return query select * from orders order by order_id
	 limit n_on_page offset (page-1)*n_on_page;
end;
$$ language plpgsql;

select insert_into_goods(10, 11, 'u', 12, 362254653)
select insert_into_orders('u', 362254653, 2, '2019-10-15T09:48:42.244522+03:00')
select * from goods
delete from orders
delete from goods
select * from orders
select * from rv_const
	
	
	