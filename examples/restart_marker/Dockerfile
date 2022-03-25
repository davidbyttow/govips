FROM debian:buster

RUN apt-get -y update && \
  apt-get -y install build-essential libtiff5-dev libpng-dev \
  libfftw3-dev librsvg2-dev libgif-dev libexif-dev \
  libexpat1-dev pkg-config glib2.0-dev libimagequant-dev libgsf-1-dev \
  liborc-0.4-dev liblcms2-dev autoconf nasm libtool automake wget

ENV VIPS_VERSION=8.12.2
ARG VIPS_URL=https://github.com/libvips/libvips/releases/download

ENV MOZJPEG_VERSION=3.3.1
ARG MOZJPEG_URL=https://github.com/mozilla/mozjpeg/archive

WORKDIR /tmp
ENV LD_LIBRARY_PATH /lib:/usr/lib:/usr/local/lib

# get vips and mozjpeg
ADD ${MOZJPEG_URL}/v${MOZJPEG_VERSION}.tar.gz \
  ${VIPS_URL}/v${VIPS_VERSION}/vips-${VIPS_VERSION}.tar.gz \
  /tmp/

# compile mozjpeg
RUN tar xf v${MOZJPEG_VERSION}.tar.gz \
  && cd mozjpeg-${MOZJPEG_VERSION} \
  && aclocal \
  && autoconf \
  && autoheader \
  && libtoolize \
  && automake --add-missing \
  && ./configure \
  && make -j4 V=0 \
  && make install libdir=/usr/lib/x86_64-linux-gnu prefix=/usr \
  && cp jpegint.h /usr/include/jpegint.h \
  && cd .. \
  && rm -rf mozjpeg

# compile libvips
RUN tar xf vips-${VIPS_VERSION}.tar.gz \
  && cd vips-${VIPS_VERSION} \
  && ./configure --without-magick --without-pangoft2 --without-ppm \
       --without-analyze --without-radiance --without-OpenEXR \
       --with-jpeg-includes=/opt/mozjpeg/include \
       --with-jpeg-libraries=/opt/mozjpeg/lib64 \
  && make -j4 V=0 \
  && make install \
  && rm -rf /tmp/*

ENV GOLANG_VERSION 1.17.8

RUN set -eux; \
	url="https://golang.org/dl/go${GOLANG_VERSION}.linux-amd64.tar.gz"; \
	wget -O go.tgz "$url"; \
	tar -C /usr/local -xzf go.tgz; \
	rm go.tgz; \
	export PATH="/usr/local/go/bin:$PATH"; \
	go version

ENV GOPATH /go
ENV PATH $GOPATH/bin:/usr/local/go/bin:$PATH
RUN mkdir -p "$GOPATH/src" "$GOPATH/bin" "$GOPATH/pkg" && chmod -R 777 "$GOPATH"

WORKDIR /restart_marker

COPY examples/restart_marker/main.go ./
COPY examples/restart_marker/go.mod ./
COPY resources/jpg-24bit-icc-adobe-rgb.jpg ./

RUN go mod vendor
# Rewrite GOVIPS source code
COPY vips ./vendor/github.com/davidbyttow/govips/v2/vips
RUN go build -o jpeg_rst_mrk ./main.go
RUN mkdir "result"

ENTRYPOINT ["./jpeg_rst_mrk"]

