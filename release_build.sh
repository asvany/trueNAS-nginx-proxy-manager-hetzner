BIN_NAME=trueNAS-nginx-proxy-manager-hetzner.static.upx
./build.sh && cp trueNAS-nginx-proxy-manager-hetzner.static $BIN_NAME && strip $BIN_NAME && upx $BIN_NAME
