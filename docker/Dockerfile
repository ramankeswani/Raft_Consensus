FROM golang:latest 
RUN mkdir -p /app  
WORKDIR /app 
ADD . /app
RUN go build ./main.go ./server.go ./client.go 
CMD ["./main"]
