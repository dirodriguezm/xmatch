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
  inputs = {
    nixpkgs.url = "github:cachix/devenv-nixpkgs/rolling";
    systems.url = "github:nix-systems/default";
    devenv.url = "github:cachix/devenv";
    devenv.inputs.nixpkgs.follows = "nixpkgs";
    healpix-src = {
      flake = false;
      owner = "healpy";
      repo = "healpixmirror";
      rev = "a44dd367e50a038ea38f1916ff78c20607fde315";
      type = "github";
    };
  };

  nixConfig = {
    extra-trusted-public-keys = "devenv.cachix.org-1:w1cLUi8dv3hnoSPGAuibQv+f9TZLr6cv/Hm9XgU50cw=";
    extra-substituters = "https://devenv.cachix.org";
  };

  outputs = {
    self,
    nixpkgs,
    devenv,
    systems,
    healpix-src,
    ...
  } @ inputs: let
    forEachSystem = nixpkgs.lib.genAttrs (import systems);
  in {
    packages = forEachSystem (system: {
      devenv-up = self.devShells.${system}.default.config.procfileScript;
      devenv-test = self.devShells.${system}.default.config.test;
    });

    devShells =
      forEachSystem
      (system: let
        pkgs = nixpkgs.legacyPackages.${system};
        healpix = pkgs.stdenv.mkDerivation {
          pname = "healpix";
          version = "1.16.5";
          src = healpix-src;

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

          enableParallelBuilding = true;

          buildPhase = ''
            make -j
          '';

          installPhase = ''
            mkdir -p $out/lib $out/include
            ls -l lib

            cp -a lib/* $out/lib/
            ls -l $out/lib
            ls -l $out/lib/pkgconfig
            cp -r include/healpix_cxx/* $out/include

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
          '';

          dontAddPrefix = true;
          meta = with pkgs.lib; {
            description = "Library for fast spherical harmonic transforms";
            homepage = "https://github.com/Libsharp/libsharp";
            license = licenses.gpl2;
            platforms = platforms.unix;
          };
        };
      in {
        default = devenv.lib.mkShell {
          inherit inputs pkgs;
          modules = [
            {
              # https://devenv.sh/reference/options/
              packages = with pkgs; [
                healpix
                just
                swig
                cfitsio
              ];

              env = {
                LD_LIBRARY_PATH = pkgs.lib.makeLibraryPath [
                  healpix
                  pkgs.cfitsio
                ];
                PKG_CONFIG_PATH = "${healpix}/lib/pkgconfig";
                CGO_CFLAGS = "-I${healpix}/include -I${healpix}/include/healpix_cxx -I${pkgs.cfitsio}/include ";
                CGO_LDFLAGS = "-L${healpix}/lib -L${pkgs.cfitsio}/lib  -lhealpix_cxx -lcfitsio ";
              };

              scripts.init-healpix.exec = ''
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

              enterShell = ''
                echo "Healpix development environment configured."
                init-healpix
              '';

              languages.go = {
                enable = true;
                package = pkgs.go_1_23;
              };
            }
          ];
        };
      });
  };
}
