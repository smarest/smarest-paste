FROM golang:1.13

WORKDIR /go/
COPY ./main.go ./

# Download all the dependencies
RUN go get -d github.com/gin-gonic/gin
RUN go get -d github.com/smarest/smarest-common/service
RUN go get -d github.com/smarest/smarest-common/domain/entity/exception
RUN go get -d github.com/smarest/smarest-paste/application

# Install the package
RUN go build -o main .

# This container exposes port 8080 to the outside world
EXPOSE 8080

# Run the executable
CMD ["./main"]
