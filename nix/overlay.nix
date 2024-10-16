final: prev: {
  devShell = final.callPackage ./sentinel_sdk.nix { };

  go = final.unstable.go;
}
