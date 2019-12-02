3.(dev) Projete (e implemente um protótipo) de um escalonador de tarefas em lote.
Descreva as simplificações adotadas no protótipo (caso seja mais simples que o projeto). 
Considere que a tarefa é um tipo que inclui um ID e uma string que descreve o comando a ser executado.
Tarefas são executadas pelo tipo worker. O resultado da execução de uma tarefa é um inteiro que indica sucesso ou não.
Esse escalonador funciona baseado em eventos. Considere os seguintes tipos de evento: 
1) adição/remoção de workers do pool considerado pelo escalonador;
2) adição/remoção de tarefas;
3) finalização da execução de tarefas.

O worker não precisa executar os comandos das tarefas de fato (ao invés disso, simule a duração da execução.
 considere uma escala de unidades de segundos).
Considere que o worker notifica o escalonador sobre os resultados das tarefas (ao invés do scheduler fazer pooling para obtenção dos resultados).

O escalonador possui uma estrutura de dados a qual armazena as atividades em uma fila e é o responsável por executá-las(definir de quem será a CPU no momento x), para isso é necessário um worker o qual adiciona e remove as tarefas na fila e quem executa a tarefa.

Função do worker:
* Adicionar ou remover tarefas (cada worker é responsável por uma fila de tarefas)
* Realizar a função da tarefa

Função do escalonador:
* Atribuir um worker para uma atividade


As notificações são emitidas pelo sistima, o escalonador vai receber o evento e tomar as devidas ações:
* adicionar
    Atualiza a lista de atividades
    Adiciona um worker para determinada atividade
* 