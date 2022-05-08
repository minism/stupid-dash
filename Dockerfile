FROM golang:alpine

RUN apk add --update gcc musl-dev git

WORKDIR /app

COPY . .

RUN go build

EXPOSE 80/tcp
ENV PORT 80

CMD ["./stupid-dash"]