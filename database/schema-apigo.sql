# api_go

create database api_go;

use api_go;

create table tickers
(
    id       int auto_increment
        primary key,
    symbol   varchar(10)            not null,
    value    float(64,2)            not null,
    quota    float(64,2)            null,
    avgPrice float(64,2)            null,
    previousClose float(64,2)       null,
    lastChangePercent float(64,2)   null,
    changeFromAvgPrice float(64,2)  null
);

insert into tickers (symbol, value, quota, avgPrice, previousClose, lastChangePercent, changeFromAvgPrice) values ('ABEV3', 75.28, 30.00, 0,0,0,0);
insert into tickers (symbol, value, quota, avgPrice, previousClose, lastChangePercent, changeFromAvgPrice) values ('B3SA3', 90.22, 25.00, 0,0,0,0);
insert into tickers (symbol, value, quota, avgPrice, previousClose, lastChangePercent, changeFromAvgPrice) values ('BIDI4', 89.51, 10.00, 0,0,0,0);
insert into tickers (symbol, value, quota, avgPrice, previousClose, lastChangePercent, changeFromAvgPrice) values ('BRKM5', 238.13, 30.00, 0,0,0,0);
insert into tickers (symbol, value, quota, avgPrice, previousClose, lastChangePercent, changeFromAvgPrice) values ('CSNA3', 130.55, 25.00, 0,0,0,0);
insert into tickers (symbol, value, quota, avgPrice, previousClose, lastChangePercent, changeFromAvgPrice) values ('ELET3', 80.92, 10.00, 0,0,0,0);
insert into tickers (symbol, value, quota, avgPrice, previousClose, lastChangePercent, changeFromAvgPrice) values ('ELET6', 40.69, 30.00, 0,0,0,0);
insert into tickers (symbol, value, quota, avgPrice, previousClose, lastChangePercent, changeFromAvgPrice) values ('EGIE3', 53.68, 25.00, 0,0,0,0);
insert into tickers (symbol, value, quota, avgPrice, previousClose, lastChangePercent, changeFromAvgPrice) values ('FLRY3', 62.58, 10.00, 0,0,0,0);
insert into tickers (symbol, value, quota, avgPrice, previousClose, lastChangePercent, changeFromAvgPrice) values ('ITUB4', 69.28, 30.00, 0,0,0,0);
insert into tickers (symbol, value, quota, avgPrice, previousClose, lastChangePercent, changeFromAvgPrice) values ('ITSA4', 67.35, 25.00, 0,0,0,0);
insert into tickers (symbol, value, quota, avgPrice, previousClose, lastChangePercent, changeFromAvgPrice) values ('KLBN11', 64.86, 10.00, 0,0,0,0);
insert into tickers (symbol, value, quota, avgPrice, previousClose, lastChangePercent, changeFromAvgPrice) values ('PRIO3', 48.75, 30.00, 0,0,0,0);
insert into tickers (symbol, value, quota, avgPrice, previousClose, lastChangePercent, changeFromAvgPrice) values ('STBP3', 152.02, 25.00, 0,0,0,0);
insert into tickers (symbol, value, quota, avgPrice, previousClose, lastChangePercent, changeFromAvgPrice) values ('SHUL4', 257.17, 10.00, 0,0,0,0);
insert into tickers (symbol, value, quota, avgPrice, previousClose, lastChangePercent, changeFromAvgPrice) values ('TAEE4', 63.78, 30.00, 0,0,0,0);
insert into tickers (symbol, value, quota, avgPrice, previousClose, lastChangePercent, changeFromAvgPrice) values ('TAEE11', 31.12, 25.00, 0,0,0,0);
insert into tickers (symbol, value, quota, avgPrice, previousClose, lastChangePercent, changeFromAvgPrice) values ('TRIS3', 82.70, 10.00, 0,0,0,0);
insert into tickers (symbol, value, quota, avgPrice, previousClose, lastChangePercent, changeFromAvgPrice) values ('WEGE3', 412.78, 10.00, 0,0,0,0);

create table buys
(
    id       int auto_increment
             primary key,
    symbol   varchar(10) not null,
    price    float(64,2) not null,
    quantity int         not null,
    date     char(10)        not null
);
create index buys_symbol_IX
    on buys (symbol);
create index buys_date_IX
    on buys (date);

insert into buys (symbol, price, quantity, date) values ('BIDI4', 72.08, 4, '15/01/2020');
insert into buys (symbol, price, quantity, date) values ('BRKM5', 68.90, 2, '15/01/2020');
insert into buys (symbol, price, quantity, date) values ('ITUB4', 69.28, 2, '15/01/2020');
insert into buys (symbol, price, quantity, date) values ('ITSA4', 67.35, 5, '15/01/2020');
insert into buys (symbol, price, quantity, date) values ('EGIE3', 53.68, 1, '15/01/2020');
insert into buys (symbol, price, quantity, date) values ('FLRY3', 62.58, 2, '15/01/2020');

insert into buys (symbol, price, quantity, date) values ('TAEE4', 63.78, 6, '20/01/2020');
insert into buys (symbol, price, quantity, date) values ('STBP3', 72.72, 9, '20/01/2020');
insert into buys (symbol, price, quantity, date) values ('B3SA3', 90.22, 2, '20/01/2020');
insert into buys (symbol, price, quantity, date) values ('ABEV3', 75.28, 4, '20/01/2020');
insert into buys (symbol, price, quantity, date) values ('WEGE3', 74.44, 2, '20/01/2020');
insert into buys (symbol, price, quantity, date) values ('TRIS3', 82.70, 5, '20/01/2020');
insert into buys (symbol, price, quantity, date) values ('CSNA3', 59.80, 4, '20/01/2020');
insert into buys (symbol, price, quantity, date) values ('ELET3', 80.92, 2, '20/01/2020');

insert into buys (symbol, price, quantity, date) values ('SHUL4', 23.52, 2, '24/01/2020');
insert into buys (symbol, price, quantity, date) values ('KLBN11', 64.86, 3, '24/01/2020');

insert into buys (symbol, price, quantity, date) values ('BIDI4', 17.43, 1, '27/01/2020');
insert into buys (symbol, price, quantity, date) values ('ELET6', 40.69, 1, '27/01/2020');
insert into buys (symbol, price, quantity, date) values ('STBP3', 79.30, 10, '27/01/2020');
insert into buys (symbol, price, quantity, date) values ('BRKM5', 71.88, 2, '27/01/2020');
insert into buys (symbol, price, quantity, date) values ('WEGE3', 82.06, 2, '27/01/2020');
insert into buys (symbol, price, quantity, date) values ('CSNA3', 70.75, 5, '27/01/2020');

insert into buys (symbol, price, quantity, date) values ('TAEE11', 31.12, 1, '31/01/2020');
insert into buys (symbol, price, quantity, date) values ('PRIO3', 48.75, 1, '31/01/2020');
insert into buys (symbol, price, quantity, date) values ('WEGE3', 81.00, 2, '31/01/2020');
insert into buys (symbol, price, quantity, date) values ('WEGE3', 41.00, 1, '31/01/2020');
insert into buys (symbol, price, quantity, date) values ('BRKM5', 97.35, 3, '31/01/2020');

insert into buys (symbol, price, quantity, date) values ('SHUL4', 89.52, 6, '17/02/2020');
insert into buys (symbol, price, quantity, date) values ('SHUL4', 14.89, 1, '17/02/2020');

insert into buys (symbol, price, quantity, date) values ('SHUL4', 129.24, 9, '02/03/2020');
insert into buys (symbol, price, quantity, date) values ('WEGE3', 134.28, 3, '02/03/2020');

insert into buys (symbol, price, quantity, date) values ('ITSA4', 54.30, 5, '11/03/2020');