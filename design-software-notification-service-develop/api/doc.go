// El paquete api contiene los DTOs (structs) generados a partir de shared-contracts (SDD).
// No edites a mano los archivos generados; regenéralos con los comandos de abajo.
//
// notification.gen.go, generado desde el contrato OpenAPI:
//
//	oapi-codegen -generate types -package api -o api/notification.gen.go \
//	  ../design-software-shared-contracts/openapi/notification.yaml
//
// event_envelope.gen.go, generado desde el JSON Schema del "sobre" de eventos de dominio:
//
//	go-jsonschema -p api -t --capitalization ID -o api/event_envelope.gen.go \
//	  ../design-software-shared-contracts/events/event-envelope.schema.json
//
// NOTA (2026-08-15): el archivo shared-contracts/events/event-envelope.schema.json
// actualmente tiene un JSON inválido (el regex del campo "pattern" usa `\.` en vez de
// `\\.`, un escape de JSON inválido — confirmado corriendo
// `python3 -c "import json; json.load(open(...))"`). go-jsonschema no puede parsearlo
// tal cual está. El archivo event_envelope.gen.go de aquí se generó desde una copia
// parcheada localmente (solo se arregló ese escape) — hay que regenerarlo de la misma
// forma una vez se arregle el bug en shared-contracts. Esta es una dependencia de
// solo-lectura para notification-service; el arreglo va en el repo shared-contracts,
// no aquí.
package api
