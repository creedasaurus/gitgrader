FROM alpine:3.7
COPY build/server /
EXPOSE 1234
CMD [ "/server" ]