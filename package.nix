{ buildGoModule, lib }:
let
  version = "0.1.1";
in
buildGoModule {
  pname = "mpd-mcp-server";
  inherit version;
  vendorHash = "sha256-yox4oz9VQlE8dolPmu4Sl91OjbPbv9ymSiBZIEOjD4I=";
  src = lib.cleanSource ./.;
  ldflags = [
    "-s"
    "-w"
    "-X github.com/harakeishi/curver.Version=v${version}"
  ];
  meta = {
    description = "MCP Server for Music Player Daemon (MPD)";
    homepage = "https://github.com/gamoutatsumi/mpd-mcp-server";
    license = lib.licenses.mit;
  };
}
