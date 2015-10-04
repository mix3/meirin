GIT_VER := $(shell git describe --tags)
DATE    := $(shell date +%Y-%m-%dT%H:%M:%S%z)
OWNER   := "mix3"

clean:
	rm -rf pkg/*

binary: clean
	gox -osarch="linux/amd64 darwin/amd64" \
		-output "pkg/{{.Dir}}-${GIT_VER}-{{.OS}}-{{.Arch}}" \
		-ldflags "-X main.version ${GIT_VER} -X main.buildDate ${DATE}"

package: binary
	cd pkg && find . -name "*${GIT_VER}*" -type f \
		-exec mkdir -p meirin \; \
		-exec cp {} meirin/meirin \; \
		-exec zip -r {}.zip meirin \; \
		-exec rm -rf meirin \;
