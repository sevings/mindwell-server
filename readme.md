# Build and run server
1. Install PostgreSQL:
```
sudo apt-get install postgresql postgresql-contrib
# set `default_text_search_config = 'pg_catalog.russian'`
sudo nano /etc/postgresql/14/main/postgresql.conf
sudo systemctl start postgresql.service
sudo -u postgres createuser -P --interactive
```
2. Install RUM:
```
sudo apt-get install systemtap-sdt-dev postgresql-server-dev-14
git clone https://github.com/postgrespro/rum
cd rum
make USE_PGXS=1
make USE_PGXS=1 install
```
3. Install the latest release of [go-swagger](https://github.com/go-swagger/go-swagger/releases):
```
wget https://github.com/go-swagger/go-swagger/releases/download/v0.30.5/swagger_linux_amd64
sudo mv swagger_linux_amd64 /usr/bin/swagger
```
4. Install Go: `sudo snap install go --classic`
6. Clone mindwell-server:
```
mkdir -p ~/go/src
cd ~/go/src
git clone https://github.com/sevings/mindwell-server.git
cd mindwell-server
```
6. Generate code: `./scripts/generate.sh`
7. Create DB:
```
psql -c 'create database mindwell'
psql -d mindwell -q -f scripts/mindwell.sql
```
8. Configure:
```
cp configs/server.sample.toml configs/server.toml
nano configs/server.toml
```
9. Run tests: `go test ./test/ --failfast`
10. Run server: `go run ./cmd/mindwell-server/ --port 8000`

# Build and run images
11. Install dependencies: `sudo apt-get install libvips-dev`
12. Configure:
```
cp configs/images.sample.toml configs/images.toml
nano configs/images.toml
```
13. Run images: `go run ./cmd/mindwell-images-server/ --port 8888`
