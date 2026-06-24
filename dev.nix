let
    # nixos-26.05 (picked at 2026-06-23)
    nixpkgs = fetchTarball "https://github.com/NixOS/nixpkgs/archive/34268251cf5547d39063f2c5ea9a196246f7f3a6.tar.gz";
    # nixos-unstable-small (picked 2026-06-23)
    nixpkgsUnstable = fetchTarball "https://github.com/NixOS/nixpkgs/archive/6f5e2d2ea55618c5197ce40e346f205c35d2213e.tar.gz";
    pkgs = import nixpkgs { config = {}; overlays = []; };
    pkgsUnstable = import nixpkgsUnstable { config = { allowUnfree = true; }; overlays = []; };
    buildBaseImage = import ./base.nix;

    name = "jrdelperu_dev";
in

buildBaseImage {
    inherit pkgs name;
    tag = "latest";
    includeStorePaths = false;
    extraContents = [
        pkgs.go
        pkgs.stdenv.cc
        pkgs.gnumake
        pkgs.pkg-config
        pkgs.templ
        pkgs.sqlc
        (pkgs.go-migrate.overrideAttrs (oldAttrs: {
            tags = [ "postgres" ];
        }))
        pkgs.nodejs_26
        pkgs.tailwindcss_4
        pkgs.esbuild
        pkgs.vips
        pkgs.vips.dev
        pkgs.glib.dev
        pkgs.air
        pkgs.poppler-utils
        pkgsUnstable.claude-code
    ];
    extraConfigEnv = [
        "PKG_CONFIG_PATH=/lib/pkgconfig:/share/pkgconfig"
    ];
    configCmd = [ "/tmp/remote.sh" ];
    enableDebugTools = true;
    extraFakeRootCommands = ''
        cat <<'EOF' > "/etc/bashrc"
        if [[ $- != *i* ]]; then
            return
        fi
        if [[ "$INSIDE_EMACS" == *vterm* ]] && [[ "$INSIDE_EMACS" == *tramp* ]]; then
            export HISTFILE=~/.bash_history
        fi
        PS1="\[\033[0;32m\](${name})\[\033[0m\] \[\033[0;34m\]\W\[\033[0m\] \$ "
        EOF

        cat <<'EOF' > "/etc/profile"
        if [[ -f /tmp/env.sh ]]; then
            source /tmp/env.sh
        fi
        if [[ -f /etc/bashrc ]]; then
            source /etc/bashrc
        fi
        EOF

        cat <<'EOF' > "/tmp/remote.sh"
        #!/bin/sh
        export -p > "/tmp/env.sh"
        echo "$SSH_PUBKEY" > /srv/sshd/.ssh/authorized_keys
        chmod 600 /srv/sshd/.ssh/authorized_keys
        exec ${pkgs.openssh}/bin/sshd -f /srv/sshd/sshd_config -D
        EOF

        chmod +x "/tmp/remote.sh"

        # SSH setup
        mkdir -p /srv/sshd/ssh_host_keys
        mkdir -p /srv/sshd/.ssh
        chmod 700 /srv/sshd /srv/sshd/.ssh /srv/sshd/ssh_host_keys
        ${pkgs.openssh}/bin/ssh-keygen -t ed25519 -f /srv/sshd/ssh_host_keys/ssh_host_ed25519_key -N "" -q

        cat <<EOF > /srv/sshd/sshd_config
        Port 2222
        AddressFamily inet
        HostKey /srv/sshd/ssh_host_keys/ssh_host_ed25519_key
        AuthorizedKeysFile /srv/sshd/.ssh/authorized_keys
        PidFile /tmp/sshd.pid
        PubkeyAuthentication yes
        PasswordAuthentication no
        KbdInteractiveAuthentication no
        UsePAM no
        PermitEmptyPasswords no
        EOF

        chmod 600 /srv/sshd/sshd_config
        chown -R runner:runner /srv/sshd
    '';
}
