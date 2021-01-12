FROM golang:1.14.3-alpine as builder

RUN mkdir -p /home/script
WORKDIR /home/script

COPY go.mod . 
COPY go.sum .

# get dependencies
RUN go mod download

# copy the source code
COPY . .

RUN go build -o shuffler .

######## Start a new stage from a minimal alpine image #######
FROM alpine:latest

# Copy the Pre-built binary file from the previous stage
COPY --from=builder /home/script/shuffler /home/script/shuffler
COPY --from=builder /home/script/email_details.yaml /home/script/email_details.yaml

CMD ["/home/script/shuffler", "--yamlPath", "/home/script/email_details.yaml"]