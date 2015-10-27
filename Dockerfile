FROM ubuntu:trusty

RUN	echo "deb http://repo.reverbrain.com/trusty/ current/amd64/" > /etc/apt/sources.list.d/reverbrain.list && \
	echo "deb http://repo.reverbrain.com/trusty/ current/all/" >> /etc/apt/sources.list.d/reverbrain.list && \
	apt-get install -y curl tzdata && \
	cp -f /usr/share/zoneinfo/posix/W-SU /etc/localtime && \
	curl http://repo.reverbrain.com/REVERBRAIN.GPG | apt-key add - && \
	apt-get update && \
	apt-get upgrade -y && \
	apt-get install -y git elliptics-client elliptics-dev g++

RUN 	VERSION=go1.5.1 && \
	curl -O https://storage.googleapis.com/golang/$VERSION.linux-amd64.tar.gz && \
	tar -C /usr/local -xf $VERSION.linux-amd64.tar.gz && \
	rm -f $VERSION.linux-amd64.tar.gz

RUN	git config --global user.email "zbr@ioremap.net" && \
	git config --global user.name "Evgeniy Polyakov" && \
	#cd /root && \
	#git clone https://go.googlesource.com/go gosource && \
	#cd gosource && \
	#git fetch https://go.googlesource.com/go refs/changes/44/13944/3 && git checkout -b stack_changes FETCH_HEAD && \
	#cd src && \
	#GOROOT_BOOTSTRAP=/usr/local/go ./all.bash && \
	#export PATH=$PATH:/root/gosource/bin && \
	export PATH=$PATH:/usr/local/go/bin && \
	mkdir /root/go && \
	export GOPATH=/root/go && \
	go get github.com/bioothod/elliptics-go/elliptics && \
	cd /root/go/src/github.com/bioothod/elliptics-go/elliptics && \
	git checkout master && \
	git pull && \
	git branch -v && \
	go install && \
	echo "Go binding has been updated" && \
	go get github.com/bioothod/backrunner && \
	cd /root/go/src/github.com/bioothod/backrunner && \
	git checkout master && \
	git pull && \
	git branch -v && \
	make && \
	go install && \
	echo "Backrunner has been updated" ;\
    	rm -rf /var/lib/apt/lists/*

EXPOSE 9090 80 443