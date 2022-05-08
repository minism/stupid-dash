FROM golang:alpine

RUN apk add --update gcc musl-dev git

WORKDIR /app

# Deps first for efficiency
COPY go.mod .
COPY go.sum .
RUN go mod download

# Build binary
COPY . .
RUN go build

EXPOSE 80/tcp
ENV PORT 80

CMD ["./stupid-dash"]