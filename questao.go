package main

import (
	"fmt"
	"runtime"
	"strconv"
	"sync"
	"time"
)

// Tarefa - tipo para tarefa que será processada
type Tarefa struct {
	ID          string
	description string
}

//Função que a tarefa executa.
func (t *Tarefa) Process() int {
	fmt.Printf("Processing job '%s'\n", t.ID)
	time.Sleep(1 * time.Second)
	return 1
}

// Worker - Faz o precessamento das tarefas
type Worker struct {
	concluida sync.WaitGroup
	atvPronta chan chan Tarefa
	tarefa    chan Tarefa

	livre chan bool
}

// Escalonador - Gerenciador das atividades
type Escalonador struct {
	tarefas           chan Tarefa
	atvPronta         chan chan Tarefa //
	workers           []*Worker
	dispatcherStopped sync.WaitGroup
	workersPausados   sync.WaitGroup
	livre             chan bool
}

// CriarEscalonador - Cria e modifica as propriedades default do escalonador
func CriarEscalonador(maxWorkers int) *Escalonador {
	workersPausados := sync.WaitGroup{}
	atvPronta := make(chan chan Tarefa, maxWorkers)
	workers := make([]*Worker, maxWorkers, maxWorkers)
	for i := 0; i < maxWorkers; i++ {
		workers[i] = CriarWorker(atvPronta, workersPausados)
	}
	return &Escalonador{
		tarefas:           make(chan Tarefa),
		atvPronta:         atvPronta,
		workers:           workers,
		dispatcherStopped: sync.WaitGroup{},
		workersPausados:   workersPausados,
		livre:             make(chan bool),
	}
}

// Executar - Inicializa todos os workes e ...
func (q *Escalonador) Executar() {
	for i := 0; i < len(q.workers); i++ {
		q.workers[i].ExecutarWorker()
	}
	go q.dispatch()
}

// Parar - Para a execução e libera a fila, sispatcher routine
func (q *Escalonador) Parar() {
	q.livre <- true
	q.dispatcherStopped.Wait()
}

func (q *Escalonador) dispatch() {
	q.dispatcherStopped.Add(1)
	for {
		select {
		case job := <-q.tarefas: // Verifica se existe alguma tarefa na fila
			workerChannel := <-q.atvPronta // Verifica se existe atividade pronta
			workerChannel <- job           // Envia a atividade para o canal
		case <-q.livre: // Se o escalonador está livre
			for i := 0; i < len(q.workers); i++ { //Pausa todos os workes
				q.workers[i].PararWorker()
			}
			q.workersPausados.Wait()
			q.dispatcherStopped.Done()
			return
		}
	}
}

// AddTarefa - Adiciona uma nova tarefa a ser processada
func (q *Escalonador) AddTarefa(job Tarefa) Tarefa {
	q.tarefas <- job
	return job
}

// RemoveTarefa - Remove uma nova tarefa a ser processada
func (q *Escalonador) RemoveTarefa(job Tarefa) {
	//job <- q.tarefas
}

// NewWorker - Cria novo worker
func CriarWorker(tarefa chan chan Tarefa, concluida sync.WaitGroup) *Worker {
	return &Worker{
		concluida: concluida,
		atvPronta: tarefa,
		tarefa:    make(chan Tarefa),
		livre:     make(chan bool),
	}
}

// RemoveWorker - Remove um worker da lista de workers
func (q *Escalonador) RemoveWorker(worker *Worker) {
	tmp := q.workers[:0]
	for _, w := range q.workers {
		if w != worker {
			tmp = append(tmp, w)
		}
	}
	q.workers = tmp
}

// ExecutarWorker - Processar tarefas do worker
func (w *Worker) ExecutarWorker() {
	go func() {
		w.concluida.Add(1)
		for {
			w.atvPronta <- w.tarefa // check the job queue in
			select {
			case job := <-w.tarefa: // see if anything has been assigned to the queue
				job.Process()
			case <-w.livre:
				w.concluida.Done()
				return
			}
		}
	}()
}

// PararWorker - Torna o worker livre
func (w *Worker) PararWorker() {
	w.livre <- true
}

//////////////// Example //////////////////

func main() {
	queue := CriarEscalonador(runtime.NumCPU())
	queue.Executar()
	defer queue.Parar()

	fmt.Println(queue.workers[1])

	for i := 0; i < 4*runtime.NumCPU(); i++ {
		tarefa := queue.AddTarefa(Tarefa{strconv.Itoa(i), "dddd"})
	}
	fmt.Printf("final job '%s'\n", tarefa.ID)
}
