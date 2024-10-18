FROM golang:1.23
RUN apt update -y
RUN apt install -y \
	autoconf \
	automake \
	libtool \
	pkg-config \
	build-essential \
	zlib1g-dev \
	libsharp-dev \
	swig \
	git

RUN wget -O /tmp/cfitsio.tar.gz https://heasarc.gsfc.nasa.gov/FTP/software/fitsio/c/cfitsio-4.5.0.tar.gz
RUN tar -C /opt -xvf /tmp/cfitsio.tar.gz
WORKDIR /opt/cfitsio-4.5.0
RUN ./configure
RUN make -j
RUN make install

WORKDIR /opt/healpix
RUN git clone https://github.com/healpy/healpixmirror.git
WORKDIR /opt/healpix/healpixmirror
RUN cd src/common_libraries/libsharp && autoreconf -i && cd -
RUN cd src/cxx && autoreconf -i && cd -
RUN FITSINC=/opt/cfitsio-4.5.0/include \
	FITSDIR=/opt/cfitsio-4.5.0/lib \
	./configure --auto=cxx
RUN make -j
RUN cp lib/pkgconfig/libsharp.pc /usr/share/pkgconfig/libsharp.pc
RUN cp lib/pkgconfig/healpix_cxx.pc /usr/share/pkgconfig/healpix_cxx.pc

WORKDIR /go/src/app
COPY ./healpix /go/src/app/healpix
RUN cd healpix && ./run_swig.sh 
COPY ./go.mod .
COPY ./go.sum .
COPY ./main.go .
RUN go build main.go
ENV LD_LIBRARY_PATH=/opt/healpix/healpixmirror/lib:/opt/cfitsio-4.5.0/lib
CMD ["./main", "ang2pix", "20", "1"]
