FROM golang:1.17-alpine as builder

WORKDIR /go/src/gwentgo


COPY go.mod .
COPY go.sum .

# téléchargement des dépendances
RUN go mod download

COPY *.go .

# build Go
RUN go build -ldflags="-s -w" -o /main

#FROM alpine
#COPY --from=builder /main .
EXPOSE 8080
WORKDIR /
RUN adduser -DHu 1000 geralt
USER geralt
ENV GIN_MODE release
ADD cards.csv /
CMD ["/main"]

