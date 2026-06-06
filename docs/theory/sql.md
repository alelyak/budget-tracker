пример: select * from courier_actions order by order_id desc limit 20
сортировка: desc - по убыванию, без desc - по возрастанию
limit - первые 20

2005-08-22 06:37:58 - как используется дата и время на самом деле, но в redash показывает по другому
select creation_time::text from orders - преобразует как надо
select (10/0.7)::text - как надо, а так по умолчанию две цифры после запятой
select (10/0.7, 2) надо писать так

правильный порядок
SELECT -- перечисление полей результирующей таблицы
FROM -- указание источника данных
WHERE -- фильтрация данных
GROUP BY -- группировка данных
HAVING -- фильтрация данных после группировки
ORDER BY -- сортировка результирующей таблицы
LIMIT -- ограничение количества выводимых записей

Приоритеты выполнения операций:
умножение и деление (* и /)
сложение и вычитание (+ и -)
операторы сравнения (=, !=, >, <, >=, <=)
NOT
AND
OR


При составлении SQL-запросов колонкам в результирующей таблице можно присваивать любые другие имена (их ещё называют «алиасами»). Это можно делать с помощью оператора AS:
SELECT name AS new_name
FROM table

# функции
Используя операторы SELECT, FROM, ORDER BY и LIMIT, а также функцию LENGTH, определите товар с самым длинным названием в таблице products. Выведите его наименование, длину наименования в символах, а также цену этого товара. Колонку с длиной наименования в символах назовите name_length.

select name, length(name) as name_length, price from products order by name_length desc limit 1

Примените последовательно функции UPPER и SPLIT_PART к колонке name и преобразуйте наименования товаров в таблице products так, чтобы от названий осталось только первое слово, записанное в верхнем регистре. Колонку с новым названием, состоящим из первого слова, назовите first_word.
В результат включите исходные наименования товаров, новые наименования из первого слова, а также цену товаров. Результат отсортируйте по возрастанию исходного наименования товара в колонке name.

select name, split_part(upper(name), ' ', 1) as first_word, price from products order by name

Измените тип колонки price из таблицы products на VARCHAR. Выведите колонки с наименованием товаров, ценой в исходном формате и ценой в формате VARCHAR. Новую колонку с ценой в новом формате назовите price_char. Результат отсортируйте по возрастанию наименования товара в колонке name. Количество выводимых записей не ограничивайте.

select name, price, cast(price as VARCHAR) as price_char from products order by name

Для первых 200 записей из таблицы orders выведите информацию в следующем виде (обратите внимание на пробелы): Заказ № [id заказа] создан [дата]. Полученную колонку назовите order_info.

select concat('Заказ № ',order_id,' создан ', creation_time::date) as order_info from orders limit 200

Выведите id всех курьеров и их годы рождения из таблицы couriers.

Год рождения необходимо получить из колонки birth_date. Новую колонку с годом назовите birth_year. Результат отсортируйте сначала по убыванию года рождения курьера (т.е. от самых младших к самым старшим), затем по возрастанию id курьера.

select courier_id, date_part('year', birth_date) as birth_year from couriers order by birth_year desc, courier_id

Как и в предыдущем задании, снова выведите id всех курьеров и их годы рождения, только теперь к извлеченному году примените функцию COALESCE. Укажите параметры функции так, чтобы вместо NULL значений в результат попадало текстовое значение unknown. Названия полей оставьте прежними.Отсортируйте итоговую таблицу сначала по убыванию года рождения курьера, затем по возрастанию id курьера.

select courier_id, coalesce(date_part('year', birth_date)::varchar, 'unknown') as birth_year from couriers order by birth_year desc, courier_id

Давайте представим, что по какой-то необъяснимой причине мы вдруг решили в одночасье повысить цену всех товаров в таблице products на 5%. Выведите id и наименования всех товаров, их старую и новую цену. Колонку со старой ценой назовите old_price, а колонку с новой — new_price. Результат отсортируйте сначала по убыванию новой цены, затем по возрастанию id товара.

select product_id, name, price as old_price, price*1.05 as new_price from products order by new_price desc, product_id

Вновь, как и в прошлом задании, повысьте цену всех товаров на 5%, только теперь к колонке с новой ценой примените функцию ROUND. Выведите id и наименования товаров, их старую цену, а также новую цену с округлением. Новую цену округлите до одного знака после запятой, но тип данных не меняйте.Результат отсортируйте сначала по убыванию новой цены, затем по возрастанию id товара.

select product_id, name, price as old_price, round(price*1.05, 1) as new_price from products order by new_price desc, product_id

Повысьте цену на 5% только на те товары, цена которых превышает 100 рублей. Цену остальных товаров оставьте без изменений. Также не повышайте цену на икру, которая и так стоит 800 рублей. Выведите id и наименования всех товаров, их старую и новую цену. Цену округлять не нужно. Результат отсортируйте сначала по убыванию новой цены, затем по возрастанию id товара.

select product_id, name, price as old_price,
case when price>100 and name!='икра' then price*1.05
else price
end as new_price
from products
order by new_price desc, product_id

Вычислите НДС каждого товара в таблице products и рассчитайте цену без учёта НДС. Выведите всю информацию о товарах, включая сумму налога и цену без его учёта. Колонки с суммой налога и ценой без НДС назовите соответственно tax и price_before_tax. Округлите значения в этих колонках до двух знаков после запятой. Результат отсортируйте сначала по убыванию цены товара без учёта НДС, затем по возрастанию id товара.

SELECT product_id, name, price, ROUND(price-price/1.2, 2) AS tax, ROUND(price/1.2, 2) AS price_before_tax
FROM products
ORDER BY price_before_tax desc, product_id

# WHERE/WHERE
Напишите SQL-запрос к таблице products и выведите всю информацию о товарах, цена которых не превышает 100 рублей. Результат отсортируйте по возрастанию id товара.

select product_id, name, price from products where price<=100 order by product_id

Отберите пользователей женского пола (female) из таблицы users. Выведите только id этих пользователей. Результат отсортируйте по возрастанию id.
Добавьте в запрос оператор LIMIT и выведите только 1000 первых id из отсортированного списка.

select user_id from users where sex='female' order by user_id limit 1000

Отберите из таблицы user_actions все действия пользователей по созданию заказов, которые были совершены ими после полуночи 6 сентября 2022 года. Выведите колонки с id пользователей, id созданных заказов и временем их создания. Результат должен быть отсортирован по возрастанию id заказа.

select user_id, order_id, time from user_actions where time > '2022-09-06 00:00:00' and action = 'create_order' order by order_id

Назначьте скидку 20% на все товары из таблицы products и отберите те, цена на которые с учётом скидки превышает 100 рублей. Выведите id товаров, их наименования, прежнюю цену и новую цену с учётом скидки. Колонку со старой ценой назовите old_price, с новой — new_price. Результат должен быть отсортирован по возрастанию id товара.

select product_id, name, price as old_price, price*0.8 as new_price from products where price*0.8>100 order by product_id
! псевдоним только в order by, до этого нет!

Отберите из таблицы products все товары, названия которых либо начинаются со слова «чай», либо состоят из пяти символов. Выведите две колонки: id товаров и их наименования.
Результат должен быть отсортирован по возрастанию id товара.

select product_id, name from products where split_part(name, ' ', 1)='чай' or length(name)=5 order by product_id

Отберите из таблицы products все товары, содержащие в своём названии последовательность символов «чай» (без кавычек). Выведите две колонки: id продукта и его название.
Результат должен быть отсортирован по возрастанию id товара.

select product_id, name from products where name like '%чай%' order by product_id

Выберите из таблицы products id и наименования только тех товаров, названия которых начинаются на букву «с» и содержат только одно слово.
Результат должен быть отсортирован по возрастанию id товара.

select product_id, name from products where name like 'с%' and name not like '% %'  order by product_id

Составьте SQL-запрос, который выбирает из таблицы products все чаи стоимостью больше 60 рублей и вычисляет для них цену со скидкой 25%.
Скидку в % менеджер попросил указать в отдельном столбце в формате текста, то есть вот так: «25%» (без кавычек). Столбцы со скидкой и новой ценой назовите соответственно discount и new_price.
Также необходимо любым известным способом избавиться от «чайного гриба»: вряд ли менеджер имел в виду и его, когда ставил нам задачу.
Результат должен быть отсортирован по возрастанию id товара.

SELECT product_id,
       name,
       price,
       '25%' as discount,
       price*0.75 as new_price
FROM   products
WHERE  name like '%чай%'
   and price > 60
   and name not like 'чайный гриб'
ORDER BY product_id

Из таблицы user_actions выведите всю информацию о действиях пользователей с id 170, 200 и 230 за период с 25 августа по 4 сентября 2022 года включительно. Результат отсортируйте по убыванию id заказа — то есть от самых поздних действий к самым первым.

select user_id, order_id, action, time from user_actions where user_id in (170,200,230) and time between '2022-08-25' and '2022-09-05' order by order_id desc

Напишите SQL-запрос к таблице couriers и выведите всю информацию о курьерах, у которых не указан их день рождения.
Результат должен быть отсортирован по возрастанию id курьера.

Напишите SQL-запрос к таблице couriers и выведите всю информацию о курьерах, у которых не указан их день рождения.
Результат должен быть отсортирован по возрастанию id курьера.

select birth_date, courier_id, sex from couriers where birth_date is NULL order by courier_id

Определите id и даты рождения 50 самых молодых пользователей мужского пола из таблицы users. Не учитывайте тех пользователей, у которых не указана дата рождения.

select user_id, birth_date from users where sex='male' and birth_date is not null order by birth_date desc limit 50

Напишите SQL-запрос к таблице courier_actions, чтобы узнать id и время доставки последних 10 заказов, доставленных курьером с id 100.

select order_id, time from courier_actions where courier_id=100 and action='deliver_order' order by order_id desc limit 10


Из таблицы user_actions получите id всех заказов, сделанных пользователями сервиса в августе 2022 года.
Результат отсортируйте по возрастанию id заказа.

select order_id from user_actions where action='create_order' and date_part('month', time)=08 and date_part('year', time)=2022 order by order_id

Из таблицы couriers отберите id всех курьеров, родившихся в период с 1990 по 1995 год включительно.
Результат отсортируйте по возрастанию id курьера.

select courier_id from couriers where date_part('year', birth_date) between 1990 and 1995 order by courier_id

Из таблицы user_actions получите информацию о всех отменах заказов, которые пользователи совершали в течение августа 2022 года по средам с 12:00 до 15:59:59.
Результат отсортируйте по убыванию id отменённых заказов.

SELECT user_id,
       order_id,
       action,
       time
FROM   user_actions
WHERE  action = 'cancel_order'
   and date_part('year', time) = 2022
   and date_part('month', time) = 8
   and date_part('dow', time) = 3
   and date_part('hour', time) between 12
   and 15
ORDER BY order_id desc
