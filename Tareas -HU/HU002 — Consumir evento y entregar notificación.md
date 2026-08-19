# HU002 — Consumir evento y entregar notificación (AMQP, idempotencia)

## 1. ¿Qué hace esta HU?

A diferencia de HU001 (donde alguien te pide algo por internet), acá
**nadie te pide nada directamente**. Otro sistema avisa que "pasó
algo" (se publicó un horario, o se disparó una alerta), y el
`notification-worker` está siempre escuchando para reaccionar solo,
armar la notificación correspondiente, y guardarla — sin repetirla
si por error le llega el mismo aviso más de una vez.

## 2. Cómo la entendí (en simple)

Lo pensé como un cartero dejando cartas en un buzón:

1. **Llega una carta al buzón** — otro sistema publica un evento (por
   ejemplo "se publicó el horario X" o "hay una alerta de
   capacidad").
2. **Alguien revisa el buzón todo el tiempo** ( consumer.go) — el worker está
   siempre escuchando, no espera que nadie se lo pida.
3. **Abre la carta y la traduce** a algo que el sistema entiende.
4. **El organizador decide qué hacer:** (consume_domain_event.go)
   - Busca a quién hay que avisar.
   - Si hay una plantilla guardada para ese tipo de aviso, arma el
     mensaje con ella.
   - Intenta "enviar" la notificación.
   - Guarda el resultado (si se envió bien, o si falló).
5. **No repetir la misma carta dos veces:** antes de guardar, la
   base de datos ya tiene puesta una regla que dice "si esta carta
   ya la procesé antes, no hagas nada de nuevo". Eso evita mandar la
   misma notificación duplicada si el mensaje llega repetido (algo
   normal en sistemas de mensajería).

**Diferencia clave con HU001:** ahí vos esperabas una respuesta
(202/503). Acá no hay nadie esperando — el worker toma el aviso, lo
procesa, y sigue con el siguiente, sin devolverle nada a nadie.

## 3. Diagrama de secuencia (simple)

```mermaid
sequenceDiagram
    actor RMQ as El cartero (RabbitMQ)
    participant E as El que revisa el buzón
    participant O as El organizador
    participant B as La base de datos

    RMQ->>E: deja una carta (evento nuevo)
    E->>E: abre y traduce la carta
    E->>O: le pasa el aviso ya traducido
    O->>O: busca a quién avisar
    O->>O: arma el mensaje (con plantilla si hay)
    O->>B: "¿ya procesé esta carta antes?"

    alt carta repetida (idempotencia)
        B--)O: ya existe, no la proceses de nuevo
        O--)E: listo, no hago nada más
    else carta nueva
        B-->>O: guardada con éxito
        O-->>E: listo
    end
```

## 4. Evidencia — cómo lo probé

Como todavía no existe un "publicador de horarios" real (ese es otro
microservicio que no está armado), le mandé una carta de prueba a
mano desde el panel de administración de RabbitMQ
(`http://localhost:15672`, pestaña **Exchanges** →
`monitoring-events` → **Publish message**), con este contenido:

**Routing key:** `monitoring.alert.triggered`

**Payload:**
```json
{
  "event_id": "aaaaaaaa-1111-1111-1111-111111111111",
  "event_type": "monitoring.alert.triggered",
  "source_service": "monitoring-service-test",
  "timestamp": "2026-08-18T20:00:00Z",
  "version": "1.0",
  "payload": {
    "affected_entity_type": "Ficha",
    "affected_entity_id": "bbbbbbbb-2222-2222-2222-222222222222",
    "alert_type_code": "CPU_HIGH"
  }
}
```
![alt text](image-5.png)

Con el `notification-worker` corriendo:
```powershell
$env:NOTIFICATION_DB_DSN="postgres://design_software_user:change-me@localhost:15432/design-software-develop?sslmode=disable"
$env:NOTIFICATION_AMQP_URL="amqp://design_software_user:change-me@localhost:5672/"
$env:NOTIFICATION_SMTP_ADDR="localhost:1025"
go run ./cmd/notification-worker
```
![alt text](image-6.png)

**Resultado esperado en la terminal:** el worker procesa el evento
(no muestra nada especial si sale bien — el silencio es buena señal;
si falla, sí lo loguea).

**Verificación en base de datos:**
```sql
select * from notification.sent_notification
where source_event_id = 'aaaaaaaa-1111-1111-1111-111111111111';
```

**Prueba de la idempotencia:** mandé la misma carta (mismo
`event_id`) una segunda vez, y confirmé que **no se creó un segundo
registro** en la tabla — quedó solo el primero.

## 5. Mejora propuesta

Hoy quien busca "a quién avisar" (`StubRecipientResolver`) es un
**reemplazo temporal** — siempre devuelve el mismo correo fijo,
porque el sistema real que tiene esos datos (`actors-service`)
todavía no existe. Mi propuesta: cuando ese servicio exista, agregar
un **caché corto** (por ejemplo, de un par de minutos) para no
preguntarle "¿quién es este destinatario?" cada vez que llega un
evento del mismo profesor o ficha en un lapso corto — hoy, tal como
está armado, cada evento dispara una consulta nueva aunque sea la
misma persona una y otra vez.