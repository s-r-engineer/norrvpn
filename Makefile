GoLinux = linux
GoArch386 = 386
GoArch64 = amd64
CGO = 0
N = norrvpn

ifeq ($(DEBUG),true)
    LDFLAGS = -ldflags "-X 'main.debug=1'"
endif

build: build_linux

build_linux: build_linux_386 build_linux_x64

build_linux_386:
	GOOS=$(GoLinux) GOARCH=$(GoArch386) CGO_ENABLED=$(CGO) go build $(LDFLAGS) -o ./builds/$(N)_$(GoLinux)_$(GoArch386) .

build_linux_x64:
	echo GOOS=$(GoLinux) GOARCH=$(GoArch64) CGO_ENABLED=$(CGO) go build $(LDFLAGS) -o ./builds/$(N)_$(GoLinux)_$(GoArch64) .
	GOOS=$(GoLinux) GOARCH=$(GoArch64) CGO_ENABLED=$(CGO) go build $(LDFLAGS) -o ./builds/$(N)_$(GoLinux)_$(GoArch64) .

install_x64: build_linux_x64
	sudo cp ./builds/$(N)_$(GoLinux)_$(GoArch64) /usr/local/bin/$(N)