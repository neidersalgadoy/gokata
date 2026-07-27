# Knowledge Tree — gokata

> No es lista plana. Es árbol. Cada respuesta abre ramas más profundas.

## 🌳 Estructura del árbol

```
Go
├── Concurrencia
│   ├── Gorutinas (qué son, costo, lifecycle)  son proceso livianos manjados por el scheduler d go que a diferencia de hilos una go rutina o multiple go rutinas corren en un hilo pero el scheduler escoje cuando ejeuctar
│   │   ├── Scheduler G/M/P  son como los partes que usa schedulerpara escojer que go rutina tien recuerso y ejeuta cunado duerme una etc..
│   │   │   ├── sysmon
│   │   │   ├── netpoller
│   │   │   └── work stealing  robo de la carga de trabajao hay un patorn de ocncurrencia peorasumo esto es ma concepp internal
│   │   └── Goroutine leaks (cómo detectar, evitar) ni idea 
│   ├── Channels (buffered, unbuffered, close, nil) aca preguna isemprl a difernecia de bufered y unbufere con tirc siemrpe hay como cascaras de preguntas cn buer y unbuffer
│   │   ├── Patrones fan-out / fan-in los ocnozco pero no los implmento de cero entineod us ocncepto
│   │   ├── Pipeline pattern masomenso entiendo el concepto dde pipes no implemtnarlo ymenso ocn chann
│   │   └── Select (timeout, default, nil cases) si entiendo com ofuncion lo veo entinedol dodifoo peor no implmento ni plaeno y menso cuando podria usar ucnado es mi too para usar y caso nada solo lo entinedo a lleer.
│   ├── sync primitives (WaitGroup, Mutex, Once, Pool) las entiendo un poc pero un em flat musculo digamos peude quep onga un add en una go rutina y buen oanda me asgura que la oguritna ejeucte antes uqe acabe lem ain ejempl oeos me causo dolor de cabeza ne le pasado... asi que mi nivel esentiendo peor no soy experto para atener ocall de proceso ocmplejos..
│   ├── errgroup (cancelación, propagación errores) nuevo para mi me gustariaaprnder eeror handlin de cero a hero
│   ├── singleflight (cache stampede)  no se que es
│   └── atomic (CAS, Load/Store, false sharing) no seque es
├── Memoria
│   ├── Stack vs Heap (escape analysis). stak lo que va en scope  y heap es la que mas permanece referencia mas permante stask acaba sale ucando no se ua defer trabaja por tack,
│   ├── GC (mark-sweep, STW, GOGC, pausas) no se que es
│   └── pprof / trace (heap, CPU, goroutine, mutex) si me han preugntado jamas he abido decir ni m
├── APIs
│   ├── REST las amo lo mas ocmuno recien sacaron query verb repasemos bien todo versiones eeguritad token headers todo como jamas tuve escuela en testo fue mas lo que aprnedo en handoo y diario a diario pepr nada omral segunro estnadar 
│   ├── gRPC igual si aplrendi pero solo lo basico mi sosprcracuando no oso existi protobugsino thrit omg. que no sooera pide dame isn ocomo muitple diretente streams raro cuotas rado velocidad esto a nviel avanzado debe ser moeusnto yyo oslo protobu serviceim pemtne llamen cliente llame rpc ejecuta devuevle  un similarapai rest pero por hhtp2.0 ajajaj  uh contor dem emodira gorutina errores manejo repsueta gateways seguridad replcia disoverin nada de eos y esucho muhco mas ocnepctoque hum..
│   └── Middleware esencial nada rofund
├── Testing unitari testtable peroperdi muculo sobretodo en mb cleinte http reuqest mocks 
├── Databases sql mysq pors peor no seria capas de ejecutar ocnsutlar rapido com oesperair un senio ni meno optimizar consutla ocmplja y menos escribir transacioanl o proagmaicon qls,
├── Arquitectura
│   ├── DDD sabisocnoco la toeir
│   ├── CQRS alguna ves la trbaja un poco pero rico hacer vario claeneges a prondudad
│   ├── Event Sourcing rico hacr diferentes archi y retos complejso de ete tema tradeoof etc
│   └── SAGA ceo hazme un master
├── Sistemas Distribuidos 
│   ├── CAP / PACELC no recuedo
│   ├── Raft no e+se 
│   ├── MVCC jama lo vi
│   └── Rate Limiting epa evitar matarnos conpeticion y evitar martar con peticiones
└── Infra
    ├── Docker poquitito
    ├── K8s queiro d9minarlo entender ser coadjyudante y poder descrubir erorore pen prodccuion o sea ce -11000% y queiro ser el +3custrillones%
    └── CI/CD un pocquito ocn jenckisn pero aca cdk terraform etc ufufff amorair framwork srverles ngnex etc de todo gh actions uff muero por saber manejar git 
```

## ⚡ Dinámica

1. Yo elijo una rama y pregunto desde nivel 1
2. Respondés con lo que sabés
3. Según tu respuesta, voy por una sub-rama u otra
4. Profundizo hasta llegar a tu límite
5. Marcamos dónde termina tu conocimiento real
6. Pasamos a otra rama
