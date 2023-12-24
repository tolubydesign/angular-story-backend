# FROM postgres
# ENV POSTGRES_PASSWORD docker
# ENV POSTGRES_DB world
# COPY world.sql /docker-entrypoint-initdb.d/

# Run 
# $ docker build -t my-postgres-db ./
# $ docker images -a
# $ docker run -d --name my-postgresdb-container -p 5432:5432 my-postgres-db

# https://docs.aws.amazon.com/lambda/latest/dg/go-image.html
FROM golang:1.21 as build
WORKDIR /helloworld
# Copy dependencies list
COPY go.mod go.sum ./
COPY .env .
COPY .env.development ./
COPY .env.production ./
# Build
COPY /app .
COPY /functions .
COPY /init_database .
COPY /cdk/cdk-golang.go ./main.go

RUN go mod download
RUN go mod tidy
RUN go build -o main main.go
# GOOS=linux GOARCH=amd64 CGO_ENABLED=0
# Copy artifacts to a clean image
FROM alpine:3.16
COPY --from=build /helloworld/main /main
ENTRYPOINT [ "/main" ]