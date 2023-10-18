{
    inputs = {
        nixpkgs.url = "nixpkgs";
    };

    outputs = { self, nixpkgs, flake-utils }:
        flake-utils.lib.eachDefaultSystem(system:
            let
                pkgs = nixpkgs.legacyPackages.${system};
            in
            {
                devShell = pkgs.mkShell {
                    packages = [
                        pkgs.go_1_21
                    ];
                };
            });
}
