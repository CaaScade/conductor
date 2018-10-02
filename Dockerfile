FROM golang:1.8

WORKDIR ./conductor

COPY . .

CMD "revel run"
#ENTRYPOINT ["/"]