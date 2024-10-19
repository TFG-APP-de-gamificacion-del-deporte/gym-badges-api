FROM alpine:3.20
RUN apk add --no-cache libc6-compat

ARG BINARY_NAME=main
ENV BINARY_NAME ${BINARY_NAME}

EXPOSE 8080

COPY gym-badges-api /bin

CMD gym-badges-api --port 8080