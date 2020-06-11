FROM golang:1.14

ADD . /go/src/
WORKDIR /go/src
RUN go get team_5_game
RUN go install team_5_game
ENTRYPOINT /go/bin/team_5_game

EXPOSE 3000
