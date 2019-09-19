ARG GO_VERSION=1.13

# Stage: Build
FROM golang:${GO_VERSION}-alpine AS builder
RUN mkdir /user && \
    echo 'nobody:x:65534:65534:nobody:/:' > /user/passwd && \
    echo 'nobody:x:65534:' > /user/group
RUN apk add --no-cache ca-certificates git
WORKDIR /src
COPY ./go.mod ./go.sum ./
RUN go mod download
COPY ./ ./
RUN CGO_ENABLED=0 go build \
    -installsuffix 'static' \
    -o /app .

# Stage: Final
FROM scratch AS final
COPY --from=builder /user/group /user/passwd /etc/
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=builder /app /app
COPY --from=builder /src/internal/ /internal/
COPY --from=builder /src/static/ /static/
COPY --from=builder /src/eye-of-sauron.yml /
ADD eye-of-sauron.yml /
EXPOSE 8080

# Unpriviliged user for security reasons.
USER nobody:nobody

# Run the app
ENTRYPOINT ["/app"]
CMD ["start"]