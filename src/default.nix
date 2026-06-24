{
    lib,
    buildGoModule,
    fetchNpmDeps,
    nodejs,
    pkg-config,
    gnumake,
    templ,
    sqlc,
    esbuild,
    tailwindcss_4,
    vips,
    glib,
}:

buildGoModule rec {
    pname = "jrdelperu-app";
    version = "1.0.0";

    src = ./.;

    vendorHash = "sha256-Lj0+sYXP0AFxXPgFOe2hg8J8kOviCmak2Im9pKlQmvw="; 

    npmDeps = fetchNpmDeps {
        src = ./.;
        hash = "sha256-Fk/MlvHJ8ZXzkSMzBG6o/TP4wgkxi/wgDFAhEWbMZzY=";
    };

    nativeBuildInputs = [
        pkg-config
        gnumake
        templ
        sqlc
        esbuild
        tailwindcss_4
        nodejs
    ];

    buildInputs = [
        vips
        glib
    ];

    env.CGO_ENABLED = "1";

    preBuild = ''
        export HOME=$(mktemp -d)
        export npm_config_cache=$(mktemp -d)
        cp -R ${npmDeps}/* $npm_config_cache/
        chmod -R +w $npm_config_cache
        npm ci --prefer-offline --no-audit
    '';

    buildPhase = ''
        runHook preBuild

        make clean
        make ./build/server
        make ./build/seeder

        runHook postBuild
    '';

    installPhase = ''
        runHook preInstall

        mkdir -p $out/bin
        cp ./build/server $out/bin/server
        cp ./build/seeder $out/bin/seeder

        mkdir -p $out/src/db
        cp -r ./db $out/src/

        runHook postInstall
    '';
}
