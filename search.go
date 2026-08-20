package main

import (
	"encoding/binary"
	"encoding/json"
	"fmt"
	"log"
	"math"
	"net/http"
	"strings"
)
func createIndex(){

	_, err:= rdb.Do(ctx, "FT.CREATE", "idx:links",
		"ON", "HASH",
		"PREFIX", "1", "link:",
		"SCHEMA",
		"url", "TEXT",
		"description", "TEXT",
		"embedding", "VECTOR", "FLAT", "6",
		"TYPE", "FLOAT32",
		"DIM", "768",
		"DISTANCE_METRIC", "COSINE",
	).Result()

	if err!=nil{
		log.Println("index already exists or error:", err)
	}else{
		fmt.Println("Search index created")
	}
}
func getEmbedding(text string) ([]byte, error){
	body:=fmt.Sprintf(`{"model": "nomic-embed-text", "input": "%s"}`, text)

	resp, err:=http.Post("http://localhost:11434/api/embed", "application/json", strings.NewReader(body))

	if err!=nil{
		return nil, err
	}
	defer resp.Body.Close()

	var result struct{
		Embeddings [][]float64 `json:"embeddings"`
	}
	if err:=json.NewDecoder(resp.Body).Decode(&result); err!=nil{
		return nil, err
	}
	buf:=make([]byte, len(result.Embeddings[0])*4)
	for i, v:=range result.Embeddings[0]{
		binary.LittleEndian.PutUint32(buf[i*4:], math.Float32bits(float32(v)))
	}
	return buf, nil
}