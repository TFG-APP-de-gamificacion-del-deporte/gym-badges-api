FROM alpine:3.20
RUN apk --no-cache add ca-certificates

ARG BINARY_NAME=main
ENV BINARY_NAME ${BINARY_NAME}

COPY gym-badges-api /bin

EXPOSE 8080

CMD ["./gym-badges-api"]