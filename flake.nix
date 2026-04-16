{
  inputs = {
    nixpkgs.url = "github:NixOS/nixpkgs/nixpkgs-unstable";
  };

  outputs =
    { nixpkgs, ... }:
    let
      forAllSystems = nixpkgs.lib.genAttrs nixpkgs.lib.systems.flakeExposed;
      pkgsFor = system: import nixpkgs { inherit system; };
    in
    {
      packages = forAllSystems (
        system:
        let
          pkgs = pkgsFor system;
        in
        rec {
          flack = pkgs.buildGoModule {
            pname = "flack";
            version = "0.1.0";
            src = ./.;
            vendorHash = "sha256-efDJq5gmBYDgE8TqW+o4a6I0H/zCsluUuO0v4zxAA2E=";
          };
          default = flack;
        }
      );

      devShells = forAllSystems (
        system:
        let
          pkgs = pkgsFor system;
        in
        {
          default = pkgs.mkShell {
            packages = with pkgs; [
              go
              gopls
              gotools
            ];
          };
        }
      );
    };
}
