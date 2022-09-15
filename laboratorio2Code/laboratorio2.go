package main

import (
	"context"
	"log"
	"math/rand"
	"net"
	"os"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"

	pb "github.com/Kendovvul/Ejemplo/Proto"
	"google.golang.org/grpc"
)

type server struct {
	pb.UnimplementedMessageServiceServer
}

var equipoAsig string

func (s *server) Intercambio(ctx context.Context, msg *pb.Message) (*pb.Message, error) {
	log.Printf("Llega Escuadrón " + msg.Body + ", conteniendo estallido...")
	equipoAsig = string(msg.Body)
	return &pb.Message{Body: "Escuadrón Recibido Lab 2"}, nil
}

func (s *server) Verificar(ctx context.Context, msg *pb.Message) (*pb.Message, error) {
	s1 := rand.NewSource(time.Now().UnixNano())
	r1 := rand.New(s1)
	var bodyStr string
	if r1.Intn(100) < 60 {
		bodyStr = "LISTO"
		log.Printf("Revisando estado Escuadrón: LISTO")
	} else {
		bodyStr = "NO LISTO"
		log.Printf("Revisando estado Escuadrón: NO LISTO")
	}
	return &pb.Message{Body: bodyStr}, nil
}

func (s *server) Apagar(ctx context.Context, msg *pb.Message) (*pb.Message, error) {
	log.Printf("Estallido Contenido, Escuadrón " + equipoAsig + " Retornando...")
	equipoAsig = "0"
	serv.Stop()
	return &pb.Message{Body: "Equipo retornado"}, nil
}

func (s *server) Terminar(ctx context.Context, msg *pb.Message) (*pb.Message, error) {
	log.Printf("Señal de termino recibida. Terminando ejecución")
	serv.Stop()
	os.Exit(1)
	return &pb.Message{Body: ""}, nil
}

func failOnError(err error, msg string) {
	if err != nil {
		log.Panicf("%s: %s", msg, err)
	}
}

var serv *grpc.Server

func main() {

	conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
	failOnError(err, "Failed to connect to RabbitMQ")
	defer conn.Close()

	ch, err := conn.Channel()
	failOnError(err, "Failed to open a channel")
	defer ch.Close()

	q, err := ch.QueueDeclare(
		"Central", // name
		false,     // durable
		false,     // delete when unused
		false,     // exclusive
		false,     // no-wait
		nil,       // arguments
	)

	failOnError(err, "Failed to declare a queue")

	var estallido bool = false

	//INICIO CICLO POSIBILIDAD ESTALLIDO EN LABORATORIO

	for {

		time.Sleep(5 * time.Second)

		s1 := rand.NewSource(time.Now().UnixNano())
		r1 := rand.New(s1)

		if r1.Intn(100) < 80 {
			estallido = true
		}

		if estallido {
			log.Printf("Analizando estado laboratorio: ESTALLIDO")

			//ENVIO MENSAJE RABBITMQ
			body := "2"
			err = ch.Publish(
				"",     // exchange
				q.Name, // routing key
				false,  // mandatory
				false,  // immediate
				amqp.Publishing{
					ContentType: "text/plain",
					Body:        []byte(body),
				})

			failOnError(err, "Failed to publish a message")
			log.Printf("SOS Enviado a Central. Esperando respuesta...")

			//CONEXION GRPC

			listener, err := net.Listen("tcp", ":50052") //conexion sincrona
			if err != nil {
				panic("La conexion no se pudo crear" + err.Error())
			}

			serv = grpc.NewServer()
			pb.RegisterMessageServiceServer(serv, &server{})
			if err = serv.Serve(listener); err != nil {
				panic("El server no se pudo iniciar" + err.Error())
			}

		} else {
			log.Printf("Analizando estado laboratorio: OK")
		}
	}
}
