FROM docker-repo.wixpress.com/com.wixpress.vips-base-image-from-source as builder

COPY . /govips
WORKDIR /govips

ENV GOOS linux
ENV PATH $GOPATH/bin:$GOROOT/bin:$PATH
ENV GOPATH /deps

RUN mkdir -p ~/.ssh/
RUN git config --global url."git@github.com:".insteadOf https://github.com/
RUN ssh-keyscan github.com >> ~/.ssh/known_hosts
RUN mv deps /deps
RUN go test -race -v ./...
