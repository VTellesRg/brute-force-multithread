package main

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"os"
	"runtime"
	"sync"
	"sync/atomic"
	"time"
)

// Função que gera todas as combinações possíveis de uma string de comprimento n usando os caracteres fornecidos.
func generateText(chars string, length int, jobs chan<- string) {
    var generate func(prefix string, length int)
    
    generate = func(prefix string, length int) {
        if length == 0 {
            jobs <- prefix
            return
        }
        for _, c := range chars {
            generate(prefix+string(c), length-1)
        }
    }
    
    for i := 1; i <= length; i++ {
        generate("", i)
    }
    close(jobs)
}

func textToMD5(text string) string {
    hash := md5.Sum([]byte(text))
    return hex.EncodeToString(hash[:])
}

// Função que escreve o resultado no arquivo resultado.txt
func writeResultToFile(result, elapsed string) error {
    file, err := os.OpenFile("resultado.txt", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
    if err != nil {
        return err
    }
    defer file.Close()

    _, err = file.WriteString(fmt.Sprintf("Senha: %s\nTempo de processamento: %s\n", result, elapsed))
    if err != nil {
        return err
    }

    return nil
}

// Worker que tenta encontrar a senha correta usando força bruta.
func worker(chars string, pwd string, flag *int32, wg *sync.WaitGroup, jobs <-chan string, resultChan chan<- string) {
    defer wg.Done()

    for text := range jobs {
        if atomic.LoadInt32(flag) == 1 {
            return
        }
        if textToMD5(text) == pwd {
            atomic.StoreInt32(flag, 1)
            resultChan <- text
            return
        }
    }
}

func main() {
    runtime.GOMAXPROCS(8)
    var wg sync.WaitGroup
    var flag int32 = 0
    chars := "abcdefghijlmnopqrstuvwxzABCDEFGHIJLMNOPQRSTUVWXZ0123456789#$%&*+-.*="
    var length int
    fmt.Print("Digite o tamanho da senha: ")
    fmt.Scanln(&length)
    var pwd string
    fmt.Print("Digite o hash MD5: ")
    fmt.Scanln(&pwd)

    // Canal para enviar a senha quebrada
    resultChan := make(chan string, 1)

    // Canal de trabalhos (combinações)
    jobs := make(chan string, 100)

    // Inicia a medição do tempo
    start := time.Now()

    // Criação do pool de workers
    numWorkers := 8
    for i := 0; i < numWorkers; i++ {
        wg.Add(1)
        go worker(chars, pwd, &flag, &wg, jobs, resultChan)
    }

    // Gerar combinações de senha e enviá-las para o canal de jobs
    go generateText(chars, length, jobs)

    // Espera todos os workers terminarem
    wg.Wait()
    close(resultChan)

    // Calcula o tempo de processamento
    elapsed := time.Since(start).String()

    var result string
    if flag == 1 {
        result = <-resultChan
        fmt.Println("Senha quebrada com sucesso!")
    } else {
        result = "Não foi possível quebrar a senha."
        fmt.Println(result)
    }
    // Exibe o tempo de processamento
    fmt.Printf("Tempo de processamento: %s\n", elapsed)
    // Escreve a senha quebrada e o tempo de processamento no arquivo txt
    err := writeResultToFile(result, elapsed)
    if err != nil {
        fmt.Println("Erro ao escrever no arquivo:", err)
    } else {
        fmt.Println("Resultado salvo em resultado.txt")
    }

    fmt.Println("Pressione enter para sair")
    fmt.Scanln()
}
