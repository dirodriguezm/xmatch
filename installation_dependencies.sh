#!/bin/bash
command_exists() {
    command -v "$1" >/dev/null 2>&1
}
OS=$(uname -s)
case $OS in
    Linux)
        if [ -f /etc/lsb-release ] && grep -q "Ubuntu" /etc/lsb-release; then
            if ! command_exists apt; then
                echo "Error: apt is not installed. Please install apt-get first."
                exit 1
            fi
            echo "Installing packages for Ubuntu..."
            sudo apt update -y
            sudo apt install -y \
                autoconf \
                automake \
                libtool \
                pkg-config \
                build-essential \
                zlib1g-dev \
                libsharp-dev \
                swig \
                git \
                golang-go \
                just \
		        sqlite3
        else
            echo "Error: This script only supports Ubuntu Linux."
            exit 1
        fi
        ;;
    Darwin)
        if ! command_exists brew; then
            echo "Error: Homebrew is not installed. Please install Homebrew first."
            exit 1
        fi
        echo "Installing packages for macOS..."
        brew update
        brew install \
            autoconf \
            automake \
            libtool \
            pkg-config \
            zlib \
            swig \
            git \
            go \
            just \
            sqlite3

        ;;
    *)
        echo "Error: Unsupported operating system: $OS"
        exit 1
        ;;
esac
echo "Package i-nstallation completed successfully!"
original_script_dir="$(pwd)"
echo "Installing and building CFITSIO..."
wget -O /tmp/cfitsio.tar.gz https://heasarc.gsfc.nasa.gov/FTP/software/fitsio/c/cfitsio-4.5.0.tar.gz
tar -C /opt -xvf /tmp/cfitsio.tar.gz
cd /opt/cfitsio-4.5.0
./configure
make -j
make install
echo "Installing and building HEALPix c++ packages..."
if [ ! -d "/opt/healpix" ]; then 
    mkdir -p /opt/healpix 
fi
cd /opt/healpix
# Remove if it already exists
if [ -d "/opt/healpix/healpixmirror" ]; then 
    rm -rf /opt/healpix/healpixmirror
fi
git clone https://github.com/healpy/healpixmirror.git
cd healpixmirror
cd src/common_libraries/libsharp && autoreconf -i && cd -
cd src/cxx && autoreconf -i && cd -
FITSINC=/opt/cfitsio-4.5.0/include \
    FITSDIR=/opt/cfitsio-4.5.0/lib \
    ./configure --auto=cxx
make -j
cp lib/pkgconfig/libsharp.pc /usr/share/pkgconfig/libsharp.pc
cp lib/pkgconfig/healpix_cxx.pc /usr/share/pkgconfig/healpix_cxx.pc
echo "Installing and building HEALPix Go code..."
cd $original_script_dir
echo "Downloading xmatch service repo"
cd $original_script_dir
if [ -d "$original_script_dir/xmatch" ]; then 
    rm -rf $original_script_dir/xmatch
fi
git clone https://github.com/dirodriguezm/xmatch.git

cd xmatch

if [ -d "$original_script_dir/xmatch/healpix" ]; then 
    rm -rf $original_script_dir/xmatch/healpix
fi
git clone https://github.com/dirodriguezm/healpix.git
cd healpix
git checkout dev
./run_swig.sh
cd cmd/tester
go build main.go
export LD_LIBRARY_PATH=/opt/healpix/healpixmirror/lib:/opt/cfitsio-4.5.0/libLD_LIBRARY_PATH=/opt/healpix/healpixmirror/lib:/opt/cfitsio-4.5.0/lib
