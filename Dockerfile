FROM golang:1.21.4-alpine AS build
WORKDIR /app
COPY . .
RUN go mod download
RUN go build -o /app/sci_review

FROM golang:1.21.4-alpine
WORKDIR /app
COPY --from=build /app/sci_review .
COPY .env .
COPY templates ./templates
COPY assets ./assets
CMD ["./sci_review"]