FROM alpine:3.20

RUN apk --no-cache add ca-certificates

WORKDIR /root/

COPY gym-badges-api .

RUN chmod +x gym-badges-api

EXPOSE 8080

CMD ["./gym-badges-api"]