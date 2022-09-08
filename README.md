# TP0: Docker + Comunicación + Sincronización
## Barreneche Franco

## Ejecución

Para la mayoría de los ejercicios alcanza con ejecutar `sudo make docker-compose-up`.

El script solicitado en el ejercicio 1.1 se ubica en la carpeta utils, en la rama **ej1**.
Al ejecutar `sh gen_compose <cant>` donde cant es la cantidad de clientes deseados se obtiene un archivo docker-compose-dev.yaml que concatena el contenido de prefix.txt, repeat.txt y sufix.txt (en ese órden), repitiendo los contenidos de repeat.txt y reemplazando la palabra REPEAT por el índice de la iteración. Con los archivos por defecto ya puede generarse un archivo docker-compose-dev.yaml válido.

Los archivos de configuración, datasets y la salida del programa (cuando corresponde) se reúnen en la carpeta **.data**.

## Protocolo

### Representación de datos
Se especifica la representación para algunos tipos de datos primitivos, con los que pueden construirse structs más complejos. Éstos son: 
- uint32: 4 bytes, big endian.
- uint64: 8 bytes, big endian.
- bool: 1 byte, falso si es 0, de lo contrario verdadero.
- string: en utf-8. Su tamaño es variable y se indica con un uint32 que la antecede. 

Los vectores se envían indicando la cantidad de elementos del mismo con un uint32 al comienzo y luego anexando directamente la representación en bytes de cada elemento.

### Flujo de mensajes

El cliente siempre inicia la comunicación con un uint32 que representa su intención.

AskWinners: Se envía un vector de registros de personas. El servidor responde con booleanos individuales en el mismo órden en que se encuentran los elementos del vector, de modo que el quinto booleano corresponde al quinto elemento e indica si esa persona ganó o no.

AskAmount: Basta con comunicar la intención. El servidor responde con un par (uint32, bool) que indican el tamaño y si el resultado es o no parcial. Si el resultado es parcial el cliente reintenta la operación tras N segundos.

El fin de la comunicación se infiere por la pérdida de conexión cuando el servidor espera una nueva intención.

