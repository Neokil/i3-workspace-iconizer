install:
	go build -o i3-workspace-iconizer
	chmod a+x i3-workspace-iconizer
	ln -sf $(shell pwd)/i3-workspace-iconizer /usr/bin/i3-workspace-iconizer
	cp -n config.json ~/.i3-workspace-iconizer