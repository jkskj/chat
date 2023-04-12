FROM golang:1.19
WORKDIR /app
COPY  . .
ENV GOPROXY="https://goproxy.cn"
RUN go mod download
RUN go build -o main ./main.go
EXPOSE 3000
CMD ["/app/main"]