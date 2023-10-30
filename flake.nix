{
  description = "lox-interpreter";

  inputs = {
    nixpkgs.url = "nixpkgs";
    flake-utils.url = github:numtide/flake-utils;
  };

  outputs = { self, nixpkgs, flake-utils }:
    flake-utils.lib.eachDefaultSystem (system:
      let
        pkgs = import nixpkgs { inherit system; };
      in
      {
        devShells.default = pkgs.mkShell {
          packages = [
            pkgs.go
            pkgs.gotools
            pkgs.gopls
            pkgs.godef
            pkgs.go-outline
          ];
        };
      }
    );
}
