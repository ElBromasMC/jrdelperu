let
    # nixos-26.05 (picked at 2026-06-23)
    nixpkgs = fetchTarball "https://github.com/NixOS/nixpkgs/archive/34268251cf5547d39063f2c5ea9a196246f7f3a6.tar.gz";
    pkgs = import nixpkgs { config = {}; overlays = []; };
    buildBaseImage = import ./base.nix;
    name = "jrdelperu_prod";

    app = pkgs.callPackage ./src/default.nix {};
in

buildBaseImage {
    inherit pkgs name;
    tag = "latest";
    includeStorePaths = true; 
    extraContents = [
        app
        pkgs.vips
        pkgs.tzdata
        (pkgs.go-migrate.overrideAttrs (oldAttrs: {
            tags = [ "postgres" ];
        }))
    ];
    configCmd = [ "${app}/bin/server" ];
    enableDebugTools = true;
    extraFakeRootCommands = ''
        mkdir -p /srv/app/uploads
        chown runner:runner /srv/app/uploads
    '';
}
