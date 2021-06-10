FROM golang
WORKDIR /app
COPY . .
RUN go build .
ENV PORT=6000
EXPOSE ${PORT}
CMD ["/app/edward-backend"]
