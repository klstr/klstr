FROM golang:1.10 as builder

RUN curl -fsSL https://github.com/Masterminds/glide/releases/download/v0.13.1/glide-v0.13.1-linux-amd64.tar.gz -o glide.tar.gz \
    && tar -zxf glide.tar.gz \
    && mv linux-amd64/glide /usr/bin/ \
    && rm -r linux-amd64 \
    && rm glide.tar.gz

WORKDIR /go/src/github.com/klstr/klstr
COPY . .
RUN glide up -v
RUN make

FROM alpine:latest
RUN apk --no-cache add ca-certificates
WORKDIR /root/
EXPOSE 3000
COPY --from=builder /go/src/github.com/klstr/klstr/klstr .
CMD ["./klstr"]