FROM golang:1.13 AS build
WORKDIR /usr/src/app

# Копируем исходный код в Docker-контейнер
COPY go.mod .
COPY go.sum .
RUN go mod download

# Build project
COPY . .
RUN go build main.go

FROM ubuntu:18.10

# Установка postgresql
ENV PGVER 10
RUN apt -y update && apt install -y postgresql-$PGVER

# Run the rest of the commands as the ``postgres`` user created by the ``postgres-$PGVER`` package when it was ``apt-get installed``
USER postgres

WORKDIR app
COPY --from=build /usr/src/app/main .
COPY --from=build /usr/src/app/database database

# Create a PostgreSQL role named ``docker`` with ``docker`` as the password and
# then create a database `docker` owned by the ``docker`` role.
RUN /etc/init.d/postgresql start &&\
    psql --command "CREATE USER docker WITH SUPERUSER PASSWORD 'docker';" &&\
    createdb -O docker docker &&\
    psql --file=database/database.sql &&\
    /etc/init.d/postgresql stop

# Adjust PostgreSQL configuration so that remote connections to the
# database are possible.
RUN echo "host all  all    0.0.0.0/0  md5" >> /etc/postgresql/$PGVER/main/pg_hba.conf

# Expose the PostgreSQL port
EXPOSE 5432

# Add VOLUMEs to allow backup of config, logs and databases
VOLUME  ["/etc/postgresql", "/var/log/postgresql", "/var/lib/postgresql"]

# Back to the root user
USER root

# Объявлем порт сервера
EXPOSE 5000

# Запускаем PostgreSQL и сервер
CMD service postgresql start && ./main --scheme=http --port=5000 --host=0.0.0.0 --database=postgres://docker:docker@localhost/docker
