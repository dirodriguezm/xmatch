# Copyright 2024-2025 Matías Medina Silva
#
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
# http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.
#
{
  description = "Development environment for Healpix with libsharp and cfitsio";

  # Required inputs (dependencies)
  inputs = {
    nixpkgs.url = "github:NixOS/nixpkgs/nixos-unstable";
    flake-utils.url = "github:numtide/flake-utils";

    # Source repositories for specific projects
    libsharp-src = {
      flake = false;
      owner = "Libsharp";
      repo = "libsharp";
      rev = "8d519468b7932858ff48ef0bdcdeb36561090994";
      type = "github";
    };

    healpix-src = {
      flake = false;
      owner = "healpy";
      repo = "healpixmirror";
      rev = "a44dd367e50a038ea38f1916ff78c20607fde315";
      type = "github";
    };
  };

  outputs = {
    nixpkgs,
    flake-utils,
    libsharp-src,
    healpix-src,
    ...
  }:
    flake-utils.lib.eachDefaultSystem (
      system: let
        pkgs = nixpkgs.legacyPackages.${system};

        # Derivation for libsharp
        libsharp = pkgs.stdenv.mkDerivation {
          pname = "libsharp";
          version = "1.0.0";
          src = libsharp-src;

          nativeBuildInputs = [
            pkgs.autoconf
            pkgs.automake
            pkgs.libtool
          ];

          buildInputs = [
            pkgs.fftw
            pkgs.fftwFloat
          ];

          preConfigure = ''
            autoconf
          '';

          configureFlags = [
            "--enable-openmp"
            "--enable-pic"
            "--prefix=${placeholder "out"}"
          ];

          enableParallelBuilding = true;

          buildPhase = ''
            make SHARP_TARGET=auto compile_all
          '';

          installPhase = ''
            mkdir -p $out/lib $out/include $out/bin
            cp auto/lib/*.a $out/lib/
            cp -r auto/include/* $out/include/
            [ -f auto/bin/sharp_testsuite ] && cp auto/bin/sharp_testsuite $out/bin/ || true
            if [ -f python/libsharp/libsharp.so ]; then
              mkdir -p $out/${pkgs.python.sitePackages}
              cp python/libsharp/*.so $out/${pkgs.python.sitePackages}/
            fi
          '';

          meta = with pkgs.lib; {
            description = "Library for fast spherical harmonic transforms";
            homepage = "https://github.com/Libsharp/libsharp";
            license = licenses.gpl2Plus;
            platforms = platforms.unix;
          };
        };

        # Derivation for cfitsio
        cfitsio = pkgs.stdenv.mkDerivation rec {
          pname = "cfitsio";
          version = "4.5.0";
          src = pkgs.fetchurl {
            url = "https://heasarc.gsfc.nasa.gov/FTP/software/fitsio/c/cfitsio-${version}.tar.gz";
            hash = "sha256-5IVPwzZcFGLkk6pYa/qi89C7jCC3WlJJVdtkwnQnzgk=";
          };

          buildInputs = [pkgs.zlib];
          buildPhase = ''
            make -j
          '';
          configureFlags = ["--enable-reentrant"];
        };

        # Derivation for healpix
        healpix = pkgs.stdenv.mkDerivation {
          pname = "healpix";
          version = "1.16.5";
          src = healpix-src;

          nativeBuildInputs = [
            pkgs.autoconf
            pkgs.automake
            pkgs.libtool
            pkgs.pkg-config
            pkgs.patchelf
          ];

          buildInputs = [
            cfitsio
            pkgs.gfortran
            libsharp
            pkgs.zlib
            pkgs.gcc.cc.lib
          ];

          dontAddPrefix = true;

          configureFlags = [
            "--auto=cxx"
            "--prefix=${placeholder "out"}"
          ];

          FITSINC = "${cfitsio}/include";
          FITSDIR = "${cfitsio}/lib";

          preConfigure = ''
            cd src/common_libraries/libsharp && autoreconf -i && cd -
            cd src/cxx && autoreconf -i && cd -
          '';

          buildPhase = ''
            make -j
          '';

          installPhase = ''
            mkdir -p $out/lib $out/include
            cp -r lib/pkgconfig $out/lib/
            cp lib/libhealpix_cxx.a $out/lib/
            cp lib/libhealpix_cxx.la $out/lib/
            cp lib/libhealpix_cxx.so.4.0.5 $out/lib/
            cp lib/libsharp.a $out/lib/
            cp lib/libsharp.la $out/lib/
            cp lib/libsharp.so.2.0.2 $out/lib/
            cp -r include/healpix_cxx/* $out/include

            ln -s $out/lib/libhealpix_cxx.so.4.0.5 $out/lib/libhealpix_cxx.so.4
            ln -s $out/lib/libhealpix_cxx.so.4 $out/lib/libhealpix_cxx.so
            ln -s $out/lib/libsharp.so.2.0.2 $out/lib/libsharp.so.2
            ln -s $out/lib/libsharp.so.2 $out/lib/libsharp.so

            patchelf --set-rpath "${cfitsio}/lib:${libsharp}/lib:$out/lib" $out/lib/libhealpix_cxx.so.4.0.5
            patchelf --set-rpath "${cfitsio}/lib:$out/lib" $out/lib/libsharp.so.2.0.2

            sed -i "s|includedir=.|includedir=$out/include|g" $out/lib/pkgconfig/healpix_cxx.pc
            sed -i "s|/build/source/|$out/|g" $out/lib/pkgconfig/healpix_cxx.pc
          '';

          dontAutoPatchelf = true;

          meta = with pkgs.lib; {
            description = "Library for fast spherical harmonic transforms";
            homepage = "https://github.com/Libsharp/libsharp";
            license = licenses.gpl2;
            platforms = platforms.unix;
          };
        };
        initHealpix = pkgs.writeShellScriptBin "init-healpix" ''
          echo "Initializing submodule and running script..."
          git submodule update --init --recursive
          cd healpix
          cd internal/healpix_cxx
          # Check if the output file already exists
          if [ -f "healpix_wrap.cxx" ]; then
            echo "healpix_wrap.cxx already exists. The script may have already been run."
            exit 0
          fi
          cd ../..
          bash run_swig.sh
          cd -
        '';
      in {
        # Development environment definition
        devShells.default = pkgs.mkShell {
          buildInputs = with pkgs; [
            libsharp
            cfitsio
            healpix
            go_1_23
            just
            pkg-config
            swig
            gcc.cc.lib
            initHealpix
          ];

          LD_LIBRARY_PATH = pkgs.lib.makeLibraryPath [
            healpix
            cfitsio
            libsharp
            pkgs.gcc.cc.lib
          ];

          PKG_CONFIG_PATH = "${healpix}/lib/pkgconfig:${cfitsio}/lib/pkgconfig:${libsharp}/lib/pkgconfig";
          CGO_CFLAGS = "-I${healpix}/include -I${healpix}/include/healpix_cxx -I${cfitsio}/include -I${libsharp}/include";
          CGO_LDFLAGS = "-L${healpix}/lib -L${cfitsio}/lib -L${libsharp}/lib -lhealpix_cxx -lcfitsio -lsharp";

          shellHook = ''
            # Configuration for pkg-config
            echo "Healpix development environment configured."
            init-healpix
          '';
        };

        packages = {
          inherit libsharp cfitsio healpix;
          default = healpix; # The default package
        };
      }
    );
}
