// Código generado por github.com/atombender/go-jsonschema, NO EDITAR A MANO.

package api

import "encoding/json"
import "fmt"
import "regexp"
import "time"

// DomainEventEnvelope es el "sobre" estándar que envuelve todo evento de dominio
// publicado por cualquier microservicio (vía outbox -> RabbitMQ).
type DomainEventEnvelope struct {
	// CorrelationID: opcional. Sirve para rastrear una cadena de eventos
	// relacionados entre varios servicios (mismo flujo de negocio).
	CorrelationID *string `json:"correlation_id,omitempty,omitzero" yaml:"correlation_id,omitempty" mapstructure:"correlation_id,omitempty"`

	// EventID: UUID único del evento. Sirve para idempotencia (evitar
	// procesar el mismo evento dos veces si RabbitMQ lo reentrega).
	EventID string `json:"event_id" yaml:"event_id" mapstructure:"event_id"`

	// EventType: sigue el formato fijo <servicio>.<entidad>.<accion>
	// ej: "notification.sent_notification.created"
	EventType string `json:"event_type" yaml:"event_type" mapstructure:"event_type"`

	// Payload: el contenido real del evento. Es genérico (mapa clave-valor)
	// porque cada tipo de evento trae datos distintos.
	Payload DomainEventEnvelopePayload `json:"payload" yaml:"payload" mapstructure:"payload"`

	// SourceService: nombre del microservicio que publicó el evento.
	SourceService string `json:"source_service" yaml:"source_service" mapstructure:"source_service"`

	// Timestamp: fecha/hora en que se generó el evento.
	Timestamp time.Time `json:"timestamp" yaml:"timestamp" mapstructure:"timestamp"`

	// Version: versión del contrato del evento (para compatibilidad futura).
	Version string `json:"version" yaml:"version" mapstructure:"version"`
}

// DomainEventEnvelopePayload es un mapa genérico (equivalente a un JSON dinámico),
// porque el contenido varía según el tipo de evento.
type DomainEventEnvelopePayload map[string]interface{}

// UnmarshalJSON sobrescribe la conversión automática de JSON -> struct para
// agregar validaciones manuales que Go no puede hacer solo con tags.
func (j *DomainEventEnvelope) UnmarshalJSON(value []byte) error {
	// 1. Primero convierte el JSON a un mapa genérico para poder inspeccionar
	//    qué campos llegaron realmente.
	var raw map[string]interface{}
	if err := json.Unmarshal(value, &raw); err != nil {
		return err
	}

	// 2. Valida manualmente que cada campo obligatorio esté presente.
	//    Si falta alguno, devuelve un error explícito diciendo cuál.
	if _, ok := raw["event_id"]; raw != nil && !ok {
		return fmt.Errorf("field event_id in DomainEventEnvelope: required")
	}
	if _, ok := raw["event_type"]; raw != nil && !ok {
		return fmt.Errorf("field event_type in DomainEventEnvelope: required")
	}
	if _, ok := raw["payload"]; raw != nil && !ok {
		return fmt.Errorf("field payload in DomainEventEnvelope: required")
	}
	if _, ok := raw["source_service"]; raw != nil && !ok {
		return fmt.Errorf("field source_service in DomainEventEnvelope: required")
	}
	if _, ok := raw["timestamp"]; raw != nil && !ok {
		return fmt.Errorf("field timestamp in DomainEventEnvelope: required")
	}
	if _, ok := raw["version"]; raw != nil && !ok {
		return fmt.Errorf("field version in DomainEventEnvelope: required")
	}

	// 3. Truco: "Plain" es un tipo idéntico a DomainEventEnvelope pero SIN este
	//    método UnmarshalJSON, para evitar que se llame a sí mismo en bucle infinito.
	type Plain DomainEventEnvelope
	var plain Plain
	if err := json.Unmarshal(value, &plain); err != nil {
		return err
	}

	// 4. Valida con regex que EventType cumpla el formato "algo.algo.algo"
	//    en minúsculas (ej: notification.sent_notification.created).
	if matched, _ := regexp.MatchString(`^[a-z_]+\.[a-z_]+\.[a-z_]+$`, string(plain.EventType)); !matched {
		return fmt.Errorf("field %s pattern match: must match %s", "EventType", `^[a-z_]+\.[a-z_]+\.[a-z_]+$`)
	}

	// 5. Si todo pasó, asigna los valores validados al struct real.
	*j = DomainEventEnvelope(plain)
	return nil
}