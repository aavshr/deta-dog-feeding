build-mermaid:
	mkdir -p mermaid
	cd mermaid-src && yarn install && yarn build
	cp mermaid-src/package.json mermaid/
	cp -R mermaid-src/build/*  mermaid/

clean:
	rm -rf build-mermaid