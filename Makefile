default:
	go build
	./diary --env=production --domain=larissavoigt.com --port=80 \
		--facebook-id=1629858967301577 --facebook-secret=36b8b62d4a6d62f3e845a2682698749d
