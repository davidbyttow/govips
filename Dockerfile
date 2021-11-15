FROM docker-repo.wixpress.com/com.wixpress.vips-base-image-from-source as builder

COPY . /govips
WORKDIR /govips

ENV GOOS linux
ENV PATH $GOPATH/bin:$GOROOT/bin:$PATH

# no private deps, so go get should suffice
RUN go get ./....

RUN go test -race -v ./...
