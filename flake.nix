{
  description = "XDG Portal Terminal App Chooser";

  inputs = {
    nixpkgs.url = "github:NixOS/nixpkgs/nixos-unstable";
    flake-utils.url = "github:numtide/flake-utils";
  };

  outputs = { self, nixpkgs, flake-utils }:
    flake-utils.lib.eachDefaultSystem (system:
      let
        pkgs = nixpkgs.legacyPackages.${system};
      in
      {
        packages.default = pkgs.buildGoModule {
          pname = "xdg-portal-termappchooser";
          version = "0.1.0";
          
          src = ./.;

          nativeBuildInputs = with pkgs; [
            pkg-config
          ];

          buildInputs = with pkgs; [
            gtk3
            glib
          ];

          postInstall = ''
            cp -r data/share $out/share
          '';
          
          vendorHash = "sha256-onlp5RPl9osHbv0RT21oNINdjuSRmIVthPy7dlNLK6Q=";
          
          meta = with pkgs.lib; {
            description = "XDG Desktop Portal AppChooser and OpenURI implementation with fuzzel integration";
            homepage = "https://github.com/MartinLoeper/xdg-portal-termappchooser";
            license = licenses.mit;
            maintainers = [ ];
          };
        };

        devShells.default = pkgs.mkShell {
          buildInputs = with pkgs; [
            go
            gopls
            gotools
            go-tools
            fuzzel
            dbus
            gtk3
            glib
            pkg-config
            gcc
            gnumake
          ];
          
          shellHook = ''
            echo "XDG Portal Terminal App Chooser Development Environment"
            echo "Available commands:"
            echo "  ./build.sh - Build the application"
            echo "  ./xdg-portal-termappchooser - Run the built binary"
            echo "  nix build - Build with Nix"
            echo "  nix run - Build and run with Nix"
            echo ""
            echo "OpenURI + AppChooser implementation with GIO and libnotify"
          '';
        };
      }) // {
        nixosModules.default = { config, lib, pkgs, ... }: 
          with lib;
          let
            cfg = config.services.xdg-portal-termappchooser;
            termappchooser = self.packages.${pkgs.system}.default;
          in
          {
            options.services.xdg-portal-termappchooser = {
              enable = mkEnableOption "XDG Portal Terminal App Chooser";
            };

            config = mkIf cfg.enable {
              xdg.portal = {
                enable = true;
                extraPortals = [ termappchooser ];
                config = {
                  hyprland = {
                    "org.freedesktop.impl.portal.AppChooser" = "termappchooser";
                    "org.freedesktop.impl.portal.OpenURI" = "termappchooser";
                  };
                  common = {
                    "org.freedesktop.impl.portal.AppChooser" = "termappchooser";
                    "org.freedesktop.impl.portal.OpenURI" = "termappchooser";
                  };
                };
              };

              systemd.user.services.xdg-desktop-portal-termappchooser = {
                unitConfig = {
                  After = ["graphical-session.target"];
                  PartOf = ["graphical-session.target"];
                  Description = "XDG Desktop Portal Terminal App Chooser Service";
                };

                wantedBy = ["graphical-session.target"];
                
                serviceConfig = {
                  ExecStart = "${lib.getExe termappchooser}";
                  Restart = "on-failure";
                  Type = "dbus";
                  BusName = "org.freedesktop.impl.portal.desktop.termappchooser";
                  Slice = "session.slice";
                };
              };
            };
          };
      };
}