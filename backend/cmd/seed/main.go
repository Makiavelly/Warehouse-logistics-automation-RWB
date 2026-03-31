// Инструмент загрузки parquet-данных в БД через REST API.
//
// Использование:
//
//	go run ./cmd/seed -file ../../train_team_track.parquet
//	go run ./cmd/seed -file ../../test_team_track.parquet -url http://localhost:8080 -key internal-api-key -batch 500
package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/parquet-go/parquet-go"
)

type Row struct {
	RouteID      string   `parquet:"route_id"`
	OfficeFromID string   `parquet:"office_from_id"`
	Timestamp    int64    `parquet:"timestamp"`
	Status1      *float64 `parquet:"status_1"`
	Status2      *float64 `parquet:"status_2"`
	Status3      *float64 `parquet:"status_3"`
	Status4      *float64 `parquet:"status_4"`
	Status5      *float64 `parquet:"status_5"`
	Status6      *float64 `parquet:"status_6"`
	Status7      *float64 `parquet:"status_7"`
	Status8      *float64 `parquet:"status_8"`
	Target2H     *float64 `parquet:"target_2h"`
}

type DataPoint struct {
	RouteID      string   `json:"route_id"`
	OfficeFromID string   `json:"office_from_id"`
	Timestamp    string   `json:"timestamp"`
	Status1      *float64 `json:"status_1"`
	Status2      *float64 `json:"status_2"`
	Status3      *float64 `json:"status_3"`
	Status4      *float64 `json:"status_4"`
	Status5      *float64 `json:"status_5"`
	Status6      *float64 `json:"status_6"`
	Status7      *float64 `json:"status_7"`
	Status8      *float64 `json:"status_8"`
	Target2H     *float64 `json:"target_2h,omitempty"`
}

type IngestRequest struct {
	DataPoints []DataPoint `json:"data_points"`
}

type IngestResponse struct {
	Inserted int `json:"inserted"`
}

func main() {
	filePath := flag.String("file", "", "путь к .parquet файлу (обязательно)")
	apiURL := flag.String("url", "http://localhost:8080", "базовый URL API")
	apiKey := flag.String("key", "internal-api-key", "API-ключ (X-API-Key)")
	batchSize := flag.Int("batch", 500, "размер батча")
	flag.Parse()

	if *filePath == "" {
		flag.Usage()
		os.Exit(1)
	}

	f, err := os.Open(*filePath)
	if err != nil {
		log.Fatalf("не удалось открыть файл: %v", err)
	}
	defer f.Close()

	stat, _ := f.Stat()
	reader := parquet.NewGenericReader[Row](f, parquet.SchemaOf(new(Row)))
	defer reader.Close()

	log.Printf("файл: %s (%d байт)", stat.Name(), stat.Size())

	var (
		batch   []DataPoint
		total   int
		rows    = make([]Row, *batchSize)
		client  = &http.Client{Timeout: 30 * time.Second}
		ingest  = *apiURL + "/api/data/ingest"
	)

	for {
		n, err := reader.Read(rows)
		for i := 0; i < n; i++ {
			r := rows[i]
			ts := time.Unix(0, r.Timestamp).UTC().Format(time.RFC3339Nano)
			batch = append(batch, DataPoint{
				RouteID:      r.RouteID,
				OfficeFromID: r.OfficeFromID,
				Timestamp:    ts,
				Status1:      r.Status1,
				Status2:      r.Status2,
				Status3:      r.Status3,
				Status4:      r.Status4,
				Status5:      r.Status5,
				Status6:      r.Status6,
				Status7:      r.Status7,
				Status8:      r.Status8,
				Target2H:     r.Target2H,
			})

			if len(batch) >= *batchSize {
				inserted := sendBatch(client, ingest, *apiKey, batch)
				total += inserted
				fmt.Printf("\rотправлено: %d", total)
				batch = batch[:0]
			}
		}
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatalf("ошибка чтения: %v", err)
		}
	}

	if len(batch) > 0 {
		total += sendBatch(client, ingest, *apiKey, batch)
	}

	fmt.Printf("\nГотово. Загружено в БД: %d записей\n", total)
}

func sendBatch(client *http.Client, url, apiKey string, points []DataPoint) int {
	body, _ := json.Marshal(IngestRequest{DataPoints: points})
	req, _ := http.NewRequest(http.MethodPost, url, bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-API-Key", apiKey)

	resp, err := client.Do(req)
	if err != nil {
		log.Fatalf("HTTP ошибка: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		b, _ := io.ReadAll(resp.Body)
		log.Fatalf("сервер вернул %d: %s", resp.StatusCode, b)
	}

	var result IngestResponse
	json.NewDecoder(resp.Body).Decode(&result)
	return result.Inserted
}