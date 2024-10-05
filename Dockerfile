FROM amd64/ubuntu:latest AS base
RUN apt update -y
RUN apt install -y autoconf automake libtool pkg-config build-essential
RUN apt install -y \
	zlib1g-dev \
	libsharp-dev \
	swig \
	wget \
	gfortran \
	git
# Install cfitsio
COPY ./cfitsio_latest.tar.gz /tmp/cfitsio.tar.gz
RUN tar -C /opt -xvf /tmp/cfitsio.tar.gz
WORKDIR /opt/cfitsio-4.5.0
RUN ./configure FC=gfortran FFLAGS="-m64 -fPIC" CFLAGS="-m64 -fPIC"
RUN make
RUN make install
#
FROM base AS healpix
# Install Healpix
WORKDIR /opt/healpix
RUN git clone https://github.com/healpy/healpixmirror.git
WORKDIR /opt/healpix/healpixmirror
RUN cd src/common_libraries/libsharp && autoreconf -i && cd -
RUN cd src/cxx && autoreconf -i && cd -
RUN FITSINC=/opt/cfitsio-4.5.0/include \
	FITSDIR=/opt/cfitsio-4.5.0/lib \
	CXXFLAGS="-m64 -DUSE_64BIT_INTEGER -O3 -fPIC" \
	./configure --auto=sharp,cxx --enable-shared
RUN make
RUN cp lib/pkgconfig/libsharp.pc /usr/share/pkgconfig/libsharp.pc
RUN cp lib/pkgconfig/healpix_cxx.pc /usr/share/pkgconfig/healpix_cxx.pc

FROM healpix AS withgo
# Install Go
RUN wget -O /tmp/go.tar.gz https://go.dev/dl/go1.23.1.linux-amd64.tar.gz
RUN rm -rf /usr/local/go && tar -C /usr/local -xzf /tmp/go.tar.gz
ENV PATH=$PATH:/usr/local/go/bin
RUN go version


FROM withgo AS app
WORKDIR /go/src/app
# COPY . .
RUN git clone https://github.com/spenczar/healpix.git
WORKDIR /go/src/app/healpix
RUN ./run_swig.sh
WORKDIR /go/src/app
COPY ./go.mod .
COPY ./go.sum .
COPY ./main.go .
# RUN go build main.go
