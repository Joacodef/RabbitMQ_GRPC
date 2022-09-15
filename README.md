# Comunicación Sincrónica y Asincrónica con GRPC y RabbitMQ

# Ejecución de archivos:

-Se recomienda iniciar primero los laboratorios, ya que si hay un mensaje en la cola de RabbitMQ, pero la central no logra establecer la conexión GRPC con el laboratorio, la central terminará su ejecución.

-La central genera el archivo SOLICITUDES.txt en la misma carpeta donde se encuentran todos los archivos.

-Para terminar la ejecución de los programas, se debe hacer CTRL+C en la central, lo que debería terminar la ejecución en todas las otras máquinas (puede que haya casos muy particulares en que algún laboratorio no termine su ejecución, en cuyo caso, hacer CTRL+C en ese laboratorio).

-En la máquina de la central y laboratorio 1, se puede hacer "make clean" para limpiar el archivo de solicitudes.

# Explicación general del código:

El código funciona principalmente en base al loop que existe en central, toda la comunicación con laboratorios ocurre en ese loop.

Iterativamente, se revisan los mensajes que se hayan recogido en la cola de RabbitMQ, y de acuerdo a qué laboratorio haya enviado el mensaje que lea, se establece una conexión GRPC con él. Si no hay equipos disponibles, el programa se "mantiene" en ese mensaje RabbitMQ mientras sigue iterando en un loop interno, hasta que se libere un equipo.

La central envía 4 tipos de mensajes GRPC: 

-Indicar que se envió equipo, e "iniciar" las comunicaciones.

-Verificar el estado del laboratorio.

-Recibir el equipo de vuelta y terminar el servidor GRPC del lab.

-Terminar el funcionamiento del laboratorio.