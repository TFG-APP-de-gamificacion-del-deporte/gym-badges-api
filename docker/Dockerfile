FROM alpine:3.20
RUN apk --no-cache add ca-certificates

COPY gym-badges-api /bin

EXPOSE 8080

CMD ./bin/gym-badges-api