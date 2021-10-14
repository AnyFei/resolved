FROM golang:1.16
RUN mkdir /app
ADD . /app
WORKDIR /app
RUN go build -o Resolved ./src/Resolved
CMD ["/app/Resolved"]
EXPOSE 8080
