{ pkgs ? import <nixpkgs> {} }:
pkgs.mkShell {
  nativeBuildInputs = with pkgs; [
    # Python
    python3

    # Go
    go

    # CLI tools
    sql-migrate
    goose
  ];
}
