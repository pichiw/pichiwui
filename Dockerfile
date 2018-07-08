FROM alpine:latest

RUN apk update && apk add ca-certificates && rm -rf /var/cache/apk/

WORKDIR /
ADD caddy .
ADD Caddyfile .

WORKDIR /pichiw
ADD app.wasm .
ADD wasm_exec.js .
ADD app.css .
ADD index.html .

CMD ["/caddy", "-conf", "/Caddyfile"] 