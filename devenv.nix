# Copyright 2024-2025 Matías Medina Silva
# Copyright 2026 Diego Rodriguez Mancini
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
{ pkgs, config, ... }:
let
  root = config.devenv.root;
  service = "${root}/service";

  go-migrate-sqlite = pkgs.go-migrate.overrideAttrs (oldAttrs: {
    tags = ["sqlite3"];
  });

  healpix = pkgs.stdenv.mkDerivation {
    pname = "healpix";
    version = "1.16.5";
    src = pkgs.fetchFromGitHub {
      owner = "healpy";
      repo = "healpixmirror";
      rev = "a44dd367e50a038ea38f1916ff78c20607fde315";
      hash = "sha256-R/lKqhr0m55lpWJy7TAflRyWgv20b1xWj8KH1TboL9w=";
    };

    nativeBuildInputs = with pkgs; [
      autoconf
      automake
      libtool
      pkg-config
      gfortran
      patchelf
    ];

    buildInputs = with pkgs; [
      cfitsio
    ];

    configureFlags = [
      "--auto=cxx"
      "--prefix=${placeholder "out"}"
    ];

    FITSINC = "${pkgs.cfitsio}/include";
    FITSDIR = "${pkgs.cfitsio}/lib";

    preConfigure = ''
      cd src/common_libraries/libsharp && autoreconf -i && cd -
      cd src/cxx && autoreconf -i && cd -
    '';

    buildPhase = ''
      make -j
    '';

    installPhase = ''
      mkdir -p $out/lib $out/include
      cp -a lib/* $out/lib/
      cp -r include/healpix_cxx/* $out/include
    '' + (if pkgs.stdenv.hostPlatform.isDarwin then ''
        install_name_tool -id $out/lib/libhealpix_cxx.4.dylib $out/lib/libhealpix_cxx.4.dylib
        install_name_tool -change /private/tmp/nix-build-healpix-1.16.5.drv-0/source/lib/libhealpix_cxx.4.dylib $out/lib/libhealpix_cxx.4.dylib $out/lib/libhealpix_cxx.4.dylib || true
        install_name_tool -change /private/tmp/nix-build-healpix-1.16.5.drv-0/source/lib/libsharp.2.dylib $out/lib/libsharp.2.dylib $out/lib/libhealpix_cxx.4.dylib || true

        install_name_tool -id $out/lib/libsharp.2.dylib $out/lib/libsharp.2.dylib
        install_name_tool -change /private/tmp/nix-build-healpix-1.16.5.drv-0/source/lib/libsharp.2.dylib $out/lib/libsharp.2.dylib $out/lib/libsharp.2.dylib || true

        substituteInPlace $out/lib/libhealpix_cxx.la \
            --replace "/private/tmp/nix-build-healpix-1.16.5.drv-0/source" "$out"

        substituteInPlace $out/lib/libsharp.la \
            --replace "/private/tmp/nix-build-healpix-1.16.5.drv-0/source" "$out"

        substituteInPlace $out/lib/pkgconfig/healpix_cxx.pc \
            --replace "/private/tmp/nix-build-healpix-1.16.5.drv-0/source" "$out"

        substituteInPlace $out/lib/pkgconfig/libsharp.pc \
            --replace "/private/tmp/nix-build-healpix-1.16.5.drv-0/source" "$out"
      ''
      else ''
        patchelf --set-rpath "${pkgs.cfitsio}/lib:$out/lib" $out/lib/libhealpix_cxx.so.4.0.5
        patchelf --set-rpath "${pkgs.cfitsio}/lib:$out/lib" $out/lib/libsharp.so.2.0.2

        sed -i "s|includedir=.|includedir=$out/include|g" $out/lib/pkgconfig/healpix_cxx.pc
        sed -i "s|/build/source/|$out/|g" $out/lib/pkgconfig/healpix_cxx.pc
      '');

    dontAddPrefix = true;
    meta = with pkgs.lib; {
      description = "Library for fast spherical harmonic transforms";
      homepage = "https://github.com/Libsharp/libsharp";
      license = licenses.gpl2;
      platforms = platforms.unix;
    };
  };
in {
  # https://devenv.sh/reference/options/

  packages = with pkgs; [
    healpix
    swig
    cfitsio
    golangci-lint
    sqlite
    go-migrate-sqlite
    grc
    air
    tailwindcss_4
    tailwindcss-language-server
  ];

  env = {
    LD_LIBRARY_PATH = pkgs.lib.makeLibraryPath [
      healpix
      pkgs.cfitsio
      pkgs.gcc.cc.lib
    ];
    PKG_CONFIG_PATH = "${healpix}/lib/pkgconfig";
    CGO_CFLAGS = "-I${healpix}/include -I${healpix}/include/healpix_cxx -I${pkgs.cfitsio}/include ";
    CGO_LDFLAGS = "-L${healpix}/lib -L${pkgs.cfitsio}/lib  -lhealpix_cxx -lcfitsio ";
    GRC_CONFIG = ''
      # Regla para "go test" y "make test"
      \b(go test)\b
      regexp==== RUN .*
      colour=bright_blue
      -
      regexp=--- PASS: .* (\(\d+\.\d+s\))
      colour=green, yellow
      -
      regexp=^PASS$
      colour=bold white on_green
      -
      regexp=^(ok|FAIL)\s+.*
      colour=default, magenta
      -
      regexp=--- FAIL: .* (\(\d+\.\d+s\))
      colour=red, yellow
      -
      regexp=^FAIL$
      colour=bold white on_red
      -
      regexp=[^\s]+\.go(:\d+)?
      colour=cyan
    '';
  };

  languages.go = {
    enable = true;
  };

  scripts = {
    xwave-init-healpix = {
      exec = "devenv tasks run xmatch:init-healpix";
      description = "Initialize the healpix submodule and generate SWIG bindings";
    };

    xwave-build = {
      exec = "cd ${service} && go build -o build/main cmd/*.go";
      description = "Build the Go binary (service/build/main)";
    };

    xwave-test = {
      exec = "cd ${service} && grc go test ./... -race";
      description = "Run all tests with race detector";
    };

    xwave-test-verbose = {
      exec = "cd ${service} && grc go test -v ./... -race";
      description = "Run all tests verbosely with race detector";
    };

    xwave-run = {
      exec = ''
        cd ${service}
        go build -o build/main cmd/*.go
        LOG_LEVEL="''${LOG_LEVEL:-debug}" ./build/main "$1" "$2"
      '';
      description = "Build and run an application (run <application> [flags])";
    };

    xwave-migrate = {
      exec = "migrate -database sqlite3://${root}/\"$1\".db -path ${service}/internal/db/migrations up";
      description = "Run database migrations (migrate <db>)";
    };

    xwave-live-server = {
      exec = "cd ${service} && USE_LOGGER=true ENVIRONMENT=local LOG_LEVEL=debug air";
      description = "Run with air for live reload";
    };

    xwave-docs = {
      exec = "cd ${service} && go run github.com/swaggo/swag/cmd/swag@v1.16.4 init --dir ./ --generalInfo ./cmd/*.go --output ./docs";
      description = "Generate Swagger documentation";
    };

    xwave-mock = {
      exec = "cd ${service} && go run github.com/vektra/mockery/v3@v3.7.0";
      description = "Generate mocks with mockery";
    };

    xwave-clean-build = {
      exec = "rm -r ${service}/build";
      description = "Remove service/build/";
    };

    xwave-clean-all = {
      exec = "cd ${service} && go clean && go clean -testcache && rm -r build";
      description = "Clean Go caches and build artifacts";
    };

    xwave-clean-db = {
      exec = "rm ${root}/\"$1\".db";
      description = "Remove a database file (clean-db <db>)";
    };
  };

  tasks."xmatch:init-healpix" = {
    status = "test -f ${root}/healpix/internal/healpix_cxx/healpix_wrap.cxx";
    exec = ''
      echo "Initializing submodule and running script..."
      git -C ${root} submodule update --init --recursive
      cd ${root}/healpix && bash run_swig.sh
    '';
    before = ["devenv:enterShell"];
  };
}
