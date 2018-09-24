FROM golang:latest 
RUN mkdir -p /app  
WORKDIR /app 
ADD . /app
RUN go get github.com/joho/godotenv github.com/mattn/go-sqlite3 github.com/op/go-logging
RUN go build ./main.go ./server.go ./client.go ./db.go
CMD ["./main 5000"]