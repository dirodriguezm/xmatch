{
  description = "Entorno de desarrollo para Healpix con libsharp y cfitsio";

  inputs = {
    nixpkgs.url = "github:NixOS/nixpkgs/nixos-unstable"; # Puedes usar una versión específica
    flake-utils.url = "github:numtide/flake-utils";

    libsharp-src = {
      url = "github:Libsharp/libsharp/master"; # Considera usar un commit específico
      flake = false;
      hash = "sha256-11czd414gmdwp7nqk2b1gy52ddyhvkfx1w88jwj0jqdbnz1c1qdk";
    };

    healpix-src = {
      url = "github:healpy/healpixmirror/trunk"; # Considera usar un commit específico
      flake = false;
      hash = "sha256-1p1gx0vdb1y2ixb5qvxlzn19c74m3wqfswk2lmjrx6zl3am4mya7";
    };

    cfitsio-src = {
      url = "https://heasarc.gsfc.nasa.gov/FTP/software/fitsio/c/cfitsio-4.5.0.tar.gz";
      flake = false;
      hash = "sha256-5IVPwzZcFGLkk6pYa/qi89C7jCC3WlJJVdtkwnQnzgk=";
    };
  };

  outputs = {
    # self,
    nixpkgs,
    flake-utils,
    libsharp-src,
    healpix-src,
    cfitsio-src,
    ...
  }:
    flake-utils.lib.eachDefaultSystem (
      system: let
        pkgs = nixpkgs.legacyPackages.${system};

        # Derivación para libsharp
        libsharp = pkgs.stdenv.mkDerivation {
          pname = "libsharp";
          version = "unstable-2024-05-31";
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

        # Derivación para cfitsio
        cfitsio = pkgs.stdenv.mkDerivation {
          pname = "cfitsio";
          version = "4.5.0";
          src = cfitsio-src;

          buildInputs = [pkgs.zlib];
          buildPhase = ''
            make -j
          '';
          configureFlags = ["--enable-reentrant"];
        };

        # Derivación para healpix
        healpix = pkgs.stdenv.mkDerivation {
          pname = "healpix";
          version = "2024-05-31";
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
            make cpp-all -j
          '';

          installPhase = ''
            mkdir -p $out/lib
            cp -r lib/pkgconfig $out/lib/
            cp lib/libhealpix_cxx.a $out/lib/
            cp lib/libhealpix_cxx.la $out/lib/
            cp lib/libhealpix_cxx.so.4.0.5 $out/lib/
            cp lib/libsharp.a $out/lib/
            cp lib/libsharp.la $out/lib/
            cp lib/libsharp.so.2.0.2 $out/lib/
            cp -r include $out/

            ln -s $out/lib/libhealpix_cxx.so.4.0.5 $out/lib/libhealpix_cxx.so.4
            ln -s $out/lib/libhealpix_cxx.so.4 $out/lib/libhealpix_cxx.so
            ln -s $out/lib/libsharp.so.2.0.2 $out/lib/libsharp.so.2
            ln -s $out/lib/libsharp.so.2 $out/lib/libsharp.so

            patchelf --set-rpath "${cfitsio}/lib:${libsharp}/lib:$out/lib" $out/lib/libhealpix_cxx.so.4.0.5
            patchelf --set-rpath "${cfitsio}/lib:$out/lib" $out/lib/libsharp.so.2.0.2

            sed -i "s|includedir=.*|includedir=$out/include|g" $out/lib/pkgconfig/healpix_cxx.pc
            sed -i "s|/build/source/|$out/|g" $out/lib/pkgconfig/healpix_cxx.pc
          '';

          dontAutoPatchelf = true;

          meta = with pkgs.lib; {
            description = "Librería para transformadas armónicas esféricas rápidas";
            homepage = "https://github.com/Libsharp/libsharp";
            license = licenses.gpl2;
            platforms = platforms.unix;
          };
        };
      in {
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
          ];

          shellHook = ''
            # Configuración para pkg-config
            export PKG_CONFIG_PATH=${healpix}/lib/pkgconfig:${cfitsio}/lib/pkgconfig:${libsharp}/lib/pkgconfig:$PKG_CONFIG_PATH
            # Configuración para CGO
            export CGO_CFLAGS="-I${healpix}/include -I${healpix}/include/healpix_cxx -I${cfitsio}/include -I${libsharp}/include"
            export CGO_CPPFLAGS="$CGO_CFLAGS"
            export CGO_CXXFLAGS="$CGO_CFLAGS"
            export CGO_LDFLAGS="-L${healpix}/lib -L${cfitsio}/lib -L${libsharp}/lib -lhealpix_cxx -lcfitsio -lsharp"
            # Asegurar que el tiempo de ejecución encuentre las bibliotecas compartidas
            export LD_LIBRARY_PATH=${healpix}/lib:${cfitsio}/lib:${libsharp}/lib:${pkgs.gcc.cc.lib}/lib:$LD_LIBRARY_PATH
            echo "Entorno de desarrollo Healpix configurado."
            echo "Variables configuradas:"
            echo "  PKG_CONFIG_PATH=$PKG_CONFIG_PATH"
            echo "  CGO_CFLAGS=$CGO_CFLAGS"
            echo "  CGO_LDFLAGS=$CGO_LDFLAGS"
            echo "  LD_LIBRARY_PATH=$LD_LIBRARY_PATH"
          '';
        };

        packages = {
          inherit libsharp cfitsio healpix;
          default = healpix; # El paquete por defecto
        };
      }
    );
}
