### JSON 

Simple json parser, el archivo main.go que esta situado en cmd/mibolsillo maneja el handler (CreateInvoice) que esta situado en la carpeta api.
El archivo Schemes contiene el struct con el cual va a recibir el JSON (Message) , una vez parseado y filtrado de errores , se creara un array del struct llamado (Responses) con pointers a (Response)\n.
El archivo handler contiene la logica para encodear y decodear el JSON.
El JSON es descerializado con json.decoder en la funcion (decodeEncodeJSONBody) , y se guarda en la memoria. Cada objecto se guarda en un mapa de ids ("mapa[string]*Response") al final en la funcion (CreateInvoice) se haze un loop para buscar cada diferente id y se entrega un array de la respuesta.
Se usa para evitar la repeticion de errores un struct en tools llamado (MalformedRequest) , y un catch que categoriza y devuelve todos los errores que se pueden ver.

### Prerequisitess

```
Docker , Golang 1.x  
```

### Installing

Para correr el server ejecutar los siguientes comandos
```
docker build .
docker stack deploy -c docker-compose.yml bolsillo
```

### Output test 
En la carpeta test esta el ejecutable para ejecutar los json y enviarlos , la forma de hacerlo es : ./test [nombre_del_directorio_de_json]
 ```
 ./test [DIR_NAME] [DIR_NAME] ...
```

### Coverage test

Testear el codigo que esta en la carpeta test/jsons (badTest / big(limite de 1mb) / correct / outputCheck) y muestra un html , del codigo usado. 

```
go test ./api -coverprofile cover
go tool cover -html=cover
```