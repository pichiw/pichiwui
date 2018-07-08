FROM scratch

WORKDIR /
ADD server.bin .

WORKDIR /pichiw
ADD app.wasm .
ADD wasm_exec.js .
ADD app.css .
ADD index.html .

CMD ["/server.bin", "/pichiw", "80"]