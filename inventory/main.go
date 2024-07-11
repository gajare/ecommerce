package main

import (
    "context"
    "encoding/json"
    "log"
    "github.com/segmentio/kafka-go"
)

func main() {
    kafkaReader := kafka.NewReader(kafka.ReaderConfig{
        Brokers: []string{"kafka:9092"},
        Topic:   "purchase-events",
        GroupID: "inventory-service",
    })

    for {
        msg, err := kafkaReader.ReadMessage(context.Background())
        if err != nil {
            log.Println("Error reading Kafka message:", err)
            continue
        }

        var purchase Purchase
        err = json.Unmarshal(msg.Value, &purchase)
        if err != nil {
            log.Println("Error unmarshaling Kafka message:", err)
            continue
        }

        log.Printf("Received purchase event: %+v\n", purchase)
        // Process the purchase event (e.g., update inventory)
    }
}
