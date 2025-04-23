{
  description = "Flake for github:mahyarmirrashed/llm-video-analyzer";

  inputs = {
    nixpkgs.url = "github:NixOS/nixpkgs/nixos-24.11";
    flake-utils.url = "github:numtide/flake-utils";
  };

  outputs = {
    self,
    nixpkgs,
    flake-utils,
  }:
    flake-utils.lib.eachDefaultSystem (
      system: let
        pkgs = nixpkgs.legacyPackages.${system};

        mahyarmirrashed-llm-video-analyzer = pkgs.stdenv.mkDerivation {
          pname = "mahyarmirrashed-llm-video-analyzer";
          version = "0.1.0";
          src = self;
        };

        runtimeEnv = with pkgs; [ffmpeg];
      in {
        formatter = pkgs.alejandra;

        devShells.default = pkgs.mkShell {
          packages = with pkgs; [git go gotools lazygit] ++ runtimeEnv;
        };

        packages.default = mahyarmirrashed-llm-video-analyzer;
      }
    );
}
