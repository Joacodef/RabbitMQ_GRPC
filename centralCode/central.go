package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"strconv"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"

	pb "github.com/Kendovvul/Ejemplo/Proto"
	"google.golang.org/grpc"
)

func failOnError(err error, msg string) {
	if err != nil {
		log.Panicf("%s: %s", msg, err)
	}
}

func init() {
	c := make(chan os.Signal)
	signal.Notify(c, os.Interrupt, os.Interrupt)
	go func() {
		<-c
		println("Terminando ejecución")
		connS, _ := grpc.Dial("localhost:50051", grpc.WithInsecure())
		connS2, _ := grpc.Dial("localhost:50052", grpc.WithInsecure())
		connS3, _ := grpc.Dial("localhost:50053", grpc.WithInsecure())
		connS4, _ := grpc.Dial("localhost:50054", grpc.WithInsecure())

		serviceCliente = pb.NewMessageServiceClient(connS)
		serviceCliente2 = pb.NewMessageServiceClient(connS2)
		serviceCliente3 = pb.NewMessageServiceClient(connS3)
		serviceCliente4 = pb.NewMessageServiceClient(connS4)

		serviceCliente.Terminar(context.Background(), &pb.Message{
			Body: "",
		})

		serviceCliente2.Terminar(context.Background(), &pb.Message{
			Body: "",
		})

		serviceCliente3.Terminar(context.Background(), &pb.Message{
			Body: "",
		})

		serviceCliente4.Terminar(context.Background(), &pb.Message{
			Body: "",
		})

		lab1 := "Laboratorio1;" + strconv.Itoa(solicitudes[0])
		lab2 := "Laboratorio2;" + strconv.Itoa(solicitudes[1])
		lab3 := "Laboratorio3;" + strconv.Itoa(solicitudes[2])
		lab4 := "Laboratorio4;" + strconv.Itoa(solicitudes[3])

		file, err := os.Create("SOLICITUDES.txt")
		if err != nil {
			panic(err)
		}
		defer file.Close()
		_, err = file.WriteString(lab1 + "\n")
		_, err = file.WriteString(lab2 + "\n")
		_, err = file.WriteString(lab3 + "\n")
		_, err = file.WriteString(lab4 + "\n")
		if err != nil {
			panic(err)
		}

		os.Exit(1)
	}()
}

var serviceCliente pb.MessageServiceClient
var serviceCliente2 pb.MessageServiceClient
var serviceCliente3 pb.MessageServiceClient
var serviceCliente4 pb.MessageServiceClient

var solicitudes [4]int

func main() {
	solicitudes[0] = 0
	solicitudes[1] = 0
	solicitudes[2] = 0
	solicitudes[3] = 0

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

	msgs, err := ch.Consume(
		q.Name, // queue
		"",     // consumer
		true,   // auto-ack
		false,  // exclusive
		false,  // no-local
		false,  // no-wait
		nil,    // args
	)
	failOnError(err, "Failed to register a consumer")

	var forever chan struct{}

	var Equipo1 int = 0 //Variables que indican a qué lab se asigna cada equipo
	var Equipo2 int = 0
	var equipoD int = 0

	go func() {
		// LECTURA DE COLA RABBITMQ:
		for d := range msgs {
			// Conexion GRPC CON LABS SEGUN MENSAJE RECIBIDO:

			//SE PUEDE AHORRAR CODIGO AQUI CONCATENANDO EN EL PORT EL NUM EN D.BODY
			if string(d.Body) == "1" { //LAB 1
				port := ":50051"                                               //puerto de la conexion con el laboratorio
				log.Println("Mensaje asíncrono de laboratorio 1 leído")        //obtiene el primer mensaje de la cola
				connS, err := grpc.Dial("localhost"+port, grpc.WithInsecure()) //crea la conexion sincrona con el laboratorio
				if err != nil {
					panic("No se pudo conectar con el servidor" + err.Error())
				}
				defer connS.Close()

				serviceCliente = pb.NewMessageServiceClient(connS) //cliente GRPC con lab 1

			} else if string(d.Body) == "2" { //LAB 2
				port := ":50052"                                                                 //puerto de la conexion con el laboratorio
				log.Println("Mensaje asíncrono de laboratorio 2 leído")                          //obtiene el primer mensaje de la cola
				connS, err := grpc.Dial("localhost"+port, grpc.WithInsecure()) //crea la conexion sincrona con el laboratorio
				if err != nil {
					panic("No se pudo conectar con el servidor" + err.Error())
				}
				defer connS.Close()

				serviceCliente2 = pb.NewMessageServiceClient(connS) //cliente GRPC con lab 2

			} else if string(d.Body) == "3" { //LAB 3
				port := ":50053"                                                                 //puerto de la conexion con el laboratorio
				log.Println("Mensaje asíncrono de laboratorio 3 leído")                          //obtiene el primer mensaje de la cola
				connS, err := grpc.Dial("localhost"+port, grpc.WithInsecure()) //crea la conexion sincrona con el laboratorio
				if err != nil {
					panic("No se pudo conectar con el servidor" + err.Error())
				}
				defer connS.Close()

				serviceCliente3 = pb.NewMessageServiceClient(connS) //cliente GRPC con lab 3

			} else if string(d.Body) == "4" { //LAB 4
				port := ":50054"                                                                 //puerto de la conexion con el laboratorio
				log.Println("Mensaje asíncrono de laboratorio 4 leído")                          //obtiene el primer mensaje de la cola
				connS, err := grpc.Dial("localhost"+port, grpc.WithInsecure()) //crea la conexion sincrona con el laboratorio
				if err != nil {
					panic("No se pudo conectar con el servidor" + err.Error())
				}
				defer connS.Close()

				serviceCliente4 = pb.NewMessageServiceClient(connS) //cliente GRPC con lab 4

			}

			//Revision si hay equipo disponible, variable equipoD indica qué equipo está disponible:
			equipoD = 0
			if Equipo1 == 0 {
				equipoD = 1
			} else if Equipo2 == 0 {
				equipoD = 2
			}

			//SI HAY EQUIPO DISPONIBLE => SE DEBE ASIGNAR:
			if equipoD != 0 {
				var err error

				//Tomar mensaje de la cola y convertirlo a int (indica a que lab se le asignará un equipo)

				str1 := string(d.Body)
				numLab, err := strconv.Atoi(str1)
				if err != nil {
					panic("Error de conversion de int " + err.Error())
				}
				println(numLab)

				//Variables Equipo1 y Equipo2 almacenan el número del lab al que serán asignados
				if equipoD == 1 {
					Equipo1 = numLab
				} else if equipoD == 2 {
					Equipo2 = numLab
				}

				//Enviar mensaje GRPC para avisar de equipo asignado:

				if numLab == 1 {
					_, err = serviceCliente.Intercambio(context.Background(), &pb.Message{
						Body: strconv.Itoa(equipoD),
					})
				} else if numLab == 2 {
					_, err = serviceCliente2.Intercambio(context.Background(), &pb.Message{
						Body: strconv.Itoa(equipoD),
					})
				} else if numLab == 3 {
					_, err = serviceCliente3.Intercambio(context.Background(), &pb.Message{
						Body: strconv.Itoa(equipoD),
					})
				} else if numLab == 4 {
					_, err = serviceCliente4.Intercambio(context.Background(), &pb.Message{
						Body: strconv.Itoa(equipoD),
					})
				}
				log.Println("Se envía Escuadrón " + strconv.Itoa(equipoD) + " a Laboratorio " + strconv.Itoa(numLab)) //respuesta del laboratorio
				if err != nil {
					panic("No se puede crear el mensaje " + err.Error())
				}

			}

			var equipoAsignado string

			for Equipo1 != 0 && Equipo2 != 0 {
				// Chequea estado de labs cada 5 segundos
				time.Sleep(5 * time.Second)

				if Equipo1 == 1 || Equipo2 == 1 {
					res, err := serviceCliente.Verificar(context.Background(), &pb.Message{
						Body: "Verificando Estado",
					})
					if err != nil {
						panic("No se puede crear el mensaje " + err.Error())
					}
					solicitudes[0]++
					if Equipo1 == 1 {
						equipoAsignado = "1"
					} else if Equipo2 == 1 {
						equipoAsignado = "2"
					}

					log.Println("Status Escuadrón " + equipoAsignado + ": " + res.Body)

					if res.Body == "LISTO" {
						if Equipo1 == 1 {
							Equipo1 = 0
							equipoD = 1
						} else if Equipo2 == 1 {
							Equipo2 = 0
							equipoD = 2
						}
						//Apagar servidor
						serviceCliente.Apagar(context.Background(), &pb.Message{
							Body: "Solicitud de retorno de escuadron a lab 1",
						})
						log.Println("Retorno a Central Escuadrón " + equipoAsignado + ", Conexión Laboratorio 1 Cerrada")
					}
				}
				if Equipo1 == 2 || Equipo2 == 2 {
					res, err := serviceCliente2.Verificar(context.Background(), &pb.Message{
						Body: "Verificando Estado",
					})
					if err != nil {
						panic("No se puede crear el mensaje " + err.Error())
					}
					solicitudes[1]++
					if Equipo1 == 2 {
						equipoAsignado = "1"
					} else if Equipo2 == 2 {
						equipoAsignado = "2"
					}

					log.Println("Status Escuadrón " + equipoAsignado + ": " + res.Body)
					if res.Body == "LISTO" {
						if Equipo1 == 2 {
							Equipo1 = 0
							equipoD = 1
						} else if Equipo2 == 2 {
							Equipo2 = 0
							equipoD = 2
						}
						//Apagar servidor
						serviceCliente2.Apagar(context.Background(), &pb.Message{
							Body: "Solicitud de retorno de equipo a lab 2",
						})
						log.Println("Retorno a Central Escuadrón " + equipoAsignado + ", Conexión Laboratorio 2 Cerrada")
					}
				}
				if Equipo1 == 3 || Equipo2 == 3 {
					res, err := serviceCliente3.Verificar(context.Background(), &pb.Message{
						Body: "Verificando Estado",
					})
					if err != nil {
						panic("No se puede crear el mensaje " + err.Error())
					}
					solicitudes[2]++
					if Equipo1 == 3 {
						equipoAsignado = "1"
					} else if Equipo2 == 3 {
						equipoAsignado = "2"
					}

					log.Println("Status Escuadrón " + equipoAsignado + ": " + res.Body)

					if res.Body == "LISTO" {
						if Equipo1 == 3 {
							Equipo1 = 0
							equipoD = 1
						} else if Equipo2 == 3 {
							Equipo2 = 0
							equipoD = 2
						}
						//Apagar servidor
						serviceCliente3.Apagar(context.Background(), &pb.Message{
							Body: "Solicitud de retorno de escuadron a lab 3",
						})
						log.Println("Retorno a Central Escuadrón " + equipoAsignado + ", Conexión Laboratorio 3 Cerrada")
					}
				}
				if Equipo1 == 4 || Equipo2 == 4 {
					res, err := serviceCliente4.Verificar(context.Background(), &pb.Message{
						Body: "Verificando Estado",
					})
					if err != nil {
						panic("No se puede crear el mensaje " + err.Error())
					}
					equipoAsignado = string(res.Body[0])
					solicitudes[3]++
					if Equipo1 == 4 {
						equipoAsignado = "1"
					} else if Equipo2 == 4 {
						equipoAsignado = "2"
					}

					log.Println("Status Escuadrón " + equipoAsignado + ": " + res.Body)

					if res.Body == "LISTO" {
						if Equipo1 == 4 {
							Equipo1 = 0
							equipoD = 1
						} else if Equipo2 == 4 {
							Equipo2 = 0
							equipoD = 2
						}
						//Apagar servidor
						serviceCliente4.Apagar(context.Background(), &pb.Message{
							Body: "Solicitud de retorno de escuadron a lab 4",
						})
						log.Println("Retorno a Central Escuadrón " + equipoAsignado + ", Conexión Laboratorio 4 Cerrada")
					}
				}
			}

			time.Sleep(5 * time.Second)
			if Equipo1 == 1 || Equipo2 == 1 {
				res, err := serviceCliente.Verificar(context.Background(), &pb.Message{
					Body: "Verificando Estado",
				})
				if err != nil {
					panic("No se puede crear el mensaje " + err.Error())
				}
				solicitudes[0]++
				if Equipo1 == 1 {
					equipoAsignado = "1"
				} else if Equipo2 == 1 {
					equipoAsignado = "2"
				}

				log.Println("Status Escuadrón " + equipoAsignado + ": " + res.Body)

				if res.Body == "LISTO" {
					if Equipo1 == 1 {
						Equipo1 = 0
						equipoD = 1
					} else if Equipo2 == 1 {
						Equipo2 = 0
						equipoD = 2
					}
					//Apagar servidor
					serviceCliente.Apagar(context.Background(), &pb.Message{
						Body: "Solicitud de retorno de escuadron a lab 1",
					})
					log.Println("Retorno a Central Escuadrón " + equipoAsignado + ", Conexión Laboratorio 1 Cerrada")
				}
			}
			if Equipo1 == 2 || Equipo2 == 2 {
				res, err := serviceCliente2.Verificar(context.Background(), &pb.Message{
					Body: "Verificando Estado",
				})
				if err != nil {
					panic("No se puede crear el mensaje " + err.Error())
				}
				solicitudes[1]++
				if Equipo1 == 2 {
					equipoAsignado = "1"
				} else if Equipo2 == 2 {
					equipoAsignado = "2"
				}

				log.Println("Status Escuadrón " + equipoAsignado + ": " + res.Body)
				if res.Body == "LISTO" {
					if Equipo1 == 2 {
						Equipo1 = 0
						equipoD = 1
					} else if Equipo2 == 2 {
						Equipo2 = 0
						equipoD = 2
					}
					//Apagar servidor
					serviceCliente2.Apagar(context.Background(), &pb.Message{
						Body: "Solicitud de retorno de equipo a lab 2",
					})
					log.Println("Retorno a Central Escuadrón " + equipoAsignado + ", Conexión Laboratorio 2 Cerrada")
				}
			}
			if Equipo1 == 3 || Equipo2 == 3 {
				res, err := serviceCliente3.Verificar(context.Background(), &pb.Message{
					Body: "Verificando Estado",
				})
				if err != nil {
					panic("No se puede crear el mensaje " + err.Error())
				}
				solicitudes[2]++
				if Equipo1 == 3 {
					equipoAsignado = "1"
				} else if Equipo2 == 3 {
					equipoAsignado = "2"
				}

				log.Println("Status Escuadrón " + equipoAsignado + ": " + res.Body)

				if res.Body == "LISTO" {
					if Equipo1 == 3 {
						Equipo1 = 0
						equipoD = 1
					} else if Equipo2 == 3 {
						Equipo2 = 0
						equipoD = 2
					}
					//Apagar servidor
					serviceCliente3.Apagar(context.Background(), &pb.Message{
						Body: "Solicitud de retorno de escuadron a lab 3",
					})
					log.Println("Retorno a Central Escuadrón " + equipoAsignado + ", Conexión Laboratorio 3 Cerrada")
				}
			}
			if Equipo1 == 4 || Equipo2 == 4 {
				res, err := serviceCliente4.Verificar(context.Background(), &pb.Message{
					Body: "Verificando Estado",
				})
				if err != nil {
					panic("No se puede crear el mensaje " + err.Error())
				}
				solicitudes[3]++
				if Equipo1 == 4 {
					equipoAsignado = "1"
				} else if Equipo2 == 4 {
					equipoAsignado = "2"
				}

				log.Println("Status Escuadrón " + equipoAsignado + ": " + res.Body)

				if res.Body == "LISTO" {
					if Equipo1 == 4 {
						Equipo1 = 0
						equipoD = 1
					} else if Equipo2 == 4 {
						Equipo2 = 0
						equipoD = 2
					}
					//Apagar servidor
					serviceCliente4.Apagar(context.Background(), &pb.Message{
						Body: "Solicitud de retorno de escuadron a lab 4",
					})
					log.Println("Retorno a Central Escuadrón " + equipoAsignado + ", Conexión Laboratorio 4 Cerrada")
				}
			}

		}

	}()

	log.Printf(" [*] Waiting for messages. To exit press CTRL+C")
	<-forever
}
