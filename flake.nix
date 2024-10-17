{
  inputs = {
    nixpkgs.url = "nixpkgs";
    nixpkgs-unstable.url = "github:nixos/nixpkgs/master";
    flake-utils.url = "github:numtide/flake-utils";
  };

  outputs = { self, nixpkgs, nixpkgs-unstable, flake-utils }:
    let
      localOverlay = import ./nix/overlay.nix;
      overlayUnstable = final: prev: {
        # allow the ability to add an unstable package as required.
        unstable = nixpkgs-unstable.legacyPackages.${prev.system};
      };
      overlays = [ overlayUnstable localOverlay ];
    in flake-utils.lib.eachDefaultSystem (system:
      let
        pkgs = import nixpkgs {
          inherit system overlays;
        };
      in {
        legacyPackages = pkgs;
        inherit (pkgs) devShell;
      }) // {
        # platform independent attrs
        overlay = final: prev: (nixpkgs.lib.composeManyExtensions overlays) final prev;
        inherit overlays;
      };
}
