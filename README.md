### Real time forume
#### Deep Explication

* Connexion simple entre client et serveur http :
    Serveur lorsqu'il démarre créer Un socket TCP 
        Un sokcet est Objet qui permet a un programme de communiquer sur le réseau  fourni par SE 


* Systeme  Créer un Socket TCP 
       - socket(AF_INET, SOCK_STREAM, 0); 
            AF_INET : IPv4 ;
            SOCK_STREAM : TCP 
            A ce stade le socket existe mais aucun port est le associe 
       - bind(socket, 0.0.0.0:8080):Association   d port 
            si un packet TCP , arrive sur l'interface réseau avec adresse IP et le port 8080 je le donnerai a ce socket 
                                Internet
                                    │
                                    │
                            Carte réseau (NIC)
                                    │
                                    ▼
                            Kernel (TCP/IP)
                                    │
                        +-----------+------------+
                        |                        |
                Port 80 |                Port 8080
                        |                        |
                        |                 Listening Socket


        - Socket en Mode Listen : listen(socket)
            ce socket devien un socket (listening socket)


        - Le client va créer un aussi un socket
            pr exmple le navigateur demande au systéme socket()+port libre
                A ce stade on :
                    Serveur
                        Listening Socket
                        192.168.1.10:8080                

                    Client 
                        Socket TCP (pourquoi ne fait pas mode listen car le serveur qui va attend les connexions cepandant le client connait  deja serveur  donc il fait socket()=>connect())
                        192.168.1.5:54021
                        le client lance le Three-Way Handshake; Il envoie un premier segment TCP : 
                            Client
                                │
                                │ SYN
                                ▼
                                Serveur
                            Client
                                ▲
                                │ SYN + ACK
                                Serveur
                            Client
                                │
                                │ ACK
                                ▼
                                Serveur

                            Connexion TCP = ESTABLISHED 
                                quand la connexion est established le SE creer un nv socket lie au client (chaque connexion est indépendante )


Connexion ESTBLASHID => est il ouvert tous le temps ?
    Cas 1 HTTP/1.0 : Apres la reponse du serveur ,  la connexion TCP est fermée

    Cas 2 HTTP/1.1 : La connexion reste ouverte Apres la premiere réponse (KEEP-ALIVE) mais selon TimeOut => 1 Three-Way-Handshake
        la connexion est fermé  lorsque fermeture d 'onglet && quitte navigateur => Client envoie FIN 

    Cas 3 WEBSOCKET : la connexion entre client A et B reste ouverte jusqu'a  le client ou serveur decide d fermeture
        Une seul connexion TCP => 3-Way-handshake

        browser prépare une requete HTTP (c est une requete HTTP mais il veutle changement d protocole )
            GET /chat HTTP/1.1 
            Host: localhost:8080

            Upgrade: websocket
            Connection: Upgrade

            Sec-WebSocket-Key: AbCdEfGh123456

            Sec-WebSocket-Version: 13

            Le serveur vérifie plusieurs éléments.

Par exemple :

✔ la méthode est-elle GET ?

✔ Upgrade: websocket est-il présent ?

✔ Connection: Upgrade est-il présent ?

✔ Sec-WebSocket-Version est-elle supportée ?

✔ Sec-WebSocket-Key est-elle présente ?

Si tout est correct...

Le serveur accepte et répond 
    HTTP/1.1 101 Switching Protocols

    Upgrade: websocket
    Connection: Upgrade

    Sec-WebSocket-Accept: HSmrc0sMlYUkAGmm5OPpG2HaGWk=

        TCP => 101 (Web Socket )

        Échange WebSocket 
            Le client envoie plus sous forme de requete HTTP ; il construit une WebSocket Frame :
                WebSocket Frame

                FIN = 1

                Opcode = TEXT

                Payload Length = 7

                Payload = Bonjour



//Architecture  d chat : 
    pr exemple ALICE && BOB 
        nrmlmt les clients ne communiquent jamais direct , il communiquent directem au serveur , et celui qui distribue le message pour l'autre client 
            si ALice ouvre le chat entre lui et le client => localhost:8080/chat/client1
                 const ws = new WebSocket(route) => une connexion TCP est créer ,  un webSocket est ouvert 


comment le back end recoit cette connexion et appartient a qui ? 
garder cette connexion d client ? 
Sans Gorilla

Tu devrais toi-même :

lire les en-têtes HTTP ;
vérifier Upgrade, Connection, Sec-WebSocket-Key, etc. ;
calculer Sec-WebSocket-Accept (SHA-1 + Base64) ;
envoyer la réponse 101 Switching Protocols ;
lire les frames WebSocket ;
décoder les frames ;
gérer le masking ;
gérer les Ping/Pong ;
gérer la fermeture de la connexion.

C'est beaucoup de travail.



WebSOcket Frame :
+----------------+
| Frame Header   |
+----------------+
| Type: Text     |
+----------------+
| Payload        |
| "Bonjour Bob"  |
+----------------+


client A : > Creation d une connexion webSocket avec Server 



RQ: la diff entre sync.Mutex && sync.RWMutex

+---------------------------+
|        Channel            |
|---------------------------|
| file d'attente            |
| capacité                  |
| goroutines en attente     |
+---------------------------+
Un channel est donc un objet qui appartient au runtime Go.

ch := make(chan int) => Creation d channel (ch contient une référence vers ce objet channel )

<- envoyer une valeur ds le channel 

Problematique:
            Connexion Client 5

                       ▲
                       │
          ┌────────────┼────────────┐
          │            │            │
          │            │            │
Goroutine Client 8   Broadcast   (autres...)



Un Client => Une connexion WebSocket => Go créer un objet con (websocket.Conn)
    Une 1er goroutine veut ecrire : con.WriteJson("bonjour")
    une 2 eme arrive : con.WriteJSON("Salut")
                   websocket.Conn

                        ▲
                        │
                ┌────────┴────────┐
                │                 │

            Goroutine A       Goroutine B

            la bib Gorilla Websocket ne supporte pas 2 ecritures simultanés concurrent write to websocket connection