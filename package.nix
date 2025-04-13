{ buildGoModule, lib }:
let
  version = "0.1.0";
in
buildGoModule {
  pname = "mpd-mcp-server";
  inherit version;
  vendorHash = "sha256-pJbkwUiAGjGiKBYiP21Ifo/PIfdlBFbmVUMKehPs+f4=";
  src = lib.cleanSource ./.;
  ldflags = [
    "-s"
    "-w"
    "-X main.version=${version}"
  ];
  meta = {
    description = "MCP Server for Music Player Daemon (MPD)";
    homepage = "https://github.com/gamoutatsumi/mpd-mcp-server";
    license = lib.licenses.mit;
  };
}
