FROM golang:1.21.4-alpine AS build
WORKDIR /
COPY . .
RUN go mod download
RUN go build -o /sci_review

FROM golang:1.21.4-alpine
WORKDIR /
COPY --from=build /sci_review .
COPY .env .
COPY templates ./templates
COPY assets ./assets
COPY db/migrations /db/migrations
CMD ["./sci_review"]