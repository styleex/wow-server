FROM golang:1.19-buster
WORKDIR /goman/

COPY . .

RUN go install -a -trimpath -ldflags='-w -s' ./cmd/...

FROM debian:buster-slim
WORKDIR /apps/

COPY --from=0 /go/bin /apps/
ENV PATH="/apps/:${PATH}"
