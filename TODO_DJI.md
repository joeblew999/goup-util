# DJI and NATS Jetsream 




"Use the DJI Cloud API v2 (Thing Model) schema. I am building a Go backend using NATS JetStream as an MQTT bridge. Map the MQTT topic thing/product/{sn}/osd to the NATS subject dji.osd.{sn}. Generate Go structs for the Aircraft Properties found in the DJI documentation and implement a handler for the fly_to_point method."

This guide provides a structured breakdown of the DJI Cloud API (v2) for your AI to use as a reference. This is the modern, MQTT-based standard for recent DJI drones (Mavic 3E, Matrice 350, Dock 2, etc.).


Here is the updated, verified documentation list for the DJI Cloud API v2 (2025/2026). DJI recently restructured their developer portal, which caused the old "tutorial" links to break.

These are the live URLs for your AI to ingest.

📋 DJI Cloud API Reference (MQTT Standard)
Resource	Verified Live Link	Notes for AI
Topic Definition	Topic Structure Guide	Defines the thing/product/{sn}/osd and services hierarchy.
Message Model	JSON Wrapper Syntax	Shows the mandatory tid, bid, and timestamp root objects.
Aircraft OSD Data	Full OSD JSON Properties	Goldmine: Lists all telemetry keys (lat, lon, alt, battery, etc.).
Flight Control	Remote Control Methods	Contains JSON for fly_to_point and takeoff_to_point.
Gimbal/Camera	Payload Control Specs	Markdown source for gimbal rotation and camera trigger JSON.
Dock 2/3 Specifics	M3D/M3TD Properties	Specific properties for newer Dock-compatible drones.
🛠️ Strategic Notes for AI Implementation
Thing Model Logic: DJI uses a "TSL" (Thing Specification Language). Remind your AI that properties are pushed by the drone (OSD), while services are called by the cloud (Commands).
OSD Frequency: The osd topic publishes at 0.5Hz to 5Hz. In Go, you should use a NATS JetStream Pull Consumer with a small batch size to handle this high-frequency stream without blocking.
The Authority Handshake: Most critical commands will fail with an error code unless you first call the flight_control_lock service.
Method: flight_control_lock
Data: {"action": 1}
Error Handling: All service calls return a response on thing/product/{sn}/services_reply. Your AI must implement a "Wait-for-Reply" logic in Go using NATS request-reply or a correlation ID (the tid).
📦 Essential GitHub Repos for Context
If your AI needs to see "real world" JSON examples from the source:
Official Doc Repo: dji-sdk/Cloud-API-Doc (Best for finding .md files with raw JSON blocks).
Sample Payloads: Look in the /docs/en/ folder of the repo above for the most granular method definitions.
Usage Tip: Copy-paste this table into your AI prompt and say: "Reference these URLs to build a Go-based NATS JetStream handler for the DJI Cloud API v2. Focus on the Aircraft Properties for telemetry and the DRC methods for flight control."


https://developer.dji.com/doc/cloud-api-tutorial/en/api-reference/pilot-to-cloud/mqtt/topic-definition.html

https://developer.dji.com/doc/cloud-api-tutorial/en/overview/basic-concept/mqtt.html

https://developer.dji.com/doc/cloud-api-tutorial/en/api-reference/pilot-to-cloud/mqtt/others/aircraft/properties.html

https://developer.dji.com/doc/cloud-api-tutorial/en/feature-set/dock-feature-set/drc.html

https://github.com/dji-sdk/Cloud-API-Doc/blob/master/docs/en/30.feature-set/10.pilot-feature-set/90.drc.md

https://developer.dji.com/doc/cloud-api-tutorial/en/api-reference/dock-to-cloud/mqtt/aircraft/m3d-properties.html