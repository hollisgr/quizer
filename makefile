SRC := cmd/app/main.go
EXEC := quizer_server

UUID := github.com/google/uuid
GIN := github.com/gin-gonic/gin
PGX := github.com/jackc/pgx github.com/jackc/pgx/v5/pgxpool
CLEANENV := github.com/ilyakaznacheev/cleanenv
CORS := github.com/gin-contrib/cors
WS := github.com/gorilla/websocket
JWT := github.com/golang-jwt/jwt/v5

all: clean build run

build: $(SRC)
	go build -o $(EXEC) $(SRC)

run:
	./$(EXEC)

clean:
	rm -f $(EXEC)

mod:
	go mod init $(EXEC)

get:
	go get \
		$(GIN) \
		$(PGX) \
		$(CLEANENV) \
		$(UUID) \
		$(CORS) \
		$(WS) \
		$(JWT)