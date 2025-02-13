FROM alpine:latest  

RUN apk --no-cache add ca-certificates  

COPY link_shortener_lin link_shortener_lin

COPY local.yaml ./config/local.yaml  
COPY .env .env

ENTRYPOINT ["/link_shortener_lin"]  

EXPOSE 44044
