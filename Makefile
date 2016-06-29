all: generator

generator: gen.go
	go build -o generator gen.go

website: generator
	./generator

compare:
	[[ ! -d prod~ ]] && mkdir prod~
	ncftpget -R -u voilokov -p `pass rsync/cooking` ftp.voilokov.com prod~/ '/public_html/cooking/'
	diff prod~ gen

deploy:
	# need rename previous dir
	# ncftpput -R -u voilokov -p `pass rsync/cooking` ftp.voilokov.com '/public_html/cooking/' gen/*
