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
          
          vendorHash = "sha256-VtGat4ek0ij8GOx68MQPNFtBuansj/d1GCOgfLOiGwM=";
          
          meta = with pkgs.lib; {
            description = "XDG Desktop Portal AppChooser implementation with fuzzel integration";
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
          ];
          
          shellHook = ''
            echo "XDG Portal Terminal App Chooser Development Environment"
            echo "Available commands:"
            echo "  ./build.sh - Build the application"
            echo "  ./xdg-portal-termappchooser - Run the built binary"
            echo "  nix build - Build with Nix"
            echo "  nix run - Build and run with Nix"
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
                };
              };

              systemd.user.services.xdg-desktop-portal-termappchooser = {
                after = ["graphical-session.target"];
                wantedBy = ["graphical-session.target"];
                serviceConfig = {
                  ExecStart = "${lib.getExe termappchooser}";
                  Restart = "on-failure";
                };
              };
            };
          };
      };
}