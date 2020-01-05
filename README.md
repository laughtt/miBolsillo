### JSON 

Simple json parser, el archivo main.go que esta situado en cmd/mibolsillo maneja el handler (CreateInvoice) que esta situado en la carpeta api. Schemes contiene el struct con el cual va a recibir el JSON (Message) , una vez parseado y filtrado de errores , se creara un array del struct llamado (Responses) con pointers a (Response). El JSON es descerializado con json.decoder en la funcion (decodeEncodeJSONBody) , y se guarda en la memoria. Cada objecto se guarda en un mapa de ids (mapa[string]*Response por si hay diferentes ids en un mismo array) al final en la funcion(CreateInvoice) se haze un loop para buscar cada diferente id y se entrega un array de la respuesta.
Se usa para evitar la repeticion de codigo un struct en tools llamado (MalformedRequest) , y un catch que categoriza y devuelve todos los errores que se pueden ver.
### Prerequisitess

```
Docker , Golang 1.x  
```

### Installing

For running the server you just need to :

```
docker build .
docker-compose deploy -c docker-compose.yml bolsillo
```

### Coverage test

Testear el codigo. Lee lo que esta en la carpeta test/jsons(badTest / big(limite de 1mb) / correct / outputCheck), 

```
go test ./api -coverprofile cover
go tool cover -html=cover
```