FROM golang:1.16-alpine

WORKDIR /app

RUN go env -w GOPROXY=https://goproxy.io
RUN go get github.com/containerd/cgroups
COPY go.mod ./
COPY go.sum ./
RUN go mod download
COPY . .
RUN go build -o /start local/local_testing.go

CMD [ "/start" ]


