#env GOOS=linux GOARCH=amd64 go build
#docker build -t mercurius:simple .
#docker run -p 8080:8080 -d mercurius:simple

FROM scratch

ADD simple /
ADD conf/ /conf
ADD public/ /public
ADD locale/ /locale

CMD [ "/simple" ]