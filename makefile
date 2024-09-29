tool=querylyctl
company=DanielVeraM
binPath = bin
executionDir:=/usr/local/bin
buildPath = build
outFile = $(buildPath)/$(tool)
tarName = clitool.tar.xz
tarFile = $(buildPath)/$(tarName)
installer_template = installer/installer.sh
installer_name=setup.sh
installer = $(buildPath)/$(installer_name)
conf = ./cfg-dev.json
tmpFile = $(buildPath)/tmpload
manifest = pkg/config/manifest.go
config_ctl = pkg/config/constants.go
cli_version := $(shell git describe)
pwdPath := $(shell pwd)

# deploy vars
host = root@127.0.0.1
key = /home/user/.ssh/yourkey
remoteDir = /tmp
remoteZip = $(remoteDir)/$(tool).zip
remoteScript = $(remoteDir)/setup.sh

all: build

.PHONY: update-dep
update-dep:
	@echo "[*] downloading dependencies"
	@go mod download
	@go mod tidy

.PHONY: compile
compile: update-dep
	@echo "[*] setting current version "$(cli_version)
	@sed -i '' "s/{{__VERSION_PLACEHOLDER__}}/$(cli_version)/g" $(manifest)
	@sed -i '' "s/{{__COMMAND_PLACEHOLDER__}}/$(tool)/g" $(config_ctl)
	@sed -i '' "s/{{__COMPANY_PLACEHOLDER__}}/$(company)/g" $(config_ctl)
	@sed -i '' "s|{{__INSTALLATION_PLACEHOLDER__}}|$(executionDir)|g" $(config_ctl)
	@echo "[*] compiling go code"
	@go build -o $(buildPath) ./cmd/...
	@mv  $(buildPath)/cmd $(outFile)
	@chmod a+x $(outFile)
	@echo "[*] rolling back manifest"
	@git checkout -- $(manifest)
	@git checkout -- $(config_ctl)

.PHONY: compress
compress:
	@echo "[*] compressing dependencies"
	@cd  $(buildPath); tar cJf $(tarName) $(tool)
	@cp $(installer_template) $(installer)
	@echo "[*] adding dependencies to installer"
	@touch $(tmpFile)
	@cat $(tarFile) >>  $(tmpFile)
	@cat $(tmpFile) >> $(installer)
	@echo "[*] applying permission to installer"
	@chmod a+x $(installer)
	@mv $(installer) $(binPath)/$(installer_name)
	@echo "[*] zipping installer"
	@cd ./$(binPath); zip $(tool).zip $(installer_name)
	@echo "[*] succesfully created installer for "$(tool)

.PHONY: clean
clean:
	@echo "[*] removing tmp files"
	@rm -f -r $(buildPath)

.PHONY: setUp
setUp:
	@echo "[*] building "$(tool)"..."
	@mkdir -p $(buildPath)
	@mkdir -p $(binPath)

.PHONY: build
build: setUp compile compress clean

.PHONY: build-run
build-run: build install
	@echo "[*] Executing the version of "$(tool)
	@$(tool) version

.PHONY: install
install: build
	@cd $(binPath); sudo ./setup.sh

.PHONY: deploy
deploy:
	@echo "[*] deploying to remote..."
	@echo "[-] transferring zip file..."
	scp $(pwdPath)/$(binPath)/$(tool).zip $(host):$(remoteDir)
	@echo "[-] unzipping remote zip..."
	ssh $(host) -i $(key) unzip -u -q  $(remoteZip) -d /tmp
	@echo "[-] installing..."
	ssh $(host) -i $(key)  "sudo /tmp/setup.sh; rm /tmp/setup.sh; $(tool) version;"

.PHONY: build-deploy
build-deploy: build deploy
	@echo "--- built and deployed ---"
