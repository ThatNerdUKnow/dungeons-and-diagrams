FROM golang
RUN apt-get update && apt-get install -y libz3-dev libz3-4 z3
COPY . .
RUN go install .
CMD [ "dungeons-and-diagrams" ]