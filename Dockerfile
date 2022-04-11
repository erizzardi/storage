FROM golang:alpine AS builder
LABEL image=builder

# download necessary packages
RUN apk add git make

ENV GO111MODULE=on

WORKDIR /

# copy everything from project directory
COPY . .

# download dependencies
RUN go mod download


# compile binary
RUN make storage
RUN chmod 111 storage

#-------------------#
FROM scratch

COPY --from=builder /storage /storage

EXPOSE 8081 8081

ENTRYPOINT [ "/storage" ]