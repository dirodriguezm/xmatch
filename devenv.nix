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
{ pkgs, config, lib, ... }:
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

  # Source tree for the release build. Only the Go module and its local
  # dependencies are needed. The SWIG wrapper is regenerated inside the
  # derivation, so it is excluded from the source to keep the build hermetic.
  releaseSrc = lib.cleanSourceWith {
    src = lib.cleanSource ./.;
    filter = path: _type:
      let
        rel = lib.removePrefix (toString ./. + "/") (toString path);
      in
      !(lib.hasPrefix ".devenv" rel)
      && !(lib.hasPrefix ".opencode" rel)
      && !(lib.hasPrefix "deploy" rel)
      && !(lib.hasPrefix "docs" rel)
      && !(lib.hasPrefix "openspec" rel)
      && !(lib.hasPrefix "service/build" rel)
      && !(lib.hasPrefix "service/data" rel)
      && !(lib.hasPrefix "service/tmp" rel)
      && !(lib.hasSuffix ".db" rel)
      && !(lib.hasSuffix ".db-shm" rel)
      && !(lib.hasSuffix ".db-wal" rel)
      && !(lib.hasSuffix "healpix_wrap.cxx" rel);
  };

  # Fully static release binary, built with the static stdenv so it can run on
  # any amd64 Linux host (e.g. Ubuntu 24.04) without Nix or any system
  # libraries. The healpix_cxx/libsharp static archives are linked in (the
  # SWIG bindings don't use the FITS code, so cfitsio is not needed);
  # `netgo`/`osusergo` avoid the static libc NSS/DNS pitfalls.
  release = pkgs.pkgsStatic.buildGoModule {
    pname = "xmatch";
    version = "0.1.0";
    src = releaseSrc;
    modRoot = "service";
    vendorHash = "sha256-j1S9Y8A/SFdz1CytpX5yNzfHiGXhoNKEZKpy/WeEdf8=";

    tags = [ "netgo" "osusergo" ];
    ldflags = [
      "-s"
      "-w"
      "-linkmode=external"
      "-extldflags=-static"
    ];

    nativeBuildInputs = [
      pkgs.swig
      pkgs.pkg-config
    ];

    env = {
      PKG_CONFIG_PATH = "${healpix}/lib/pkgconfig";
      CGO_CFLAGS = "-I${healpix}/include -I${healpix}/include/healpix_cxx";
      CGO_LDFLAGS = "-L${healpix}/lib -lhealpix_cxx -lsharp -lstdc++ -fopenmp -lm";
    };

    # `go install` names the binary after the package directory ("cmd").
    postInstall = ''
      mv $out/bin/cmd $out/bin/xmatch
    '';

    # Regenerate the SWIG bindings instead of relying on the gitignored file
    # produced by the xmatch:init-healpix task. This hook also runs for the
    # go-modules vendor derivation, which inherits postPatch.
    postPatch = ''
      pushd healpix/internal/healpix_cxx
      swig -c++ -go -intgosize 64 \
        $(pkg-config --cflags-only-I libsharp healpix_cxx) \
        -o healpix_wrap.cxx healpix_amd64.i
      popd
    '';
  };
in {
  # https://devenv.sh/reference/options/

  outputs = {
    # Static binary for deployment: `devenv build outputs.xmatch`
    xmatch = release;
  };

  packages = with pkgs; [
    healpix
    swig
    cfitsio
    golangci-lint
    sqlite
    go-migrate-sqlite
    air
    jq
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

    xwave-release = {
      exec = ''
        set -eo pipefail
        cd ${root}
        out=$(devenv build outputs.xmatch | jq -r '.["outputs.xmatch"]')
        mkdir -p ${service}/build
        install -m 0755 "$out/bin/xmatch" ${service}/build/main \
          && echo "Static release binary written to ${service}/build/main"
      '';
      description = "Build the fully static release binary (service/build/main)";
    };

    xwave-test = {
      exec = "cd ${service} && go test ./... -race";
      description = "Run all tests with race detector";
    };

    xwave-test-verbose = {
      exec = "cd ${service} && go test -v ./... -race";
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
