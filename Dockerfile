FROM golang:1.21-alpine AS builder

WORKDIR /src
RUN echo $(ls -a)
COPY . .

RUN echo $(ls -a)
RUN go mod tidy
RUN go build -o ./app src/main.go

FROM alpine:latest AS runner
COPY --from=builder src/app /app
EXPOSE 8080
ENTRYPOINT ["/app"]