{
  description = "Crazy Note Taking TUI";

  inputs = {
    nixpkgs.url = "github:nixos/nixpkgs?ref=25.05";
    systems.url = "github:nix-systems/default";
    flake-utils = {
      url = "github:numtide/flake-utils";
      inputs.systems.follows = "systems";
    };
  };

  outputs = { self, systems, nixpkgs, flake-utils }:
    flake-utils.lib.eachDefaultSystem (system:
      let
        pkgs = import nixpkgs { inherit system; };
      in
      {
        devShells.default = pkgs.mkShell {
          buildInputs = with pkgs; [
            go
          ];
        };
        packages.default = pkgs.callPackage ./build.nix { };
      });
}
