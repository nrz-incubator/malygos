{ pkgs ? import <nixpkgs> { },
}:

pkgs.mkShell {
  LOCALE_ARCHIVE = "${pkgs.glibcLocales}/lib/locale/locale-archive";
  env.LANG = "C.UTF-8";
  env.LC_ALL = "C.UTF-8";
  env.GOPATH = "${builtins.getEnv "PWD"}/.go";

  packages = [
	pkgs.go
	pkgs.delve
  ];
  shellHook = ''
    ${(import ./pre-commit.nix).pre-commit-check.shellHook}
  '';
}
