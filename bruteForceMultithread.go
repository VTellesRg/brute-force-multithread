
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
    "os/exec"
)
// Função que impede a suspensão do sistema operacional.
func preventSleep() {
	cmd := exec.Command("SetThreadExecutionState")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err := cmd.Start()
	if err != nil {
		fmt.Println("Erro ao tentar impedir a suspensão:", err)
	}
}
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
func writeResultToFile(hash, result, elapsed string) error {
    file, err := os.OpenFile("resultado.txt", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
    if err != nil {
        return err
    }
    defer file.Close()

    _, err = file.WriteString(fmt.Sprintf("\nHash: %s\nSenha: %s\nTempo de processamento: %s\n\n", hash, result, elapsed))
    if err != nil {
        return err
    }

    return nil
}

// Worker que tenta encontrar a senha correta usando força bruta.
func worker(chars string, hash string, flag *int32, wg *sync.WaitGroup, jobs <-chan string, resultChan chan<- string) {
    defer wg.Done()

    for text := range jobs {
        if atomic.LoadInt32(flag) == 1 {
            return
        }
        if textToMD5(text) == hash {
            atomic.StoreInt32(flag, 1)
            resultChan <- text
            return
        }
    }
}

func main() {
    preventSleep()
    runtime.GOMAXPROCS(10)
    chars := "0123456789#$%&*+-.*=abcdefghijlmnopqrstuvwxzABCDEFGHIJLMNOPQRSTUVWXZ"
    var length int
    fmt.Print("Digite o tamanho máximo da senha: ")
    fmt.Scanln(&length)

    var numHashes int
    fmt.Print("Digite o número de hashes para quebrar: ")
    fmt.Scanln(&numHashes)

    hashes := make([]string, numHashes)
    for i := 0; i < numHashes; i++ {
        fmt.Printf("Digite o hash MD5 %d: ", i+1)
        fmt.Scanln(&hashes[i])
    }

    for _, hash := range hashes {
        var wg sync.WaitGroup
        var flag int32 = 0

        // Canal para enviar a senha quebrada
        resultChan := make(chan string, 1)

        // Canal de trabalhos (combinações)
        jobs := make(chan string, 100)

        // Inicia a medição do tempo
        start := time.Now()

        // Criação do pool de workers
        numWorkers := runtime.NumCPU()
        for i := 0; i < numWorkers; i++ {
            wg.Add(1)
            go worker(chars, hash, &flag, &wg, jobs, resultChan)
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
            fmt.Printf("Hash %s quebrado com sucesso! Senha: %s\n", hash, result)
        } else {
            result = "Não foi possível quebrar a senha."
            fmt.Printf("Hash %s: %s\n", hash, result)
        }

        // Escreve a senha quebrada e o tempo de processamento no arquivo txt
        err := writeResultToFile(hash, result, elapsed)
        if err != nil {
            fmt.Println("Erro ao escrever no arquivo:", err)
        } else {
            fmt.Printf("Resultado para o hash %s salvo em resultado.txt\n", hash)
        }
    }

    fmt.Println("Todos os hashes foram processados.")
    fmt.Println("Pressione enter para sair")
    fmt.Scanln()
}