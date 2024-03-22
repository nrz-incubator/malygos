let
  nix-pre-commit-hooks = import (builtins.fetchTarball "https://github.com/cachix/pre-commit-hooks.nix/tarball/ea758da1a6dcde6dc36db348ed690d09b9864128");
in {
  pre-commit-check = nix-pre-commit-hooks.run {
    src = ./.;
    # If your hooks are intrusive, avoid running on each commit with a default_states like this:
    # default_stages = ["manual" "push"];
    hooks = {
      shellcheck.enable = true;
      statix.enable = true;
	  gofmt.enable = true;
	  # disabled for now due to wrong golang version
	  # govet.enable = true;
    };
    settings = {};
  };
}