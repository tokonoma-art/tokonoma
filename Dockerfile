# Build the server
FROM golang:latest AS server-builder
WORKDIR /app/server
# Download dependencies first (to allow good caching)
COPY server/go.mod server/go.sum ./
RUN go mod download
# Copy the sources
COPY server/. .
# Build the app and include every single dependencies (see https://rollout.io/blog/building-minimal-docker-containers-for-go-applications/)
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main .

# Build the client
FROM node:latest AS client-builder
WORKDIR /app/client
# Download dependencies first (to allow good caching)
RUN npm install webpack webpack-cli -g
COPY client/package.json client/package-lock.json ./
RUN npm config set registry http://registry.npmjs.org/ && npm install
# Copy the sources
COPY client/. .
# Build the app
RUN webpack

# Prepare final container
FROM scratch
# Ensure root certificates are there
COPY --from=server-builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
# Copy server & client
COPY --from=server-builder /app/server/main /app/server/main
COPY --from=client-builder /app/client/dist /app/client/dist
# Run the server on start
CMD ["/app/server/main", "--app", "/app", "--storage", "/var/lib/tokonoma"]
