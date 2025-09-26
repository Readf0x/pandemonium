rec {
  description = "Description for the project";

  inputs = {
    flake-parts.url = "github:hercules-ci/flake-parts";
    nixpkgs.url = "github:NixOS/nixpkgs/nixos-unstable";
  };

  outputs = inputs @ {flake-parts, ...}:
    flake-parts.lib.mkFlake {inherit inputs;} {
      systems = ["x86_64-linux"];
      perSystem = {
        system,
        pkgs,
        ...
      }: let
        info = {
          projectName = "slander"; # Don't forget to change this!
          # You can set the module name as well
          # moduleName = "github.com/example/${projectName}";
        };
      in
        (
          {
            projectName,
            moduleName ? projectName,
          }: rec {
            devShells.default = pkgs.mkShell {
              packages = with pkgs; [
                go
                delve
                air
                templ
                nodejs
                sqlite
                sqlitebrowser
              ];
            };
            packages = {
              ${projectName} = pkgs.buildGoModule rec {
                name = projectName;
                pname = name;
                version = "0.1";

                src = ./.;

                vendorHash = "sha256-AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA=";

                meta = {
                  inherit description;
                  # homepage = "";
                  # license = lib.licenses.;
                  # maintainers = with lib.maintainers; [  ];
                };
              };
              default = packages.${projectName};
              dev = pkgs.writeShellScriptBin "dev" ''

                trap 'kill $(jobs -p)' EXIT
                mkdir tmp/bin -p &>/dev/null

                templ generate \
                --watch \
                --proxy="http://localhost:3000" \
                --open-browser=false -v &

                air \
                --build.cmd "go build -o tmp/bin ./cmd/server" \
                --build.bin "tmp/bin/server" \
                --build.include_ext "go" \
                --build.stop_on_error "false" \
                --misc.clean_on_exit true

              '';
            };
          }
        )
        info;
    };
}
