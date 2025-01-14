{
  description = "Utility for bumping and pushing git tags";

  inputs = {
    nixpkgs.url = "github:NixOS/nixpkgs/nixos-unstable";
    flake-utils.url = "github:numtide/flake-utils";
  };

  outputs = { self, nixpkgs, flake-utils }:
    flake-utils.lib.eachDefaultSystem (system:
      let
        pkgs = nixpkgs.legacyPackages.${system};

        lastTag = "v0.2.0";

        revision =
          if (self ? shortRev)
          then "${self.shortRev}"
          else "${self.dirtyShortRev or "dirty"}";

        # Add the commit to the version string for flake builds
        version = "${lastTag}";

        # Run `devbox run update-hash` to update the vendor-hash
        vendorHash =
          if builtins.pathExists ./vendor-hash
          then builtins.readFile ./vendor-hash
          else "";

        buildGoModule = pkgs.buildGo123Module;

      in
      {
        inherit self;
        packages.default = buildGoModule {
          pname = "bump";
          inherit version vendorHash;

          src = ./.;

          subpackage = [ ./cmd/bump ];

          ldflags = [
            "-s"
            "-w"
            "-X github.com/mrvinkel/bump/cmd/bump/internal.BumpVersion=${version}"
          ];

          # Disable tests if they require network access or are integration tests
          doCheck = false;

          nativeBuildInputs = [ pkgs.installShellFiles ];

          meta = with pkgs.lib; {
            description = "Utility for bumping and pushing git tags";
            homepage = "https://github.com/mrvinkel/bump";
            license = licenses.unlicense;
            maintainers = with maintainers; [ mrvinkel ];
          };
        };
      }
    );
}