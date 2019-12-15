# DataBase Project

This is Go project realised with PostgreSQL DB

[Swagger docs](https://tech-db-forum.bozaro.ru/)

### Run
```
docker build -t db_go https://github.com/Ledka17/TP_DB.git
docker run -p 5000:5000 -p 5432:5432--name db_go -t db_go
```