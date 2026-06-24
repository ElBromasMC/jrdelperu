{
    pkgs,
    name,
    tag ? "latest",
    includeStorePaths ? true,
    extraContents ? [],
    extraFakeRootCommands ? "",
    configCmd,
    extraConfigEnv ? [],
    enableDebugTools ? false,
}:
pkgs.dockerTools.streamLayeredImage {
    inherit name tag includeStorePaths;
    contents = [
        pkgs.dockerTools.caCertificates
        (pkgs.dockerTools.fakeNss.override {
            extraPasswdLines = [ "runner::1000:1000:Runner User:/home/runner:/bin/bash" ];
            extraGroupLines = [ "runner:x:1000:" ];
        })
    ]
    ++ pkgs.lib.optionals enableDebugTools [
        pkgs.dockerTools.usrBinEnv
        pkgs.dockerTools.binSh
        pkgs.bash
        pkgs.getent
        pkgs.curlMinimal
        pkgs.busybox
        pkgs.perl
        pkgs.inotify-tools
        pkgs.git
        pkgs.ripgrep
        pkgs.fd
    ]
    ++ extraContents;
    enableFakechroot = true;
    fakeRootCommands = ''
        # Tmp folder setup
        mkdir /tmp /var/tmp /run
        chmod 1777 /tmp
        chmod 1777 /var/tmp
        ln -s ../run /var/run

        # Runner user setup
        for dir in \
            "/tmp/runtime-dir" \
            "/home/runner" \
            "/srv/app"
        do
            mkdir -p "$dir"
            chown runner:runner "$dir"
        done
        chmod 0700 /tmp/runtime-dir

        ${extraFakeRootCommands}
    '';
    config = {
        Cmd = configCmd;
        User = "1000:1000";
        WorkingDir = "/srv/app";
        Env = [
            "XDG_RUNTIME_DIR=/tmp/runtime-dir"
            "HOME=/home/runner"
            "USER=runner"
            "LANG=C.UTF-8"
        ]
        ++ extraConfigEnv;
    };
}
