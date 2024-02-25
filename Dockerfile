FROM golang:1.19-rc-bullseye as builder

WORKDIR /app
COPY . ./


RUN go install golang.org/x/tools/gopls@latest
RUN apt-get install ca-certificates && update-ca-certificates
RUN apt-get install git
RUN curl --proto '=https' --tlsv1.2 -sSf https://just.systems/install.sh | bash -s -- --to /usr/local/bin

RUN mkdir -p -m 0600 ~/.ssh && ssh-keyscan github.com >> ~/.ssh/known_hosts

RUN git config --global url."git@github.com:".insteadOf "https://github.com/"
COPY go.mod .
COPY go.sum .


RUN just build

FROM alpine

ENV SERVER_PORT=3000
WORKDIR /app
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=builder /app/dist/ /bin
COPY --from=builder /usr/bin/git /bin
EXPOSE $SERVER_PORT

ARG release
ENV RELEASE_SHA $release

CMD ["app"]
