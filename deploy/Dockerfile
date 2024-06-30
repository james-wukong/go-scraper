FROM golang:latest

ENV GO111MODULE=ON
ENV GOFLAGS=-mod=vendor

ENV APP_HOME /go/src/app
RUN mkdir -p "$APP_HOME"

WORKDIR "$APP_HOME"

COPY go.mod go.sum ./
RUN go mod download && go mod vendor && go mod verify

# COPY . .

EXPOSE 9001

# CMD ["air", "-c", ".air.toml"]
# CMD ["make", "serve"]
ENTRYPOINT ["tail", "-f", "/dev/null"]