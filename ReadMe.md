# Brute Force Multithread

Este projeto é um script Go que realiza um ataque de força bruta multithread para encontrar uma senha a partir de seu hash MD5. O script gera todas as combinações possíveis de caracteres e compara o hash MD5 de cada combinação com o hash fornecido.

## Estrutura do Projeto

- `bruteForceMultithread.go`: O arquivo principal contendo o código fonte do script.
- `bruteForceMultithread.exe`: O executável gerado após a compilação do script.
- `resultado.txt`: Arquivo onde os resultados são armazenados.

## Funcionalidades

- **Força Bruta Multithread**: O script utiliza goroutines para realizar tentativas de senha em paralelo, aumentando a eficiência do ataque de força bruta.
- **Geração de Combinações**: Gera todas as combinações possíveis de caracteres até um comprimento especificado.
- **Verificação de Hash**: Compara o hash MD5 de cada combinação gerada com o hash fornecido.
- **Armazenamento de Resultados**: Escreve a senha encontrada e o tempo de processamento no arquivo `resultado.txt`.

## Uso

### Pré-requisitos

- Go instalado na máquina.

### Compilação

Para compilar o script, execute o seguinte comando no terminal:

```sh
go build bruteForceMultithread.go
```

Para compilar em um .exe no powershell do windows:
```sh
go build -o bruteForceMultithread.exe bruteForceMultithread.go
```

## Detalhes e desafios

Inicialmente, a função utilizada para comparar os hashes estava causando falhas na aplicação. Para contornar esse problema, foi implementado um sistema de `workers` que distribui as tarefas de comparação entre várias goroutines, garantindo uma execução mais estável e eficiente.

Além disso, para evitar a sobrecarga da máquina durante a execução do script, foi adicionado o comando `runtime.GOMAXPROCS(10)`, que limita o processamento a `10 núcleos`. Isso ajuda a balancear a carga de trabalho e a manter a performance do sistema durante o ataque de força bruta.

A fim de agilizar o trabalho, o programa foi rodado em 2 máquinas diferentes (ambas não possuem núcleos virtualizados):
- **Máquina 1**: 
    - Processador: Intel Core Ultra 7 165U 2.10 GHz 12 Cores (rodando em overclock 3.4 GHz)
    - Memória RAM: 32GB
    - Sistema Operacional: Windows 11
     
- **Máquina 2**: 
    - Processador: AMD Ryzen 5 3500X 3.60 GHz 6 Cores
    - Memória RAM: 32GB
    - Sistema Operacional: Windows 10

Os hashs de senhas de 5 dígitos foi rodado na `Máquina 2`, enquanto a `Máquina 1` rodava o primeiro Hash de senha de 6 dígitos. Antes da `Máquina 1` terminar o primeiro hash de 6 dígitos, a `Máquina 2` já estava trabalhando no segundo, porém após 24h rodando encerrei o processamento na `Máquina 2` e iniciei o mesmo hash na `Máquina 1`.

### Por que `Go`?

Go é uma excelente escolha para um script de força bruta multithread por vários motivos:

1. **Concorrência Simples e Eficiente**: Go possui goroutines, que são threads leves gerenciadas pelo runtime da linguagem, permitindo a execução de milhares de tarefas em paralelo sem sobrecarregar o sistema.

2. **Desempenho**: Sendo uma linguagem compilada, Go oferece execução rápida e um garbage collector eficiente que mantém a performance alta.

3. **Simplicidade e Legibilidade**: A sintaxe simples e direta de Go facilita a escrita e manutenção do código, essencial para scripts complexos.

4. **Bibliotecas Padrão Poderosas**: Go inclui pacotes para manipulação de strings, geração de hashes (como MD5), e manipulação de arquivos, reduzindo a necessidade de dependências externas.

5. **Gerenciamento de Recursos**: Com `runtime.GOMAXPROCS`, Go permite limitar o número de threads do sistema, evitando sobrecarga durante a execução de tarefas intensivas.

6. **Portabilidade**: Go é altamente portátil, permitindo a compilação do script para diferentes sistemas operacionais e arquiteturas de forma simples.

A combinação desses fatores faz de Go uma escolha ideal para scripts que exigem alta performance e execução paralela, como um ataque de força bruta multithread.

## Refatoração v0.2

### Alterações realizadas para essa versão:
- Receber input de todos os hashs ao iniciar o programa como solicitado o exercício
- Correção do dicionário de dados, pois gera uma diferença com relação ao resultado dos colegas
- Formato do dicionário alterado para um array, supostamente torna o processamento um pouco mais rápido
- Função para impedir o computador de entrar em suspensão, supostamente caso o computador entre em suspensão o processamento pode ser interrompido

## Conclusão

A ordem do dicionário de dados afeta diretamente o processamento de senhas, para uma aplicação real do código seria válido ordenar o dicionário de forma inteligente. outro fator interessante é a semelhança de processamento entre um processador de 12 núcleos de laptop e um de 6 núcleos de desktop. neste projeto também foi realizada uma branch com prefixo de caracter especial e numérico, porém não identifiquei melhora na velocidade de processamento.

### Melhor caminho para otimizar o processamento de brute force

Por fim, com base em tudo que estudei para otimizar esse script, concluo que o melhor caminho para tornar o processamento mais rápido é a utilização dos núcleos da GPU para o processamento, porém não foi possível a implementação por incompatibilidade da minha GPU com softwares e bibliotecas disponíveis hoje em dia, por questões de modelo e fabricante (AMD Radeon RX550).
