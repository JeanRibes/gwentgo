FROM golang:1.17-alpine as builder

WORKDIR /go/src/gwentgo

COPY go.mod .
COPY go.sum .

# téléchargement des dépendances
RUN go mod download

COPY *.go ./
COPY gwent ./gwent
COPY ia ./ia

# build Go
RUN go build -ldflags="-s -w" -o /main

FROM alpine:3.14
COPY --from=builder /main .

WORKDIR /
RUN adduser -DHu 1000 geralt
ENV GIN_MODE release

ADD cards.csv /
ADD templates/ /templates
ADD static /static

RUN touch /user_db.json && chown geralt /user_db.json
RUN touch /db.json && chown geralt /db.json

USER geralt
EXPOSE 8080
CMD ["/main"]
