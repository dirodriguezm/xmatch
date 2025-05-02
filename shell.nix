{pkgs ? import <nixpkgs> {}}: let
  libsharp = pkgs.stdenv.mkDerivation {
    pname = "libsharp";
    version = "unstable-2024-05-31";

    src = pkgs.fetchFromGitHub {
      owner = "Libsharp";
      repo = "libsharp";
      rev = "master"; # Consider pinning to a specific commit
      sha256 = "11czd414gmdwp7nqk2b1gy52ddyhvkfx1w88jwj0jqdbnz1c1qdk"; # Update if needed
    };

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
      # "--enable-mpi" # Optional
      "--enable-pic" # Needed for Python extensions
      "--prefix=${placeholder "out"}" # Tell configure where to install
    ];

    enableParallelBuilding = true;

    # The Makefile doesn't have 'install' but has 'hdrcopy' and builds in auto/
    buildPhase = ''
      make SHARP_TARGET=auto compile_all
    '';

    installPhase = ''
      # Create necessary directories
      mkdir -p $out/lib $out/include $out/bin

      # Install libraries
      cp auto/lib/*.a $out/lib/

      # Install headers
      cp -r auto/include/* $out/include/

      # Install binaries if they exist
      [ -f auto/bin/sharp_testsuite ] && cp auto/bin/sharp_testsuite $out/bin/

      # Install Python modules if built
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

  healpix = pkgs.stdenv.mkDerivation rec {
    pname = "healpix";
    version = "2024-05-31";
    src = pkgs.fetchFromGitHub {
      owner = "healpy";
      repo = "healpixmirror";
      rev = "trunk";
      sha256 = "1p1gx0vdb1y2ixb5qvxlzn19c74m3wqfswk2lmjrx6zl3am4mya7"; # Calcula con nix-prefetch-git
    };
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
      "--prefix=${placeholder "out"}" # Tell configure where to install
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
        # Create necessary directories
      mkdir -p $out/lib

      # Install libraries
      ls lib -l
      cp -r lib/pkgconfig $out/lib/
      cp lib/libhealpix_cxx.a $out/lib/
      cp lib/libhealpix_cxx.la $out/lib/
      cp lib/libhealpix_cxx.so.4.0.5 $out/lib/
      cp lib/libsharp.a $out/lib/
      cp lib/libsharp.la $out/lib/
      cp lib/libsharp.so.2.0.2 $out/lib/

      cp -r include $out/

      # Crear enlaces simbólicos para las versiones sin número de versión
      ln -s $out/lib/libhealpix_cxx.so.4.0.5 $out/lib/libhealpix_cxx.so.4
      ln -s $out/lib/libhealpix_cxx.so.4 $out/lib/libhealpix_cxx.so
      ln -s $out/lib/libsharp.so.2.0.2 $out/lib/libsharp.so.2
      ln -s $out/lib/libsharp.so.2 $out/lib/libsharp.so

      # Corregir RPATH de las bibliotecas
      patchelf --set-rpath "${cfitsio}/lib:${libsharp}/lib:$out/lib" $out/lib/libhealpix_cxx.so.4.0.5
      patchelf --set-rpath "${cfitsio}/lib:$out/lib" $out/lib/libsharp.so.2.0.2

      # Corregir el archivo pkg-config para que incluya las rutas correctas
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
in
  pkgs.mkShell {
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

    LD_LIBRARY_PATH = "${pkgs.lib.makeLibraryPath [cfitsio healpix]}";
    PKG_CONFIG_PATH = "${healpix}/lib/pkgconfig";

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
  }
