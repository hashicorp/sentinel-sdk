{ go
, delve
, nodejs
, zlib
, mozjpeg
, libtool
, golangci-lint
, kgt
, mkShell
}:

mkShell rec {
  name = "sentinel-sdk";

  hardeningDisable = [ "fortify" ];

  packages = [
    go

    # dev tools
    delve
    golangci-lint
  ];
}
