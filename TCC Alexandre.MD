# Algoritmo de operações de compra e venda de ativos na bolsa de valores com base em previsões obtidas com machine learning

Como havíamos conversado antes, estava procurando a melhor fonte de dados para
as cotações históricas das empresas que compõem o balcão de mercado da **B3**.
Encontrei algumas como **Quandl**, **AlphaVantage**, **yahoo Finance** e a 
própria **B3**. Nessa busca conclui que a melhor fonte seria a **B3** por 
conter todas as empresas listadas assim como dados de volume e 
OHLC (open, high, low, close) que serão
 necessários para o projeto, e não existem algumas limitações como limite de 
requisições de dados e subscrições pagas. Os dados brutos disponibilizados no site 
da **B3** são compactados em um arquivo '.txt', o que tornou necessário 
realizar um tratamento de dados para que a leitura desses fosse possível.

## Estruturação do projeto

Dividi o projeto em três etapas:

 * Tratamento dos dados obtidos no site da **B3** - o que envolve a leitura
 dos preços OHLC, que não estão com os devidos desdobramentos de cada empresa
 , ajustar os preços aplicando os desdobramentos, remoção de ativos com 
 liquidez muito baixa (existem ativos que são negociados pouquíssimas datas)
 , a remoção de ativos estrangeiros e fundos imobiliários (**FII's**). Essa 
 parte eu já fiz e o senhor pode dar uma olhada melhor nesse repositório: https://github.com/alexandrealvino/finalProject
 
 * **(TODO)** Normalizar os dados tratados, criar uma estrutura de redes neurais do tipo
  **LSTM** e alimentá-la com os dados normalizados. Ainda estou pensando qual 
  ferramenta vou utilizar para construir as redes neurais, estou estudando 
  **PyTorch** e **Tensorflow**. Ouvi falar muito bem do framework ***Keras***.
  
 * Criação de um algoritmo que execute ordens de compra e venda de ativos na 
 bolsa de valores através de uma corretora. A plataforma de *trades* escolhida
  para execução do algoritmo foi o **MetaTrader 5**. Essa parte já foi 
  desenvolvida e testada.
  
 * **(TODO)** Desenvolver uma **API** para conectar o ambiente de previsão e
  análise feita em **Python**, com a plataforma **MetaTrader 5**.
   
  

