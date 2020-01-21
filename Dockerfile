FROM erdiugurlu/golang-alpine-git:1.13.6 AS builder
 
LABEL MAINTAINER=erdi.ugurlu@isbank.com.tr
 
RUN mkdir /user && \
    echo 'nobody:x:65534:65534:nobody:/:' > /user/passwd && \
    echo 'nobody:x:65534:' > /user/group
 
RUN apk add --update --no-cache ca-certificates
  
 
# Set the working directory outside $GOPATH to enable the support for modules.
#WORKDIR /src
WORKDIR /go/src/restapi

ENV GO111MODULE=on
COPY ./go.mod ./go.sum ./
COPY ./airportresponse.json ./
RUN go mod download

COPY countryairportlist.go .
 
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o app .


FROM golang:1.13.6-alpine AS final
# Import the user and group files from the first stage.
COPY --from=builder /user/group /user/passwd /etc/
WORKDIR /opt/
COPY --from=builder /go/src/restapi/app .
COPY ./airportresponse.json ./
RUN echo '192.168.64.2  country-api.info' >> /etc/resolv.conf
#USER nobody:nobody
EXPOSE 8000
#CMD ["nobody"]
CMD ["./app"]
